package tree_parser

type NamedType interface {
	GetName() string
	GetGoName() string
}

type ComplexType struct {
	Name         string
	Namespace    string
	GoName       string
	Fields       []*Field
	BaseType     NamedType
	BaseTypeName *QName
}

func (t *ComplexType) GetName() string {
	return t.Name
}

func (t *ComplexType) GetGoName() string {
	if t.GoName != "" {
		return t.GoName
	} else {
		return t.Name
	}
}

type SimpleType struct {
	Name         string
	GoName       string
	BaseType     NamedType
	BaseTypeName *QName
}

func (t *SimpleType) GetName() string {
	return t.Name
}

func (t *SimpleType) GetGoName() string {
	if t.GoName != "" {
		return t.GoName
	} else {
		return t.Name
	}
}

type Field struct {
	Name      string
	Type      NamedType
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
