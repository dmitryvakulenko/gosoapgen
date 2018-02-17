package generator

import "gosoapgen/xsd"

func GenerateTypes(v interface{}) []*Struct {
	res := []*Struct{}
	var newTypes []*Struct
	switch v.(type) {
	case xsd.Schema:
		tmp := v.(xsd.Schema)
		newTypes = generateFromSchema(&tmp)
	}
	res = append(res, newTypes...)

	return res
}

func generateFromSchema(s *xsd.Schema) []*Struct {
	var res []*Struct

	for _, elem := range s.Element {
		newStruct := &Struct{}
		newStruct.Name = elem.Name
		res = append(res, newStruct)
	}

	return res
}