package translator

import (
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"strings"
	"strconv"
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
		t.generateFromNamedSimpleType(elem)
	}

	for _, elem := range s.Element {
		t.generateFromElement(elem)
	}

	for _, elem := range s.ComplexType {
		newType := t.parseComplexType(elem, "")
		t.addType(newType)
	}

	t.resolveTypes()
	t.resolveBaseTypes()
	t.prepareGoNames()
}

func (t *decoder) GetTypes() []NamedType {
	return t.typesList
}

func (t *decoder) addType(newType NamedType) {
	if !t.typesListCache.has(t.curTargetNamespace, newType.GetName()) {
		t.typesListCache.put(t.curTargetNamespace, newType)
		t.typesList = append(t.typesList, newType)
	}
}

func (t *decoder) findAttributeGroup(fullTypeName string) (interface{}, bool) {
	qName := t.parseFullName(fullTypeName)
	return t.attributeGroupCache.find(qName.Namespace, qName.Name)
}

func (t *decoder) findType(fullTypeName string) (interface{}, bool) {
	qName := t.parseFullName(fullTypeName)
	return t.typesListCache.find(qName.Namespace, qName.Name)
}

func (t *decoder) parseFullName(fullTypeName string) *QName {
	parts := strings.Split(fullTypeName, ":")
	if len(parts) == 2 {
		return &QName{
			parts[1],
			t.curXmlns[parts[0]]}
	} else {
		return &QName{
			parts[0],
			t.curTargetNamespace}
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
	if element == nil || element.MaxOccurs == "0" {
		return nil
	}

	if element.SimpleType != nil {
		t.generateFromNamedSimpleType(element.SimpleType)
		return nil
	} else if element.ComplexType != nil {
		cType := t.parseComplexType(element.ComplexType, "")
		cType.Name = element.Name
		t.addType(cType)
		return nil
	} else {
		field := &Field{}
		field.Name = element.Name
		if element.MinOccurs == "0" {
			field.MinOccurs, _ = strconv.Atoi(element.MinOccurs)
		}

		if element.MaxOccurs == "" {
			field.MaxOccurs = 0
		} else {
			maxOccurs, err := strconv.Atoi(element.MaxOccurs)
			if err != nil {
				field.MaxOccurs = 100
			} else {
				field.MaxOccurs = maxOccurs
			}
		}

		if element.Type != "" {
			field.TypeName = t.parseFullName(element.Type)
		} else if element.Ref != "" {
			field.TypeName = t.parseFullName(element.Ref)
		}
		return field
	}
}

func (t *decoder) generateFromAttribute(attribute *xsd.Attribute) *Field {
	field := &Field{
		Name:   attribute.Name,
		IsAttr: true}

	if attribute.Type != "" {
		field.TypeName = t.parseFullName(attribute.Type)
	} else if attribute.SimpleType != nil {
		field.TypeName = t.parseFullName(attribute.SimpleType.Restriction.Base)
	}

	return field
}

func (t *decoder) parseComplexType(complexType *xsd.ComplexType, baseTypeName string) *ComplexType {
	if complexType == nil {
		return nil
	}

	var curStruct = &ComplexType{Name: complexType.Name, Namespace: t.curTargetNamespace}

	if complexType.Sequence != nil {
		curStruct.Fields = append(curStruct.Fields, t.generateFromSequence(complexType.Sequence)...)
	}

	curStruct.Fields = append(curStruct.Fields, t.parseAttributes(complexType.Attribute)...)
	curStruct.Fields = append(curStruct.Fields, t.parseAttributeGroupsRef(complexType.AttributeGroup)...)

	if complexType.SimpleContent != nil {
		curStruct.Fields = append(curStruct.Fields, t.generateFromSimpleContent(complexType.SimpleContent)...)
	}

	if complexType.ComplexContent != nil {
		fields, baseType := t.generateFromComplexContent(complexType.ComplexContent, curStruct.Name)
		curStruct.Fields = append(curStruct.Fields, fields...)
		if baseType != "" {
			curStruct.BaseTypeName = t.parseFullName(baseType)
		}
	}

	return curStruct
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
		Name: "Value"}

	res = append(res, valField)

	if simpleContent.Extension != nil {
		valField.TypeName = t.parseFullName(simpleContent.Extension.Base)
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
		valField.TypeName = t.parseFullName(simpleContent.Restriction.Base)
		// парсить ограничения смысла нет
	}

	return res
}

func (t *decoder) generateFromComplexContent(complexContent *xsd.Content, baseTypeName string) ([]*Field, string) {
	var (
		res      []*Field
		baseType string
	)

	if complexContent.Extension != nil {
		baseType, ok := t.findType(complexContent.Extension.Base)
		if !ok {
			baseType = complexContent.Extension.Base
		} else {
			switch tp := baseType.(type) {
			case *SimpleType:
				res = append(res, &Field{Type: tp})
			case *ComplexType:
				res = append(res, tp.Fields...)
			}
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

	t.attributeGroupCache.put(t.curTargetNamespace, curType)
}

func (t *decoder) generateFromNamedSimpleType(simpleType *xsd.SimpleType) {
	if simpleType == nil {
		return
	}

	if simpleType.Name == "" {
		panic("No name for simple type")
	}

	typeType := simpleType.Restriction.Base
	if typeType == "" && (simpleType.Union != nil || simpleType.List != nil) {
		typeType = "string"
	}

	curType := &SimpleType{
		BaseTypeName: t.parseFullName(typeType)}

	curType.Name = simpleType.Name
	t.addType(curType)
}

func (t *decoder) resolveTypes() {
	for _, curType := range t.typesList {
		switch realType := curType.(type) {
		case *ComplexType:
			if realType.BaseTypeName != nil {
				realType.BaseType = t.resolve(realType.BaseTypeName)
			}
			for _, curField := range realType.Fields {
				curField.Type = t.resolve(curField.TypeName)
			}
		case *SimpleType:
			realType.BaseType = t.resolve(realType.BaseTypeName)
		}

	}
}

func (t *decoder) resolve(typeName *QName) NamedType {
	stdType := mapStandardType(typeName.Name)
	if stdType != "" {
		return &SimpleType{Name: stdType}
	}

	curType, ok := t.typesListCache.find(typeName.Namespace, typeName.Name)
	if !ok {
		panic("Type " + typeName.Name + " not found")
	}

	return curType
}

func mapStandardType(xmlType string) string {
	switch xmlType {
	case "integer", "positiveInteger", "nonNegativeInteger":
		return "int"
	case "decimal":
		return "float64"
	case "boolean":
		return "bool"
	case "date", "dateTime":
		return "time.Time"
	case "string", "NMTOKEN", "anyURI", "language", "base64Binary", "duration":
		return "string"
	default:
		return ""
	}
}

func (t *decoder) resolveBaseTypes() {
	for _, cType := range t.typesList {
		realType, ok := cType.(*ComplexType)
		if !ok {
			continue
		}

		t.resolveBaseTypesImpl(realType)
	}
}

func (t *decoder) resolveBaseTypesImpl(cType *ComplexType) {
	if cType.BaseType == nil {
		return
	}

	baseType := cType.BaseType.(*ComplexType)
	t.resolveBaseTypesImpl(baseType)
	cType.BaseType = nil

	cType.Fields = append(baseType.Fields, cType.Fields...)
}

func (t *decoder) prepareGoNames() {
	usedNames := make(map[string]bool)
	for _, cType := range t.typesList {
		switch realType := cType.(type) {
		case *SimpleType:
			baseName := realType.Name
			realName := baseName
			index := 1
			for _, ok := usedNames[realName]; ok; index++ {
				realName += "_" + strconv.Itoa(index)
			}
			realType.GoName = strings.Title(realName)
		case *ComplexType:
			baseName := realType.Name
			realName := baseName
			index := 1
			for _, ok := usedNames[realName]; ok; index++ {
				realName += "_" + strconv.Itoa(index)
			}
			realType.GoName = strings.Title(realName)
		}
	}
}

//func (t *decoder) resolveTypeImpl(qName QName) string {
//	tmpType := mapStandardType(qName.Name)
//	if tmpType != "" {
//		return tmpType
//	}
//
//	curType, ok := t.typesListCache.find(qName.Namespace, qName.Name)
//	if !ok {
//		panic("TypeName " + qName.Name + " not found")
//	}
//
//	switch v := curType.(type) {
//	case *SimpleType:
//		return v.GoName
//	case *ComplexType:
//		return v.GoName
//	default:
//		return ""
//	}
//}
