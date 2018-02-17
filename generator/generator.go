package generator

import (
	"gosoapgen/xsd"
	"strings"
)

func GenerateTypes(s []xsd.Schema) []*Struct {
	var res []*Struct

	for _, v := range s {
		res = append(res, generateFromSchema(&v)...)
	}

	return res
}

func generateFromSchema(s *xsd.Schema) []*Struct {
	var resTypes []*Struct

	for _, elem := range s.Element {
		newTypes := generateFromComplexType(&elem.ComplexType, elem.Name)
		resTypes = append(resTypes, newTypes...)
	}

	for _, elem := range s.ComplexType {
		newTypes := generateFromComplexType(&elem, "")
		resTypes = append(resTypes, newTypes...)
	}

	return resTypes
}

func generateFromElement(element *xsd.Element) (*Field, []*Struct) {
	var resTypes []*Struct

	field := &Field{}
	field.Name = strings.ToUpper(element.Name[0:1]) + element.Name[1:]
	field.XmlExpr = element.Name

	newTypes := generateFromComplexType(&element.ComplexType, field.Name)

	if element.Type == "" {
		field.Type = field.Name
	} else {
		field.Type = parseStandardTypes(element.Type)
	}

	resTypes = append(resTypes, newTypes...)

	return field, resTypes
}

func generateFromAttribute(attribute *xsd.Attribute) *Field {
	field := &Field{}
	field.Name = strings.ToUpper(attribute.Name[0:1]) + attribute.Name[1:]
	field.XmlExpr = attribute.Name + ",attr"
	field.Type = parseStandardTypes(attribute.Type)

	return field
}

// Первое возвращаемое значение - текущий тип, второе - подтипы
func generateFromComplexType(complexType *xsd.ComplexType, name string) []*Struct {
	if len(complexType.Sequence.Element) == 0 && len(complexType.Attribute) == 0 {
		return []*Struct{}
	}

	var (
		resTypes  []*Struct
		curStruct = &Struct{Name: name}
	)

	resTypes = append(resTypes, curStruct)

	if complexType.Name != "" {
		curStruct.Name = complexType.Name
	}

	for _, childElem := range complexType.Sequence.Element {
		field, types := generateFromElement(&childElem)
		resTypes = append(resTypes, types...)
		curStruct.Fields = append(curStruct.Fields, field)
	}

	for _, childElem := range complexType.Attribute {
		field := generateFromAttribute(&childElem)
		curStruct.Fields = append(curStruct.Fields, field)
	}

	return resTypes
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
