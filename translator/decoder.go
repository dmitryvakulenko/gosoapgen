package translator

import (
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"strings"
)

func (t *decoder) decode(s *xsd.Schema, targetNamespace string) {
	t.curTargetNamespace = targetNamespace

	if _, ok := t.namespacesList[targetNamespace]; !ok {
		t.namespacesList[targetNamespace] = true
	}

	t.parseNamespaces(s)

	for _, attrGr := range s.AttributeGroup {
		t.parseAttributeGroupTypes(attrGr)
	}

	for _, elem := range s.SimpleType {
		t.generateFromSimpleType(elem)
	}

	for _, elem := range s.Element {
		t.generateFromComplexType(elem.ComplexType, elem.Name)
	}

	for _, elem := range s.ComplexType {
		t.generateFromComplexType(elem, "")
	}

	t.resolveBaseTypes()
}


func (t *decoder) GetTypes() []interface{} {
	return t.typesList
}


func (t *decoder) addType(newType Namespaceable) {
	if !t.typesListCache.has(newType) {
		t.typesList = append(t.typesList, newType)
		t.typesListCache.put(newType)
	}
}

func (t *decoder) findAttributeGroup(fullTypeName string) (interface{}, bool) {
	ns, name := t.parseFullName(fullTypeName)
	return t.attributeGroupCache.find(ns, name)
}

func (t *decoder) findType(fullTypeName string) (interface{}, bool)  {
	ns, name := t.parseFullName(fullTypeName)
	return t.typesListCache.find(ns, name)
}

func (t *decoder) parseFullName(fullTypeName string) (string, string) {
	parts := strings.Split(fullTypeName, ":")
	if len(parts) == 2 {
		return t.curXmlns[parts[0]], parts[1]
	} else {
		return t.curTargetNamespace, parts[0]
	}
}

func (t *decoder) parseNamespaces(s *xsd.Schema) {
	t.curXmlns = make(map[string]string)
	for _, v := range s.Attrs {
		if v.Name.Space != "xmlns" {
			continue
		}
		t.curXmlns[v.Name.Local] = v.Value
	}
}

func (t *decoder) GetNamespaces() []string {
	res := make([]string, len(t.namespacesList))

	index := 0
	for ns := range t.namespacesList {
		res[index] = ns
		index++
	}

	return res
}


func (t *decoder) generateFromElement(element *xsd.Element) *Field {
	if element == nil {
		return nil
	}

	field := &Field{}
	field.Name = element.Name
	field.XmlExpr = element.Name
	field.Namespace = t.curTargetNamespace

	t.generateFromComplexType(element.ComplexType, field.Name)

	if element.Type == "" {
		field.Type = field.Name
	} else {
		field.Type = parseType(element.Type)
	}

	return field
}

func (t *decoder) generateFromAttribute(attribute *xsd.Attribute) *Field {
	field := &Field{}
	field.Name = attribute.Name
	field.XmlExpr = attribute.Name + ",attr"
	field.Type = parseType(attribute.Type)

	return field
}

func (t *decoder) generateFromComplexType(complexType *xsd.ComplexType, name string) {
	if complexType == nil {
		return
	}

	var curStruct = &ComplexType{Name: name}
	if complexType.Name != "" {
		curStruct.Name = complexType.Name
	}
	curStruct.Namespace = t.curTargetNamespace

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
		fields, baseType := t.generateFromComplexContent(complexType.ComplexContent)
		curStruct.Fields = append(curStruct.Fields, fields...)
		if baseType != "" {
			curStruct.BaseType = baseType
		}
	}
}

func (t *decoder) parseAttributes(attributes []*xsd.Attribute) []*Field {
	var res []*Field
	for _, childElem := range attributes {
		res = append(res, t.generateFromAttribute(childElem))
	}

	return res
}


func (t *decoder) parseAttributeGroupsRef(attributeGroups []*xsd.AttributeGroup) []*Field {
	var res []*Field
	for _, curGroup := range attributeGroups {
		groupType, _ := t.findAttributeGroup(curGroup.Ref)
		res = append(res, groupType.(*attributeGroup).Fields...)
	}

	return res
}

func (t *decoder) generateFromSequence(sequence *xsd.Sequence) []*Field {
	var res []*Field
	for _, childElem := range sequence.Element {
		field := t.generateFromElement(childElem)
		if field != nil {
			res = append(res, field)
		}
	}

	return res
}

func (t *decoder) generateFromSimpleContent(simpleContent *xsd.Content) []*Field {
	var res []*Field

	valField := &Field{
		Name: "Value",
		Namespace: t.curTargetNamespace,
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

func (t *decoder) generateFromComplexContent(complexContent *xsd.Content) ([]*Field, string) {
	var (
		res []*Field
		baseType string
	)

	if complexContent.Extension != nil {
		baseType, ok := t.findType(complexContent.Extension.Base)
		if !ok {
			baseType = complexContent.Extension.Base
		} else {
			res = append(res, baseType.(*ComplexType).Fields...)
		}

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

	return res, baseType
}

func (t *decoder) parseAttributeGroupTypes(attrGr *xsd.AttributeGroup) {
	curType := &attributeGroup{Name: attrGr.Name, Namespace: t.curTargetNamespace}

	for _, attr := range attrGr.Attribute {
		field := t.generateFromAttribute(attr)
		curType.Fields = append(curType.Fields, field)
	}

	t.attributeGroupCache.put(curType)
}

func (t *decoder) generateFromSimpleType(simpleType *xsd.SimpleType) {
	curType := &SimpleType{
		Name:      simpleType.Name,
		Type:      parseType(simpleType.Restriction.Base),
		Namespace: t.curTargetNamespace}
	t.addType(curType)
}


func (t *decoder) resolveBaseTypes() {
	for _, tp := range t.typesList {
		cType, ok := tp.(*ComplexType)
		if !ok || cType.BaseType == "" {
			continue
		}

		t.resolveBaseTypesImpl(cType)
	}
}

func (t *decoder) resolveBaseTypesImpl(cType *ComplexType) {
	bType, ok := t.findType(cType.BaseType)
	if !ok {
		panic("Type " + cType.BaseType + " not found")
	}

	baseType := bType.(*ComplexType)
	if baseType.BaseType != "" {
		t.resolveBaseTypesImpl(baseType)
	}
	cType.Fields = append(baseType.Fields, cType.Fields...)
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
