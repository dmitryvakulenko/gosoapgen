/*
Парсер схемы xsd.
Использует ручной обход всего дерева, вместо загрузки его с помощь Unmarshal.
 */
package tree_parser

import (
    "io"
    "encoding/xml"
    "strings"
    "strconv"
    xsd "github.com/dmitryvakulenko/gosoapgen/xsd-model"
    "log"
    "container/list"
)

var (
    globalTypesCache map[xml.Name]*Type

    stringQName = xml.Name{Local: "string", Space: "http://www.w3.org/2001/XMLSchema"}
)

// Интерфейс загрузки xsd
// должен отслеживать уже загруженные файлы
// и правильно отрабатывать относительные пути
type Loader interface {
    /*
    Загрузить файл по указанному пути (или url)
    Второй параметр - ошибка, которую должен уметь анализировать метод IsAlreadyLoadedError
     */
    Load(path string) (io.ReadCloser, error)
    IsAlreadyLoadedError(error) bool
}

type parser struct {
    loader   Loader
    elStack  *elementsStack
    nsStack  *stringsStack
    curNs    map[string]string
    rootNode *node

    rootSchemas  []*xsd.Schema
    schemasStack *list.List
}

func NewParser(l Loader) *parser {
    return &parser{
        loader:       l,
        elStack:      &elementsStack{},
        nsStack:      &stringsStack{},
        curNs:        make(map[string]string),
        rootNode:     &node{elemName: "schema"},
        schemasStack: list.New()}
}

func (p *parser) Load(inputFile string) {
    p.rootSchemas = append(p.rootSchemas, p.loadSchema(inputFile, ""))
}

func (p *parser) loadSchema(inputFile string, ns string) *xsd.Schema {
    reader, err := p.loader.Load(inputFile)
    defer reader.Close()

    var s *xsd.Schema
    if err == nil {
        s = xsd.Load(reader)
        // to processing include
        if ns != "" {
            s.TargetNamespace = ns
        }
    } else if !p.loader.IsAlreadyLoadedError(err) {
        panic(err)
    }

    for _, i := range s.ChildrenByName("include") {
        inc := p.loadSchema(i.AttributeValue("schemaLocation"), s.TargetNamespace)
        s.ChildSchemas = append(s.ChildSchemas, inc)
    }

    for _, i := range s.ChildrenByName("import") {
        inc := p.loadSchema(i.AttributeValue("schemaLocation"), "")
        s.ChildSchemas = append(s.ChildSchemas, inc)
    }

    return s
}

func (p *parser) decodeXsd(decoder *xml.Decoder) {
    for {
        token, err := decoder.Token()
        if err == io.EOF {
            break
        }

        if err != nil {
            panic(err)
        }

        switch elem := token.(type) {
        case xml.StartElement:
            p.parseStartElement(&elem)
        case xml.EndElement:
            p.parseEndElement(&elem)
        }
    }
}

func (p *parser) parseStartElement(elem *xml.StartElement) {
    switch elem.Name.Local {
    case "schema":
        p.schemaStarted(elem)
    case "node", "simpleTypeNode", "complexType", "extension", "restriction", "sequence", "attribute", "attributeGroup", "element", "union",
        "simpleContent", "complexContent", "choice":
        p.elementStarted(elem)
    case "include", "import":
        p.includeStarted(elem)
    }
}

// Начало элемента, любого, а не только node
func (p *parser) elementStarted(e *xml.StartElement) {
    p.elStack.Push(newNode(e))
}

func (p *parser) parseEndElement(elem *xml.EndElement) {
    // switch elem.Name.Local {
    // case "schema":
    //     p.nsStack.Pop()
    // case "simpleTypeNode":
    //     p.simpleTypeNode()
    // case "extension", "restriction":
    //     p.endExtensionRestriction()
    // case "complexType":
    //     p.endComplexType()
    // case "sequence":
    //     p.endSequence()
    // case "node":
    //     p.endElement()
    // case "attribute":
    //     p.endAttribute()
    // case "attributeGroup":
    //     p.endAttributeGroup()
    // case "element":
    //     p.endElement()
    // case "union":
    //     p.endUnion()
    // case "simpleContent":
    //     p.endSimpleContent()
    // case "complexContent":
    //     p.endComplexContent()
    // case "choice":
    //     p.endChoice()
    // }
}

func (p *parser) schemaStarted(elem *xml.StartElement) {
    ns := findAttributeByName(elem.Attr, "targetNamespace")
    if ns != nil {
        p.nsStack.Push(ns.Value)
    } else {
        // Используем родительский ns. А правильно ли?
        p.nsStack.Push(p.nsStack.GetLast())
    }

    // p.curNs = make(map[string]string)
    for _, attr := range elem.Attr {
        if attr.Name.Space == "xmlns" && attr.Name.Local != "" {
            p.curNs[attr.Name.Local] = attr.Value
        }
    }

    p.elStack.Push(p.rootNode)
}

func (p *parser) GetTypes() []*Type {
    globalTypesCache = make(map[xml.Name]*Type)
    p.generateTypes(p.rootSchemas)
    // p.linkTypes()
    p.renameDuplicatedTypes()
    resolveBaseTypes()
    l := filterUnusedTypes()
    embedFields(l)

    return l
}

// Generate types list according to previously built tree
func (p *parser) generateTypes(schemas []*xsd.Schema) {
    for _, sc := range schemas {
        e := p.schemasStack.PushBack(sc)
        for _, n := range sc.Children() {
            p.schemaNode(n)
        }
        p.schemasStack.Remove(e)
    }
}

func (p *parser) renameDuplicatedTypes() {
    names := make(map[string]int)
    for _, t := range globalTypesCache {
        if _, exist := names[t.Name.Local]; exist {
            names[t.Name.Local]++
            t.Name.Local = t.Name.Local + "_" + strconv.Itoa(names[t.Name.Local])
        } else {
            names[t.Name.Local] = 0
        }
    }
}

func foldSimpleTypes() {
    for _, t := range globalTypesCache {
        for _, f := range t.Fields {
            if len(f.Type.Fields) == 0 {
                f.Type = getLastType(f.Type)
            }
        }
    }
}

func getLastType(t *Type) *Type {
    if t.baseType == nil || len(t.Fields) != 0 {
        return t
    } else {
        return getLastType(t.baseType)
    }
}

func (p *parser) findTypeByName(name xml.Name) *Type {
    if name.Space == "http://www.w3.org/2001/XMLSchema" {
        return &Type{Name: name}
    }

    if t, ok := globalTypesCache[name]; ok {
        return t
    }

    panic("Can't find type " + name.Local)
}

func (p *parser) endElement() {
    e := p.elStack.Pop()
    e.name.Space = p.nsStack.GetLast()

    maxAttr := findAttributeByName(e.startElem.Attr, "maxOccurs")
    if maxAttr != nil {
        e.isArray = true
    }

    typeAttr := findAttributeByName(e.startElem.Attr, "type")
    if typeAttr != nil {
        e.typeName = p.createQName(typeAttr.Value)
    }

    refAttr := findAttributeByName(e.startElem.Attr, "ref")
    if refAttr != nil {
        e.typeName = p.createQName(refAttr.Value)
    }

    context := p.elStack.GetLast()
    context.children = append(context.children, e)

    if e.name.Local == "" && len(e.children) == 0 {
        e.name = stringQName
    }
}

func (p *parser) createQName(qName string) xml.Name {
    var name, namespace string
    typesParts := strings.Split(qName, ":")

    if len(typesParts) == 1 {
        name = typesParts[0]
        // namespace = curSchema.TargetNamespace
    } else {
        // namespace = curSchema.ResolveSpace(typesParts[0])
        name = typesParts[1]
    }

    return xml.Name{
        Local: name,
        Space: namespace}
}

func findAttributeByName(attrsList []xml.Attr, name string) *xml.Attr {
    for _, attr := range attrsList {
        if attr.Name.Local == name {
            return &attr
        }
    }

    return nil
}

func (p *parser) endSequence() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    context.children = append(context.children, e.children...)
}

func (p *parser) endComplexType() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()

    if context.elemName == "schema" {
        nameAttr := findAttributeByName(e.startElem.Attr, "name")
        e.name.Local = nameAttr.Value
        e.name.Space = p.nsStack.GetLast()
        p.rootNode.addChild(e)
    } else {
        context.isSimpleContent = e.isSimpleContent
        context.isAttr = e.isAttr
        context.children = append(context.children, e.children...)
        context.name = e.name
    }
}

func (p *parser) simpleTypeNode(n *xsd.Node) *Type {
    var tp *Type
    name := n.AttributeValue("name")
    if name != "" {
        tp = p.createAndAddType(n)
        tp.isSimpleContent = true
    }

    restr := n.ChildByName("restriction")
    if restr != nil {
        p.restrictionNode(restr)
    }

    return tp
}

func (p *parser) endExtensionRestriction() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    baseType := findAttributeByName(e.startElem.Attr, "base")
    context.name = p.createQName(baseType.Value)
    context.children = append(context.children, e.children...)
}

func (p *parser) endAttribute() {
    e := p.elStack.Pop()
    e.isAttr = true
    if e.name.Local == "" {
        typeAttr := findAttributeByName(e.startElem.Attr, "type")
        if typeAttr != nil {
            e.name = p.createQName(typeAttr.Value)
        } else {
            e.name = stringQName
        }
    }

    context := p.elStack.GetLast()
    context.children = append(context.children, e)
}

func (p *parser) endAttributeGroup() {
    e := p.elStack.Pop()
    e.isAttr = true
    e.name.Space = p.nsStack.GetLast()

    context := p.elStack.GetLast()
    nameAttr := findAttributeByName(e.startElem.Attr, "name")
    refAttr := findAttributeByName(e.startElem.Attr, "ref")
    if context.elemName == "schema" {
        e.name.Local = nameAttr.Value
        p.rootNode.addChild(e)
    } else if refAttr != nil {
        e.name = p.createQName(refAttr.Value)
        context.addChild(e)
    } else {
        panic("No elemName and no ref for attribute group")
    }
}

func (p *parser) includeStarted(e *xml.StartElement) {
    l := findAttributeByName(e.Attr, "schemaLocation")
    p.Load(l.Value)
}

func (p *parser) endUnion() {
    p.elStack.Pop()
    context := p.elStack.GetLast()
    context.name = stringQName
}

func (p *parser) endSimpleContent() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    context.isSimpleContent = true
    context.name = e.name
    context.children = append(context.children, e.children...)
}

func (p *parser) endComplexContent() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    context.isSimpleContent = false
    context.name = e.name
    context.children = append(context.children, e.children...)
}

func (p *parser) endChoice() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    context.children = append(context.children, e.children...)
}

// Remove type that made not from elements
func filterUnusedTypes() []*Type {
    var res []*Type
    dep := buildDependencies()
    for _, t := range globalTypesCache {
        if _, ok := dep[t.Name]; ok || t.SourceNode.Name() == "element" {
            res = append(res, t)
        }
    }

    return res
}

func (p *parser) createAndAddType(n *xsd.Node) *Type {
    sc := p.schemasStack.Back().Value.(*xsd.Schema)
    t := newType(n, sc.TargetNamespace)
    if _, ok := globalTypesCache[t.Name]; ok {
        log.Fatalf("Duplicated types %+v", t)
    }

    globalTypesCache[t.Name] = t
    return t
}

func (p *parser) schemaNode(n *xsd.Node) {
    for _, ch := range n.Children() {
        var tp *Type
        switch ch.Name() {
        case "element":
            tp = p.elementNode(n)
        case "simpleType":
            tp = p.simpleTypeNode(n)
        }

        if tp != nil {
            globalTypesCache[tp.Name] = tp
        }
    }
}

func (p *parser) elementNode(n *xsd.Node) *Type {
    switch ch := n.FirstChild(); ch.Name() {
    case "simpleType":
        // name := ch.AttributeValue("name")

    case "complexType":
    }
    // if st != nil {
    //     tp := createAndAddType(n)
    //     tp.baseTypeName = p.simpleTypeNode(st)
    // }
    //
    // base := n.AttributeValue("type")
    // if base != "" {
    //     tp.baseTypeName = p.createQName(base)
    // }

    // com := n.ChildrenByName("complexType")
    return nil
}

func (p *parser) restrictionNode(n *xsd.Node) xml.Name {
    base := n.AttributeValue("base")
    if base == "" {
        panic("Restriction element without base")
    }

    return p.createQName(base)
}

// Made embedding ref, attributeGroup fields
// also adding XMLName and Value fields
func embedFields(typs []*Type) {
    dep := buildDependencies()
    for _, t := range typs {
        // adding XMLName field
        if _, ok := dep[t.Name]; !ok && t.SourceNode.Name() == "element" {
            t.Fields = append([]*Field{newXMLNameField()}, t.Fields...)
        }

        // if t.isSimpleContent {
            // t.Fields = append(t.Fields, newValueField(t.baseTypeName))
        // }
    }
}

// build dependencies
func buildDependencies() map[xml.Name][]*Type {
    usedTypes := make(map[xml.Name][]*Type)
    for _, t := range globalTypesCache {
        var typeDep []*Type
        if t.baseType != nil {
            typeDep = append(typeDep, t.baseType)
        }

        for _, f := range t.Fields {
            typeDep = append(typeDep, f.Type)
        }

        for _, tp := range typeDep {
            if _, ok := usedTypes[tp.Name]; !ok {
                usedTypes[tp.Name] = []*Type{}
            }
            usedTypes[tp.Name] = append(usedTypes[tp.Name], tp)
        }
    }

    return usedTypes
}

func resolveBaseTypes() {
    for _, t := range globalTypesCache {
        t.Fields = append(collectBaseFields(t), t.Fields...)
    }
}

func collectBaseFields(t *Type) []*Field {
    var res []*Field
    if t.baseType == nil {
        return res
    }

    res = append(res, t.Fields...)
    res = append(res, collectBaseFields(t.baseType)...)
    t.isSimpleContent = t.baseType.isSimpleContent
    t.baseType = nil

    return res
}
