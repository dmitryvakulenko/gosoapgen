package translator

import (
	"gosoapgen/xsd"
	"strings"
)

var types []*Struct

func GenerateTypes(s []xsd.Schema) []*Struct {
	types = []*Struct{}

	for _, v := range s {
		generateFromSchema(&v)
	}

	return types
}

func generateFromSchema(s *xsd.Schema) {
	for _, elem := range s.Element {
		generateFromComplexType(&elem.ComplexType, elem.Name)
	}

	for _, elem := range s.ComplexType {
		generateFromComplexType(&elem, "")
	}

	for _, elem := range s.SimpleType {
		generateFromSimpleType(&elem)
	}
}

func generateFromElement(element *xsd.Element) *Field {
	field := &Field{}
	field.Name = strings.ToUpper(element.Name[0:1]) + element.Name[1:]
	field.XmlExpr = element.Name

	generateFromComplexType(&element.ComplexType, field.Name)

	if element.Type == "" {
		field.Type = field.Name
	} else {
		field.Type = parseStandardTypes(element.Type)
	}

	return field
}

func generateFromAttribute(attribute *xsd.Attribute) *Field {
	field := &Field{}
	field.Name = strings.ToUpper(attribute.Name[0:1]) + attribute.Name[1:]
	field.XmlExpr = attribute.Name + ",attr"
	field.Type = parseStandardTypes(attribute.Type)

	return field
}

// Первое возвращаемое значение - текущий тип, второе - подтипы
func generateFromComplexType(complexType *xsd.ComplexType, name string) {
	if len(complexType.Sequence.Element) == 0 && len(complexType.Attribute) == 0 {
		return
	}

	var curStruct = &Struct{Name: name}
	types = append(types, curStruct)

	if complexType.Name != "" {
		curStruct.Name = complexType.Name
	}

	for _, childElem := range complexType.Sequence.Element {
		field := generateFromElement(&childElem)
		curStruct.Fields = append(curStruct.Fields, field)
	}

	for _, childElem := range complexType.Attribute {
		field := generateFromAttribute(&childElem)
		curStruct.Fields = append(curStruct.Fields, field)
	}
}

func generateFromSimpleType(simpleType *xsd.SimpleType) {
	curType := &Struct{Name: simpleType.Name, Type: parseStandardTypes(simpleType.Restriction.Base)}
	types = append(types, curType)
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
