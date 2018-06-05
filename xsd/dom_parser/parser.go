package dom_parser

import (
    dom "github.com/subchen/go-xmldom"
    "encoding/xml"
    "container/list"
    "log"
)

var (
    schemaNs = "http://www.w3.org/2001/XMLSchema"
    tString  = &Type{Name: xml.Name{Space: "string", Local: schemaNs}}

    fXmlName = &Field{Name: "XMLName", Type: &Type{Name: xml.Name{Local: "xml.Name"}}}
)

type parser struct {
    types      map[xml.Name]*Type
    parseStack *typesStack
    nsStack    *list.List
}

func NewParser() *parser {
    return &parser{
        types:      make(map[xml.Name]*Type),
        parseStack: new(typesStack),
        nsStack:    list.New()}
}

func (p *parser) LoadFile(fileName string) {
    doc, err := dom.ParseFile(fileName)
    if err != nil {
        log.Panic(err)
    }

    p.parseNs(doc)

    for _, n := range doc.Root.Children {
        p.parseNode(n)
    }
}

func (p *parser) parseNs(n *dom.Document) {
    nsAttr := n.Root.GetAttribute("targetNamespace")
    if nsAttr != nil {
        p.nsStack.PushBack(nsAttr.Value)
    } else {
        prevNs := p.nsStack.Back().Value.(string)
        p.nsStack.PushBack(prevNs)
    }
}

func (p *parser) parseNode(n *dom.Node) {
    switch n.Name {
    // case "simpleType":
    // 	p.simpleType(n)
    case "element":
        p.element(n)
    }
}

func (p *parser) simpleType(n *dom.Node) {
    nameAttr := n.GetAttribute("name")
    if nameAttr != nil {
        // newType := p.createAndAddType(nameAttr.Value)
        // newType.IsSimpleContent = true

        // typeAttr  := n.GetAttribute("type")
        // newType.BaseTypeName = p.createAndAddType()
    }
}

func (p *parser) createAndAddType(name string) *Type {
    newType := &Type{Name: p.createName(name)}

    if _, ok := p.types[newType.Name]; ok {
        log.Panicf("Type %#v already exist", newType.Name)
    }

    p.types[newType.Name] = newType

    return newType
}

// Просто сокращение для создания xml.Name
func (p *parser) createName(name string) xml.Name {
    lastNs := p.nsStack.Back().Value.(string)
    return xml.Name{Local: name, Space: lastNs}
}

func (p *parser) GetTypes() []*Type {
    p.addXMLNames()

    var res []*Type
    for _, v := range p.types {
        res = append(res, v)
    }
    return res
}

func (p *parser) addXMLNames() {
    depMap := p.buildDependencyMap()
    for _, v := range p.types {
        if _, ok := depMap[v.Name]; !ok {
            v.Fields = append([]*Field{fXmlName}, v.Fields...)
        }
    }
}

// Просто подсчёт количества, сколько полей ссылаются на данный тип
func (p *parser) buildDependencyMap() map[xml.Name]bool {
    res := make(map[xml.Name]bool)
    for _, v := range p.types {
        for _, f := range v.Fields {
            res[f.TypeName] = true
        }
    }
    return res
}

func (p *parser) element(n *dom.Node) {
    nameAttr := n.GetAttribute("name")
    if nameAttr != nil {
        newType := p.createAndAddType(nameAttr.Value)
        p.parseStack.Push(newType)
        ct := n.GetChild("complexContent")
        ct2 := n.GetChild("simpleType")
        if ct != nil || ct2 != nil {
            f := newField("XMLValue", tString)
            newType.addField(f)
        } else {

        }
    }
}
