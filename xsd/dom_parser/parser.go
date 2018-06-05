package dom_parser

import (
    dom "github.com/subchen/go-xmldom"
    "encoding/xml"
    "container/list"
    "log"
    "strings"
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
    nsAlias    map[string]string
}

func NewParser() *parser {
    return &parser{
        types:      make(map[xml.Name]*Type),
        parseStack: new(typesStack),
        nsStack:    list.New(),
        nsAlias:    make(map[string]string)}
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
    targetNs := n.Root.GetAttributeValue("targetNamespace")
    if targetNs != "" {
        p.nsStack.PushBack(targetNs)
    } else {
        prevNs := p.nsStack.Back().Value.(string)
        p.nsStack.PushBack(prevNs)
    }

    for _, ns := range n.Root.Attributes {
        if ns.Name != "xmlns" {
            continue
        }
    }
}

func (p *parser) parseNode(n *dom.Node) {
    switch n.Name {
    case "simpleType":
        p.simpleType(n)
    case "element":
        p.element(n)
    }
}

func (p *parser) simpleType(n *dom.Node) {
    nameAttr := n.GetAttribute("name")
    if nameAttr != nil {
        newType := p.createAndAddType(nameAttr.Value)
        switch ch := n.FirstChild(); ch.Name {
        case "restriction":
            newType.BaseTypeName = p.parseRestriction(ch)
        }
    }
}

func (p *parser) parseRestriction(n *dom.Node) xml.Name {
    base := n.GetAttributeValue("base")
    return p.createName(base)
}

func (p *parser) createName(typ string) xml.Name {
    parts := strings.Split(typ, ":")
    if len(parts) == 2 {
        return xml.Name{Local: parts[1], Space: p.nsAlias[parts[0]]}
    } else {
        return xml.Name{Local: parts[0], Space: p.nsStack.Back().Value.(string)}
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
        ct := n.FirstChild()
        if ct.Name == "simpleType" || ct.Name == "simpleContent" {
            f := newField("XMLValue", tString)
            newType.addField(f)
        } else {

        }
    }
}
