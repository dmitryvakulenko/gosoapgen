package tree_parser

import "encoding/xml"

// Абстрактное представление элемента схемы
type node struct {
	// название элемента схемы xsd
	// из которого был создан данный тип
	// node, complexType и т.д.
	// чисто для сокрачения, поскольку вся эта информация содержится в startElem
	elemName  string
	// Имя типа. По сути, значение атрибута name
	name      string
	typeName  *QName
	startElem *xml.StartElement
	namespace string
	children  []*node
	// список типов встраиваемых элементов
	refs []string
	isSimple      bool
	isAttr        bool
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

//type node struct {
//	// название элемента схемы xsd
//	elementName string
//	// название самого типа - не обязательно
//	typeName string
//	// пространство имён
//	namespace string
//	// дочерние узлы
//	children []*node
//	// это атрибут?
//	isAttr bool
//}
//
//func (n *node) Find(ns, elemName string) *node {
//	if n.namespace == ns && n.typeName == elemName {
//		return n
//	}
//
//	for _, ch := range n.children {
//		res := ch.Find(ns, elemName)
//		if res != nil {
//			return res
//		}
//	}
//
//	return nil
//}

func newNode(startElem *xml.StartElement) *node {
	name := ""
	for _, a := range startElem.Attr {
		if a.Name.Local == "name" {
			name = a.Value
			break
		}
	}

	return &node{
		name: name,
		elemName:  startElem.Name.Local,
		startElem: startElem}
}

type Type struct {
	Name         string
	IsSimple     bool
	Namespace    string
	GoName       string
	Fields       []*Field

	// Только для simpleType
	BaseType     *Type
	BaseTypeName *QName
}

func (t *Type) addField(f *Field) {
	t.Fields = append(t.Fields, f)
}

func newType(n * node) *Type {
	return &Type{
		Name: n.name,
		Namespace: n.namespace,
		IsSimple: n.isSimple}
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

func newField() *Field {
	return &Field{}
}

type attributeGroup struct {
	Name      string
	Namespace string
	Fields    []*Field
}

func (t *attributeGroup) GetName() string {
	return t.Name
}

func (t *attributeGroup) GetGoName() string {
	return t.Name
}

type QName struct {
	Name      string
	Namespace string
}
