package xsd

type StructField struct {
	Name string
	Comment string
	TypeName string
}

type StructAttribute struct {
	Name string
}

type Struct struct {
	Name string
	Comment string
	Namespace string
	Fields []*StructField
	Attributes []*Attribute
}

type WsdlTypes []*Struct

func newStruct(name string) *Struct {
	ret := &Struct{Name: name}
	ret.Attributes = make([]*Attribute, 0)
	ret.Fields = make([]*StructField, 0)

	return ret
}


func newField(name, typeName string) *StructField {
	ret := &StructField{Name: name, TypeName: typeName}
	return ret
}