package tree

type ComplexType struct {
	Name      string
	Namespace string
	Attributes []*Attribute
	Fields    []*Field
	Parent    string
}

func NewComplexType(name, ns string) *ComplexType {
	return &ComplexType{
		Name:      name,
		Namespace: ns}
}

type Field struct {
	Name      string
	Type      string
	TypeName  string
	Namespace string
	Comment   string
}

type Field struct {
	Name      string
	Type      string
}

func NewField(name, fieldType, expr string) *Field {
	return &Field{
		Name:     name,
		TypeName: fieldType,
		XmlExpr:  expr}
}
