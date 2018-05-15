package tree_parser

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
