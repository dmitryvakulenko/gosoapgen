package tree_parser

import (
    "encoding/xml"
    xsd "github.com/dmitryvakulenko/gosoapgen/xsd-model"
    "strconv"
)

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
