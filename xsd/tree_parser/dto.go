package tree_parser

import "encoding/xml"

// Абстрактное представление элемента схемы
type node struct {
    // название элемента схемы xsd
    // из которого был создан данный тип
    // node, complexType и т.д.
    // чисто для сокрачения, поскольку вся эта информация содержится в startElem
    elemName string
    // Имя типа. По сути, значение атрибута name
    name            string
    typeName        xml.Name
    startElem       *xml.StartElement
    namespace       string
    children        []*node
    isSimpleContent bool
    isAttr          bool
    isArray         bool

    // сгенерированный тип
    genType *Type
}

func (r *node) find(ns, name string) *node {
    for _, n := range r.children {
        if ns == n.namespace && name == n.name {
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
        name:      name,
        elemName:  startElem.Name.Local,
        startElem: startElem}
}

type Type struct {
    xml.Name
    Fields []*Field

    // Если это не simpleType, значит создано из extension/restriction
    BaseType     *Type
    BaseTypeName xml.Name
}

func (t *Type) addField(f *Field) {
    t.Fields = append(t.Fields, f)
}

func (t *Type) Hash() {
    // надо реализовать и учитывать при проверке дублей
}

func newType(n *node) *Type {
    return &Type{
        Name:            n.name,
        Namespace:       n.namespace}
}

type Field struct {
    Name     string
    Type     *Type
    TypeName xml.Name
    IsArray  bool
    IsAttr   bool
    Comment  string
}

func newField(n *node) *Field {
    return &Field{
        Name:     n.name,
        TypeName: n.typeName,
        IsAttr:   n.isAttr,
        IsArray:  n.isArray}
}

