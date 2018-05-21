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
	name      string
	typeName  *QName
	startElem *xml.StartElement
	namespace string
	children  []*node
	// список типов встраиваемых элементов
	refs            []string
	isSimpleContent bool
	isAttr          bool

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

func (r *node) add(e *node) {
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
	Name string
	// Если нет полей, то тип по определению - simple
	// Если поля есть, но установлен этот флаг - значит, поля атрибуты,
	// а сам элемент содержит простой контент
	IsSimpleContent bool
	Namespace       string
	GoName          string
	Fields          []*Field

	// Если это не simpleType, значит создано из extension/restriction
	BaseType     *Type
	BaseTypeName *QName
}

func (t *Type) addField(f *Field) {
	t.Fields = append(t.Fields, f)
}

func newType(n *node) *Type {
	return &Type{
		Name:            n.name,
		Namespace:       n.namespace,
		IsSimpleContent: n.isSimpleContent}
}

type Field struct {
	Name      string
	Type      *Type
	TypeName  *QName
	MinOccurs int
	MaxOccurs int
	IsAttr    bool
	Comment   string
}

func newField(n *node) *Field {
	return &Field{
		Name:     n.name,
		TypeName: n.typeName,
		IsAttr:   n.isAttr}
}

type QName struct {
	Name      string
	Namespace string
}
