package xsd

type StructAttribute struct {
	Name string
}

type Type struct {
	Name string
	Comment string
	TypeName string
	Namespace string
	Fields []*Type
	Attributes []*Attribute
}

func (s *Type) appendField(f *Type) {
	s.Fields = append(s.Fields, f)
}

type WsdlTypes []*Type