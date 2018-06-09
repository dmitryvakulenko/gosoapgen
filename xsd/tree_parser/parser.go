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
    resolveBaseTypes()

    var types []*Type
    for _, t := range globalTypesCache {
        types = append(types, t)
        types = append(types, extractInnerTypes(t, 0)...)
    }

    // renameDuplicatedTypes()
    l := filterUnusedTypes(types)
    embedFields(l)

    return l
}

func extractInnerTypes(t *Type, deep int) []*Type {
    var types []*Type
    for _, f := range t.Fields {
        if (!f.Type.isSimpleContent || len(f.Type.Fields) > 0) && deep > 0 {
            types = append(types, extractInnerTypes(f.Type, deep+1)...)
        }
    }

    return types
}

// Generate types list according to previously built tree
func (p *parser) generateTypes(schemas []*xsd.Schema) {
    for _, sc := range schemas {
        e := p.schemasStack.PushBack(sc)
        p.schemaNode(&sc.Node)
        p.schemasStack.Remove(e)
    }
}

func renameDuplicatedTypes() {
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

// Find type by name. If type not found, try to find node and create type
func (p *parser) findOrCreateGlobalType(name string) *Type {
    if name == "" {
        panic("Type should has name")
    }
    qName := p.createQName(name)
    if qName.Space == "http://www.w3.org/2001/XMLSchema" {
        return &Type{Name: qName, isSimpleContent: true}
    }

    if t, ok := globalTypesCache[qName]; ok {
        return t
    }

    node := p.findGlobalNode(qName)
    if node != nil {
        t := p.parseSomeRootNode(qName, node)
        globalTypesCache[qName] = t
        return t
    }

    panic("Can't find type " + name)
}

func (p *parser) parseSomeRootNode(name xml.Name, n *xsd.Node) *Type {
    if name.Local == "" {
        panic("Root node without name")
    }

    if _, ok := globalTypesCache[name]; ok {
        return nil
    }

    switch n.Name() {
    case "element":
        return p.elementNode(n)
    case "simpleType":
        return p.simpleTypeNode(n)
    }

    return nil
}


func (p *parser) createType(n *xsd.Node) *Type {
    sc := p.schemasStack.Back().Value.(*xsd.Schema)
    t := newType(n, sc.TargetNamespace)

    // this not global, internal type with no name
    if t.Local == "" {
        return t
    }

    if _, ok := globalTypesCache[t.Name]; ok {
        log.Fatalf("Duplicated types %+v", t)
    }

    globalTypesCache[t.Name] = t
    return t
}

// Find schema node by name and element
func (p *parser) findGlobalNode(name xml.Name) *xsd.Node {
    for _, s := range p.rootSchemas {
        n := s.FindGlobalTypeByName(name)
        if n != nil {
            return n
        }
    }

    return nil
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
    sc := p.schemasStack.Back().Value.(*xsd.Schema)
    if len(typesParts) == 1 {
        name = typesParts[0]
        namespace = sc.TargetNamespace
    } else {
        namespace = sc.ResolveSpace(typesParts[0])
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

func (p *parser) sequenceNode(n *xsd.Node) *Type {
    t := p.createType(n)
    for _, ch := range n.Children() {
        switch ch.Name() {
        case "element":
            tp := p.findOrCreateGlobalType(ch.AttributeValue("type"))
            f := newField(ch, tp)
            t.addField(f)
        }
    }

    return t
}

func (p *parser) complexTypeNode(n *xsd.Node) *Type {
    var t = p.createType(n)
    for _, ch := range n.Children() {
        switch ch.Name() {
        case "sequence":
            t.baseType = p.sequenceNode(ch)
        case "attribute":
            a := p.attributeNode(ch)
            t.addField(a)
        }
    }

    return t
}

func (p *parser) simpleTypeNode(n *xsd.Node) *Type {
    tp := p.createType(n)
    tp.isSimpleContent = true

    restr := n.ChildByName("restriction")
    if restr != nil {
        tp.baseType = p.restrictionNode(restr)
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

func (p *parser) attributeNode(n *xsd.Node) *Field {
    tp := p.findOrCreateGlobalType(n.AttributeValue("type"))
    res := newField(n, tp)
    res.IsAttr = true

    return res
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
func filterUnusedTypes(types []*Type) []*Type {
    var res []*Type
    dep := buildDependencies(types)
    for _, t := range types {
        if _, ok := dep[t.Name]; ok || t.SourceNode.Name() == "element" {
            res = append(res, t)
        }
    }

    return res
}

func (p *parser) schemaNode(n *xsd.Node) {
    for _, ch := range n.Children() {
        ns := p.schemasStack.Back().Value.(*xsd.Schema).TargetNamespace
        name := ch.AttributeValue("name")
        tp := p.parseSomeRootNode(xml.Name{Local: name, Space: ns}, ch)
        if tp != nil {
            globalTypesCache[tp.Name] = tp
        }
    }
}

func (p *parser) elementNode(n *xsd.Node) *Type {
    t := p.createType(n)
    elType := n.AttributeValue("type")
    if elType != "" {
        t.baseType = p.findOrCreateGlobalType(elType)
    } else {
        ch := n.FirstChild()
        switch ch.Name() {
        case "simpleType":
            t.baseType = p.simpleTypeNode(ch)
        case "complexType":
            t.baseType = p.complexTypeNode(ch)
        }
    }

    // if st != nil {
    //     t := createType(n)
    //     t.baseTypeName = p.simpleTypeNode(st)
    // }
    //
    // base := n.AttributeValue("type")
    // if base != "" {
    //     t.baseTypeName = p.createQName(base)
    // }

    // com := n.ChildrenByName("complexType")
    return t
}

func (p *parser) restrictionNode(n *xsd.Node) *Type {
    base := n.AttributeValue("base")
    if base == "" {
        panic("Restriction element without base")
    }
    return p.findOrCreateGlobalType(base)
}

// Made embedding ref, attributeGroup fields
// also adding XMLName and Value fields
func embedFields(typs []*Type) {
    dep := buildDependencies(typs)
    for _, t := range typs {
        // adding XMLName field
        if _, ok := dep[t.Name]; !ok && t.SourceNode.Name() == "element" {
            t.Fields = append([]*Field{newXMLNameField()}, t.Fields...)
        }
    }
}

// build dependencies
func buildDependencies(types []*Type) map[xml.Name][]*Type {
    usedTypes := make(map[xml.Name][]*Type)
    for _, t := range types {
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
        t.Fields = collectBaseFields(t)
    }
}

func collectBaseFields(t *Type) []*Field {
    res := make([]*Field, len(t.Fields))
    copy(res, t.Fields)

    if t.baseType == nil {
        if t.isSimpleContent {
            res = append(res, newValueField(t.Name.Space))
        }
        return res
    }

    baseFields := collectBaseFields(t.baseType)
    res = append(baseFields, res...)
    t.isSimpleContent = t.baseType.isSimpleContent
    t.baseType = nil

    return res
}
