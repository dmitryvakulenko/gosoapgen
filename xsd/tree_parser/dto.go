package tree_parser

import (
    "encoding/xml"
    xsd "github.com/dmitryvakulenko/gosoapgen/xsd-model"
    "strconv"
)

// Абстрактное представление элемента схемы
type node struct {
    // название элемента схемы xsd из которого был создан данный node, complexType и т.д.
    // чисто для сокращения, поскольку вся эта информация содержится в startElem
    elemName string

    // elements name, value of a "name" attribute
    name xml.Name

    // element type name
    typeName xml.Name

    // сам элемент, из которого создавался node
    startElem       *xml.StartElement
    children        []*node
    isSimpleContent bool
    isAttr          bool
    isArray         bool

    // сгенерированный тип
    genType *Type
}

func (r *node) find(ns, name string) *node {
    for _, n := range r.children {
        if ns == n.name.Space && name == n.name.Local {
            return n
        }
    }

    return nil
}

func (r *node) addChild(e *node) {
    r.children = append(r.children, e)
}

type Type struct {
    xml.Name
    Fields            []*Field
    SourceNode        *xsd.Node
    baseType          *Type
    isSimpleContent   bool
    simpleContentType *Type
    // for this type base type fields was resolved
    resolved bool
    // on this element has reference
    referenced bool
}

func (t *Type) addField(f *Field) {
    t.Fields = append(t.Fields, f)
}

func (t *Type) append(addType *Type) {
    t.Fields = append(t.Fields, addType.Fields...)
    t.isSimpleContent = addType.isSimpleContent
}

func (t *Type) Hash() {
    // надо реализовать и учитывать при проверке дублей
}

func newType(n *xsd.Node, ns string) *Type {
    name := n.AttributeValue("name")
    return &Type{
        Name:       xml.Name{Local: name, Space: ns},
        SourceNode: n}
}

func newStandardType(name string) *Type {
    return &Type{Name: xml.Name{Local: name, Space: "http://www.w3.org/2001/XMLSchema"}, isSimpleContent: true}
}

type Field struct {
    Name      string
    Type      *Type
    MinOccurs int
    MaxOccurs int
    IsAttr    bool
    Comment   string
}

func newField(n *xsd.Node, typ *Type) *Field {
    name := n.AttributeValue("name")
    if name == "" {
        name = n.AttributeValue("ref")
    }

    var min int
    switch m := n.AttributeValue("minOccurs"); m {
    case "unqualified", "":
        min = 0
    default:
        min, _ = strconv.Atoi(m)
    }

    var max int
    switch m := n.AttributeValue("maxOccurs"); m {
    case "unbounded":
        max = 1000
    case "":
        max = 0
    default:
        max, _ = strconv.Atoi(m)
    }

    return &Field{
        Name:      name,
        Type:      typ,
        MinOccurs: min,
        MaxOccurs: max}
}

func newXMLNameField() *Field {
    return &Field{
        Name: "XMLName",
        Type: newStandardType("string")}
}

func newValueField(baseType string) *Field {
    return &Field{
        Name: "XMLValue",
        Type: newStandardType(baseType)}
}
