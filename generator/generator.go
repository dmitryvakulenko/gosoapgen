package generator

import (
	"gosoapgen/xsd"
	"strings"
)

func GenerateTypes(s []xsd.Schema) []*Struct {
	var res []*Struct

	for _, v := range s {
		_, newTypes := generateTypesImpl(v)
		res = append(res, newTypes...)
	}

	return res
}

func generateTypesImpl(v interface{}) (*Field, []*Struct) {
	var (
		resTypes, newTypes []*Struct
		newField           *Field
	)
	switch v.(type) {
	case xsd.Schema:
		tmp := v.(xsd.Schema)
		newTypes = generateFromSchema(&tmp)
	case xsd.Element:
		tmp := v.(xsd.Element)
		newField, newTypes = generateFromElement(&tmp)
	case xsd.Attribute:
		tmp := v.(xsd.Attribute)
		newField = generateFromAttribute(&tmp)
	}

	resTypes = append(resTypes, newTypes...)

	return newField, resTypes
}


func generateFromSchema(s *xsd.Schema) []*Struct {
	var resTypes []*Struct

	for _, elem := range s.Element {
		newStruct := &Struct{}
		resTypes = append(resTypes, newStruct)

		newStruct.Name = elem.Name

		for _, childElem := range elem.ComplexType.Sequence.Element {
			field, types := generateTypesImpl(childElem)
			resTypes = append(resTypes, types...)
			newStruct.Fields = append(newStruct.Fields, field)
		}

		for _, childElem := range elem.ComplexType.Attribute {
			field, types := generateTypesImpl(childElem)
			resTypes = append(resTypes, types...)
			newStruct.Fields = append(newStruct.Fields, field)
		}

	}

	return resTypes
}

func generateFromElement(element *xsd.Element) (*Field, []*Struct) {
	field := &Field{}
	field.Name = strings.ToUpper(element.Name[0:1]) + element.Name[1:]
	field.XmlExpr = element.Name
	field.Type = parseStandardTypes(element.Type)

	return field, []*Struct{}
}

func generateFromAttribute(attribute *xsd.Attribute) *Field {
	field := &Field{}
	field.Name = strings.ToUpper(attribute.Name[0:1]) + attribute.Name[1:]
	field.XmlExpr = attribute.Name + ",attr"
	field.Type = parseStandardTypes(attribute.Type)

	return field
}


func parseStandardTypes(xmlType string) string {
	parts := strings.Split(xmlType, ":")

	var fieldType string
	if len(parts) == 2 {
		fieldType = parts[1]
	} else {
		fieldType = parts[0]
	}

	switch fieldType {
	case "integer":
		return "int"
	case "decimal":
		return "float64"
	case "boolean":
		return "bool"
	case "date":
		return "time.time"
	default:
		return fieldType
	}
}