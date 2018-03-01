package translator

import (
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"strings"
)

var types []*ComplexType


func Parse(s *xsd.Schema) *SchemaTypes {
	res := SchemaTypes{}

	for _, elem := range s.Element {
		generateFromComplexType(elem.ComplexType, elem.Name)
	}

	for _, attrGr := range s.AttributeGroup {
		generateFromAttributeGroup(attrGr)
	}

	for _, elem := range s.ComplexType {
		generateFromComplexType(elem, "")
	}

	for _, elem := range s.SimpleType {
		res.generateFromSimpleType(elem)
	}

	return &res
}

func generateFromElement(element *xsd.Element) *Field {
	if element == nil {
		return nil
	}

	field := &Field{}
	field.Name = strings.ToUpper(element.Name[0:1]) + element.Name[1:]
	field.XmlExpr = element.Name

	generateFromComplexType(element.ComplexType, field.Name)

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

func generateFromComplexType(complexType *xsd.ComplexType, name string) {
	if complexType == nil {
		return
	}

	if complexType.Sequence == nil && len(complexType.Attribute) == 0 && len(complexType.AttributeGroup) == 0 {
		return
	}

	var curStruct = &ComplexType{Name: name}
	types = append(types, curStruct)

	if complexType.Name != "" {
		curStruct.Name = complexType.Name
	}

	if complexType.Sequence != nil {
		for _, childElem := range complexType.Sequence.Element {
			field := generateFromElement(childElem)
			if field != nil {
				curStruct.Fields = append(curStruct.Fields, field)
			}
		}
	}

	for _, childElem := range complexType.Attribute {
		field := generateFromAttribute(childElem)
		curStruct.Fields = append(curStruct.Fields, field)
	}

	for _, attrGr := range complexType.AttributeGroup {
		curStruct.Embed = append(curStruct.Embed, attrGr.Ref)
	}
}
func generateFromAttributeGroup(attrGr *xsd.AttributeGroup) *ComplexType {
	curType := &ComplexType{Name: attrGr.Name}
	types = append(types, curType)

	for _, attr := range attrGr.Attribute {
		field := generateFromAttribute(attr)
		curType.Fields = append(curType.Fields, field)
	}

	for _, inAttrGr := range attrGr.AttributeGroup {
		curType.Embed = append(curType.Embed, inAttrGr.Ref)
	}

	return curType
}

func (t *SchemaTypes) generateFromSimpleType(simpleType *xsd.SimpleType) {
	curType := &ComplexType{Name: simpleType.Name, Type: parseStandardTypes(simpleType.Restriction.Base)}
	t.sType = append(t.sType, curType)
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
