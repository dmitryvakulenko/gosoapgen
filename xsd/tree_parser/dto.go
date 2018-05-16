package tree_parser

import "encoding/xml"

// Абстрактное представление элемента схемы
type element struct {
	// название элемента схемы xsd
	// из которого был создан данный тип
	// element, complexType и т.д.
	// чисто для сокрачения, поскольку вся эта информация содержится в startElem
	name      string
	typeName  *QName
	isSimple  bool
	startElem *xml.StartElement
	namespace string
	children  []*element
	isAttr    bool
}

func newElement(startElem *xml.StartElement) *element {
	return &element{
		name:      startElem.Name.Local,
		startElem: startElem}
}

type Type struct {
	Name         string
	IsSimple     bool
	Namespace    string
	GoName       string
	Fields       []*Field
	BaseType     *Type
	BaseTypeName *QName
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
