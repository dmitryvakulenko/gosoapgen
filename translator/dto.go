package translator

type Struct struct {
	Name string
	// признак того, что это не структура, а SimpleType
	// идея не очень хорошая, но пока так
	Type string
	Embed []string
	Fields []*Field
}

type Field struct {
	Name string
	Type string
	XmlExpr string
}