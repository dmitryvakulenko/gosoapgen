package tree_parser

import (
    "encoding/xml"
    xsd "github.com/dmitryvakulenko/gosoapgen/xsd-model"
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

func newNode(startElem *xml.StartElement) *node {
    name := ""
    for _, a := range startElem.Attr {
        if a.Name.Local == "name" {
            name = a.Value
            break
        }
    }

    return &node{
        name:      xml.Name{Local: name},
        elemName:  startElem.Name.Local,
        startElem: startElem}
}

type Type struct {
    xml.Name
    Fields     []*Field
    SourceNode *xsd.Node

    // Only for simple types
    BaseType     *Type
    BaseTypeName xml.Name
}

func (t *Type) addField(f *Field) {
    t.Fields = append(t.Fields, f)
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

type Field struct {
    Name     string
    Type     *Type
    TypeName xml.Name
    IsArray  bool
    IsAttr   bool
    Comment  string
}

func newField(n *xsd.Node, typ xml.Name) *Field {
    return &Field{
        Name:     n.AttributeValue("name"),
        TypeName: typ}
        // IsAttr:   n.isAttr,
        // IsArray:  n.isArray}
}
