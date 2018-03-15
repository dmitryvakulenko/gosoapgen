package tree

type ComplexType struct {
	Name      string
	Namespace string
	Fields    []*Field
	Parent  string
}

type Field struct {
	Name      string
	Type      string
	Namespace string
	XmlExpr   string
	Comment   string
}
