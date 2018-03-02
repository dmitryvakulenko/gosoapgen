package translator

import (
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"strings"
)

var attrGroups map[string]*attributeGroup

func Parse(s *xsd.Schema) *SchemaTypes {
	attrGroups = make(map[string]*attributeGroup)
	for _, attrGr := range s.AttributeGroup {
		parseAttributeGroup(attrGr)
	}

	res := SchemaTypes{}
	res.targetNamespace = s.TargetNamespace

	for _, elem := range s.Element {
		res.generateFromComplexType(elem.ComplexType, elem.Name)
	}

	for _, elem := range s.ComplexType {
		res.generateFromComplexType(elem, "")
	}

	for _, elem := range s.SimpleType {
		res.generateFromSimpleType(elem)
	}

	return &res
}

func (t *SchemaTypes) generateFromElement(element *xsd.Element) *Field {
	if element == nil {
		return nil
	}

	field := &Field{}
	field.Name = strings.ToUpper(element.Name[0:1]) + element.Name[1:]
	field.XmlExpr = element.Name
	field.Namespace = t.targetNamespace

	t.generateFromComplexType(element.ComplexType, field.Name)

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

func (t *SchemaTypes) generateFromComplexType(complexType *xsd.ComplexType, name string) {
	if complexType == nil {
		return
	}

	if complexType.Sequence == nil && len(complexType.Attribute) == 0 && len(complexType.AttributeGroup) == 0 {
		return
	}

	var curStruct = &ComplexType{Name: name}
	t.cType = append(t.cType, curStruct)

	if complexType.Name != "" {
		curStruct.Name = complexType.Name
	}

	if complexType.Sequence != nil {
		for _, childElem := range complexType.Sequence.Element {
			field := t.generateFromElement(childElem)
			if field != nil {
				curStruct.Fields = append(curStruct.Fields, field)
			}
		}
	}

	for _, childElem := range complexType.Attribute {
		field := generateFromAttribute(childElem)
		curStruct.Fields = append(curStruct.Fields, field)
	}

	for _, gr := range complexType.AttributeGroup {
		group := attrGroups[gr.Ref]
		curStruct.Fields = append(curStruct.Fields, group.Fields...)
	}
}

func parseAttributeGroup(attrGr *xsd.AttributeGroup) {
	curType := &attributeGroup{}

	for _, attr := range attrGr.Attribute {
		field := generateFromAttribute(attr)
		curType.Fields = append(curType.Fields, field)
	}

	attrGroups[attrGr.Name] = curType
}

func (t *SchemaTypes) generateFromSimpleType(simpleType *xsd.SimpleType) {
	curType := &SimpleType{
		Name: simpleType.Name,
		Type: parseStandardTypes(simpleType.Restriction.Base),
		Namespace: t.targetNamespace}
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
