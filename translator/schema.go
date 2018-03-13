package translator

import (
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"strings"
)

func Parse(s *xsd.Schema, targetNamespace string) *SchemaTypes {
	res := newSchemaTypes()
	res.targetNamespace = targetNamespace

	res.fillNamespaces(s)

	for _, attrGr := range s.AttributeGroup {
		res.parseAttributeGroupTypes(attrGr)
	}

	for _, elem := range s.SimpleType {
		res.generateFromSimpleType(elem)
	}

	for _, elem := range s.Element {
		res.generateFromComplexType(elem.ComplexType, elem.Name)
	}

	for _, elem := range s.ComplexType {
		res.generateFromComplexType(elem, "")
	}

	return &res
}

func (t *SchemaTypes) GetTypes() []interface{} {
	return t.typesList
}


func (t *SchemaTypes) addType(newType Namespaceable) {
	t.typesList = append(t.typesList, newType)
	t.typesListCache.put(newType)
}

func (t *SchemaTypes) findAttributeGroup(fullTypeName string) (interface{}, bool) {
	ns, name := t.parseFullName(fullTypeName)
	return t.attributeGroupCache.find(ns, name)
}

func (t *SchemaTypes) findType(fullTypeName string) (interface{}, bool)  {
	ns, name := t.parseFullName(fullTypeName)
	return t.typesListCache.find(ns, name)
}

func (t *SchemaTypes) parseFullName(fullTypeName string) (string, string) {
	parts := strings.Split(fullTypeName, ":")
	if len(parts) == 2 {
		return t.curXmlns[parts[0]], parts[1]
	} else {
		return t.targetNamespace, parts[0]
	}
}

func (t *SchemaTypes) fillNamespaces(s *xsd.Schema) {
	t.curXmlns = make(map[string]string)
	for _, v := range s.Attrs {
		if v.Name.Space != "xmlns" {
			continue
		}
		t.curXmlns[v.Name.Local] = v.Value
	}
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
		field.Type = parseType(element.Type)
	}

	return field
}

func (t *SchemaTypes) generateFromAttribute(attribute *xsd.Attribute) *Field {
	field := &Field{}
	field.Name = attribute.Name
	field.XmlExpr = attribute.Name + ",attr"
	field.Type = parseType(attribute.Type)

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

	t.addType(curStruct)

	if complexType.Sequence != nil {
		curStruct.Fields = append(curStruct.Fields, t.generateFromSequence(complexType.Sequence)...)
	}

	curStruct.Fields = append(curStruct.Fields, t.parseAttributes(complexType.Attribute)...)
	curStruct.Fields = append(curStruct.Fields, t.parseAttributeGroupsRef(complexType.AttributeGroup)...)

	if complexType.SimpleContent != nil {
		curStruct.Fields = append(curStruct.Fields, t.generateFromSimpleContent(complexType.SimpleContent)...)
	}

	if complexType.ComplexContent != nil {
		curStruct.Fields = append(curStruct.Fields, t.generateFromComplexContent(complexType.ComplexContent)...)
	}
}

func (t *SchemaTypes) parseAttributes(attributes []*xsd.Attribute) []*Field {
	var res []*Field
	for _, childElem := range attributes {
		res = append(res, t.generateFromAttribute(childElem))
	}

	return res
}


func (t *SchemaTypes) parseAttributeGroupsRef(attributeGroups []*xsd.AttributeGroup) []*Field {
	var res []*Field
	for _, curGroup := range attributeGroups {
		groupType, _ := t.findAttributeGroup(curGroup.Ref)
		res = append(res, groupType.(*attributeGroup).Fields...)
	}

	return res
}

func (t *SchemaTypes) generateFromSequence(sequence *xsd.Sequence) []*Field {
	var res []*Field
	for _, childElem := range sequence.Element {
		field := t.generateFromElement(childElem)
		if field != nil {
			res = append(res, field)
		}
	}

	return res
}

func (t *SchemaTypes) generateFromSimpleContent(simpleContent *xsd.Content) []*Field {
	var res []*Field

	valField := &Field{
		Name: "Value",
		Namespace: t.targetNamespace,
		XmlExpr: ",chardata"}

	res = append(res, valField)

	if simpleContent.Extension != nil {
		valField.Type = parseType(simpleContent.Extension.Base)
		for _, v := range simpleContent.Extension.Attribute {
			res = append(res, t.generateFromAttribute(v))
		}

		for _, v := range simpleContent.Extension.AttributeGroup {
			groupI, ok := t.findAttributeGroup(v.Ref)
			if !ok {
				panic("No attribute group " + v.Ref + " found")
			}
			res = append(res, groupI.(*attributeGroup).Fields...)
		}
	}

	if simpleContent.Restriction != nil {
		valField.Type = parseType(simpleContent.Restriction.Base)
		// парсить ограничения смысла нет
	}

	return res
}

func (t *SchemaTypes) generateFromComplexContent(complexContent *xsd.Content) []*Field {
	var res []*Field

	if complexContent.Extension != nil {
		baseType, ok := t.findType(complexContent.Extension.Base)
		if !ok {
			panic("No type " + complexContent.Extension.Base + " found")
		}

		res = append(res, baseType.(*ComplexType).Fields...)
		if complexContent.Extension.Sequence != nil {
			res = append(res, t.generateFromSequence(complexContent.Extension.Sequence)...)
		}

		for _, attr := range complexContent.Extension.Attribute {
			res = append(res, t.generateFromAttribute(attr))
		}
	}

	if complexContent.Restriction != nil {
		// это обрабатывать смысла нет, т.к. там ограничиваются только значения полей, но не сами поля
	}

	return res
}

func (t *SchemaTypes) parseAttributeGroupTypes(attrGr *xsd.AttributeGroup) {
	curType := &attributeGroup{Name: attrGr.Name, Namespace: t.targetNamespace}

	for _, attr := range attrGr.Attribute {
		field := t.generateFromAttribute(attr)
		curType.Fields = append(curType.Fields, field)
	}

	t.attributeGroupCache.put(curType)
}

func (t *SchemaTypes) generateFromSimpleType(simpleType *xsd.SimpleType) {
	curType := &SimpleType{
		Name:      simpleType.Name,
		Type:      parseType(simpleType.Restriction.Base),
		Namespace: t.targetNamespace}
	t.addType(curType)
}


func parseType(xmlType string) string {
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
