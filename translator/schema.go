package translator

import (
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"strings"
)

func Parse(s *xsd.Schema) *SchemaTypes {
	res := CreateSchemaTypes()
	res.targetNamespace = s.TargetNamespace

	for _, attrGr := range s.AttributeGroup {
		res.parseAttributeGroup(attrGr)
	}

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
	field.Name = element.Name
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
	field.Name = attribute.Name
	field.XmlExpr = attribute.Name + ",attr"
	field.Type = parseStandardTypes(attribute.Type)

	return field
}

func (t *SchemaTypes) generateFromComplexType(complexType *xsd.ComplexType, name string) {
	if complexType == nil {
		return
	}

	var curStruct = &ComplexType{Name: name}
	if complexType.Name != "" {
		curStruct.Name = complexType.Name
	}
	curStruct.Namespace = t.targetNamespace

	t.cType.put(curStruct)

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

	//for _, gr := range complexType.AttributeGroup {
		//group := attrGroups[gr.Ref]
		//curStruct.Fields = append(curStruct.Fields, group.Fields...)
	//}

	if complexType.SimpleContent != nil {

	}
}

func (t *SchemaTypes) generateFromSimpleContent(simpleContent *xsd.Content) {

}

func (t *SchemaTypes) parseAttributeGroup(attrGr *xsd.AttributeGroup) {
	curType := &attributeGroup{Name: attrGr.Name, Namespace: t.targetNamespace}

	for _, attr := range attrGr.Attribute {
		field := generateFromAttribute(attr)
		curType.Fields = append(curType.Fields, field)
	}

	t.attributeGroup.put(curType)
}

func (t *SchemaTypes) generateFromSimpleType(simpleType *xsd.SimpleType) {
	curType := &SimpleType{
		Name: simpleType.Name,
		Type: parseStandardTypes(simpleType.Restriction.Base),
		Namespace: t.targetNamespace}
	t.sType.put(curType)
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
