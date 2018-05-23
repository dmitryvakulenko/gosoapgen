package dom_parser

//import (
//	dom "github.com/subchen/go-xmldom"
//)

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


type Field struct {
	Name     string
	Type     *Type
	TypeName *QName
	IsArray  bool
	IsAttr   bool
	Comment  string
}

//func newField(n *dom.Node) *Field {
	//return &Field{
	//	Name:     n.name,
	//	TypeName: n.typeName,
	//	IsAttr:   n.isAttr,
	//	IsArray:  n.isArray}
//}

type QName struct {
	Name      string
	Namespace string
}
