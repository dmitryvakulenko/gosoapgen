package tree

type ComplexType struct {
	Name      string
	Namespace string
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
	Namespace string
	XmlExpr   string
	Comment   string
}

func NewField(name, fieldType, expr string) *Field {
	return &Field{
		Name:    name,
		Type:    fieldType,
		XmlExpr: expr}
}
