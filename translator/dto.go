package translator

type SchemaTypes struct {
	include         []string
	targetNamespace string
	cType           []*ComplexType
	sType           []*SimpleType
}

type ComplexType struct {
	Name   string
	Type   string
	Fields []*Field
}

type SimpleType struct {
	Name      string
	Type      string
	Namespace string
}

type Field struct {
	Name      string
	Type      string
	XmlExpr   string
	Comment   string
	Namespace string
}

type attributeGroup struct {
	Fields []*Field
}
