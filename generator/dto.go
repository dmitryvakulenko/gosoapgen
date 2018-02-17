package generator

type Struct struct {
	Name string
	Fields []*Field
}

type Field struct {
	Name string
	Type string
	XmlExpr string
}