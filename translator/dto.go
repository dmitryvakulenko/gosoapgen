package translator

type SchemaTypes struct {
	cType []*ComplexType
	sType []*ComplexType
}

type ComplexType struct {
	Name   string
	Type   string
	Embed  []string
	Fields []*Field
}

type SimpleType struct {
	Name string
	Type string
}

type Field struct {
	Name    string
	Type    string
	XmlExpr string
	Comment string
}
