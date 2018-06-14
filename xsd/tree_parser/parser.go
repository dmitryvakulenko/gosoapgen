package tree_parser

import (
    "io"
    "encoding/xml"
    "strings"
    xsd "github.com/dmitryvakulenko/gosoapgen/xsd-model"
    "log"
    "container/list"
)

var (
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

    rootSchemas  []*xsd.Schema
    schemasStack *list.List
    resultTypes  *typesList
}

func NewParser(l Loader) *parser {
    return &parser{
        loader:       l,
        elStack:      &elementsStack{},
        nsStack:      &stringsStack{},
        curNs:        make(map[string]string),
        schemasStack: list.New(),
        resultTypes: newTypesList()}
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

func (p *parser) GetTypes() []*Type {
    p.resultTypes.Reset()

    p.generateTypes(p.rootSchemas)

    var types []*Type
    for _, t := range p.resultTypes.Iterate() {
        types = append(types, t)
        types = append(types, extractInnerTypes(t, 0)...)
    }

    resolveBaseTypes(types)
    upFieldsTypes(types)

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
        p.generateTypes(sc.ChildSchemas)
        e := p.schemasStack.PushBack(sc)
        p.schemaNode(&sc.Node)
        p.schemasStack.Remove(e)
    }
}

// func renameDuplicatedTypes() {
//     names := make(map[string]int)
//     for _, t := range globalTypesCache {
//         if _, exist := names[t.Name.Local]; exist {
//             names[t.Name.Local]++
//             t.Name.Local = t.Name.Local + "_" + strconv.Itoa(names[t.Name.Local])
//         } else {
//             names[t.Name.Local] = 0
//         }
//     }
// }

func upFieldsTypes(types []*Type) {
    for _, t := range types {
        for _, f := range t.Fields {
            f.Type = lastType(f.Type)
        }
    }
}

func lastType(t *Type) *Type {
    if (len(t.Fields) == 0 || t.isSimpleContent && len(t.Fields) == 1) && t.baseType != nil {
        return lastType(t.baseType)
    } else {
        return t
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

    if p.resultTypes.Has(qName) {
        return p.resultTypes.Get(qName)
    }

    node := p.findGlobalNode(qName)
    if node != nil {
        return p.parseSomeRootNode(qName, node)
    }

    panic("Can't find type " + name)
}

func (p *parser) parseSomeRootNode(name xml.Name, n *xsd.Node) *Type {
    if p.resultTypes.Has(name) {
        return p.resultTypes.Get(name)
    }

    var tp *Type
    switch n.Name() {
    case "element":
        tp = p.elementNode(n)
    case "simpleType":
        tp =  p.simpleTypeNode(n)
    case "complexType":
        tp =  p.complexTypeNode(n)
    case "attributeGroup":
        tp =  p.attributeGroupNode(n)
    }

    return tp
}

func (p *parser) createType(n *xsd.Node) *Type {
    sc := p.schemasStack.Back().Value.(*xsd.Schema)
    t := newType(n, sc.TargetNamespace)

    // this not global, internal type with no name
    if t.Local == "" {
        return t
    }

    if p.resultTypes.Has(t.Name) {
        log.Fatalf("Duplicated types %+v", t)
    }

    p.resultTypes.Add(t)
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
            var tp *Type
            if typName := ch.AttributeValue("type"); typName != "" {
                tp = p.findOrCreateGlobalType(typName)
            } else {
                tp = p.elementNode(ch)
            }
            f := newField(ch, tp)
            t.addField(f)
        case "attribute":
            t.addField(p.attributeNode(ch))
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
        case "attributeGroup":
            t.baseType = p.attributeGroupNode(ch)
        case "simpleContent":
            t.baseType = p.simpleContentNode(ch)
        }
    }

    return t
}

func (p *parser) simpleTypeNode(n *xsd.Node) *Type {
    tp := p.createType(n)
    tp.isSimpleContent = true

    for _, ch := range n.Children() {
        switch ch.Name() {
        case "restriction":
            tp.baseType = p.restrictionNode(ch)
        case "union":
            tp.baseType = newStandardType("string")
        }
    }


    return tp
}

func (p *parser) extensionNode(n *xsd.Node) *Type {
    tp := p.createType(n)
    base := n.AttributeValue("base")
    tp.baseType = p.findOrCreateGlobalType(base)

    for _, ch := range n.Children() {
        switch ch.Name() {
        case "sequence", "all":
            s := p.sequenceNode(ch)
            tp.append(s)
        case "attribute":
            f := p.attributeNode(ch)
            tp.addField(f)
        case "attributeGroup":
            tp.append(p.attributeGroupNode(ch))
        }
    }

    return tp
}

func (p *parser) attributeNode(n *xsd.Node) *Field {
    typName := n.AttributeValue("type")
    ch := n.ChildByName("simpleType")
    var tp *Type
    if typName != "" {
        tp = p.findOrCreateGlobalType(typName)
    } else if ch != nil {
        tp = p.simpleTypeNode(ch)
    } else {
        panic("Unknown attribute definition")
    }

    res := newField(n, tp)
    res.IsAttr = true

    return res
}

func (p *parser) attributeGroupNode(n *xsd.Node) *Type {
    tp := p.createType(n)
    name := n.AttributeValue("name")
    ref := n.AttributeValue("ref")
    if name != "" {
        for _, ch := range n.Children() {
            switch ch.Name() {
            case "attribute":
                f := p.attributeNode(ch)
                tp.addField(f)
            }
        }
    } else if ref != "" {
        tp.append(p.findOrCreateGlobalType(ref))
    } else {
        panic("No elemName and no ref for attribute group")
    }

    return tp
}

func (p *parser) endUnion() {
    p.elStack.Pop()
    context := p.elStack.GetLast()
    context.name = stringQName
}

func (p *parser) simpleContentNode(n *xsd.Node) *Type {
    tp := p.createType(n)
    for _, ch := range n.Children() {
        switch ch.Name() {
        case "restriction":
            tp.baseType = p.restrictionNode(ch)
        case "extension":
            tp.baseType = p.extensionNode(ch)
        }
    }

    return tp
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
        if _, ok := dep[t.Name]; ok || t.SourceNode.Name() == "element" && !t.referenced {
            res = append(res, t)
        }
    }

    return res
}

func (p *parser) schemaNode(n *xsd.Node) {
    for _, ch := range n.Children() {
        if ch.Name() == "include" || ch.Name() == "import" {
            continue
        }
        ns := p.schemasStack.Back().Value.(*xsd.Schema).TargetNamespace
        name := ch.AttributeValue("name")
        p.parseSomeRootNode(xml.Name{Local: name, Space: ns}, ch)
    }
}

func (p *parser) elementNode(n *xsd.Node) *Type {
    t := p.createType(n)
    elType := n.AttributeValue("type")
    ref := n.AttributeValue("ref")
    if elType != "" {
        t.baseType = p.findOrCreateGlobalType(elType)
    } else if ref != "" {
        t.baseType = p.findOrCreateGlobalType(ref)
        t.baseType.referenced = true
    } else {
        for _, ch := range n.Children() {
            switch ch.Name() {
            case "simpleType":
                t.baseType = p.simpleTypeNode(ch)
            case "complexType":
                t.baseType = p.complexTypeNode(ch)
            }
        }
    }

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

// move fields from base type to current for inheritance avoiding
func resolveBaseTypes(types []*Type) {
    for _, t := range types {
        var baseType string
        t.Fields, baseType = collectBaseFields(t)
        if t.isSimpleContent {
            t.Fields = append(t.Fields, newValueField(baseType))
        }
    }
}

func collectBaseFields(t *Type) ([]*Field, string) {
    res := make([]*Field, len(t.Fields))
    copy(res, t.Fields)

    if t.baseType == nil {
        return res, t.Local
    }

    baseFields, baseType := collectBaseFields(t.baseType)
    res = append(baseFields, res...)
    t.isSimpleContent = t.baseType.isSimpleContent

    return res, baseType
}
