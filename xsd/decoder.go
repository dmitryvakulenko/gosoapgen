package xsd

import (
	"strings"
	"strconv"
	"github.com/dmitryvakulenko/gosoapgen/xsd/type"
)

func (t *Decoder) Decode(schemaFileName string) {
	schemas := t.schemaParser.Parse(schemaFileName, "")

	for _, s := range schemas {
		t.curTargetNamespace = s.TargetNamespace
		if _, ok := t.namespacesList[t.curTargetNamespace]; !ok {
			t.namespacesList[t.curTargetNamespace] = true
		}

		t.parseNamespaces(s)
		for _, attrGr := range s.AttributeGroup {
			t.parseAttributeGroupTypes(attrGr)
		}

		for _, elem := range s.SimpleType {
			t.generateFromSimpleType(elem, "")
		}

		for _, elem := range s.Element {
			t.generateFromElement(elem, true)
		}

		for _, elem := range s.ComplexType {
			t.parseComplexType(elem, "")
		}

		t.resolveTypes()
		t.resolveBaseTypes()
		t.prepareGoNames()
	}
}


func (t *Decoder) GetTypes() []NamedType {
	return t.typesList
}

func (t *Decoder) addType(newType NamedType) {
	if !t.typesListCache.has(t.curTargetNamespace, newType.GetName()) {
		t.typesListCache.put(t.curTargetNamespace, newType)
		t.typesList = append(t.typesList, newType)
	}
}

func (t *Decoder) findAttributeGroup(fullTypeName string) (interface{}, bool) {
	qName := t.parseFullName(fullTypeName)
	return t.attributeGroupCache.find(qName.Namespace, qName.Name)
}

func (t *Decoder) findType(fullTypeName string) (interface{}, bool) {
	qName := t.parseFullName(fullTypeName)
	return t.typesListCache.find(qName.Namespace, qName.Name)
}

func (t *Decoder) parseFullName(fullTypeName string) *QName {
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

func (t *Decoder) parseNamespaces(s *_type.Schema) {
	t.curXmlns = make(map[string]string)
	for _, v := range s.Attrs {
		if v.Name.Space != "xmlns" {
			continue
		}
		t.curXmlns[v.Name.Local] = v.Value
	}
}

func (t *Decoder) GetNamespaces() []string {
	res := make([]string, len(t.namespacesList))

	index := 0
	for ns := range t.namespacesList {
		res[index] = ns
		index++
	}

	return res
}

// Если передали fieldName - это означает, что этот элемент - поле
func (t *Decoder) generateFromElement(element *_type.Element, isParentSchema bool) *Field {
	if element == nil || element.MaxOccurs == "0" {
		return nil
	}

	typeName := strings.Title(element.Name)
	if element.SimpleType != nil {
		t.generateFromSimpleType(element.SimpleType, typeName)
	} else if element.ComplexType != nil {
		t.parseComplexType(element.ComplexType, typeName)
	}

	field := &Field{Name: element.Name}
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
	} else {
		field.TypeName = &QName{Name: typeName, Namespace: t.curTargetNamespace}
	}

	if isParentSchema {
		newType := &SimpleType{Name: field.Name, BaseTypeName: field.TypeName}
		t.addType(newType)
		return nil
	} else {
		return field
	}
}

func (t *Decoder) generateFromAttribute(attribute *_type.Attribute) *Field {
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

func (t *Decoder) parseComplexType(complexType *_type.ComplexType, fieldName string) {
	if complexType == nil {
		return
	}

	var name string
	if complexType.Name != "" {
		name = complexType.Name
	} else {
		name = fieldName
	}
	var curStruct = &ComplexType{Name: name, Namespace: t.curTargetNamespace}

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

	t.addType(curStruct)
}

func (t *Decoder) parseAttributes(attributes []*_type.Attribute) []*Field {
	var res []*Field
	for _, childElem := range attributes {
		res = append(res, t.generateFromAttribute(childElem))
	}

	return res
}

func (t *Decoder) parseAttributeGroupsRef(attributeGroups []*_type.AttributeGroup) []*Field {
	var res []*Field
	for _, curGroup := range attributeGroups {
		groupType, _ := t.findAttributeGroup(curGroup.Ref)
		res = append(res, groupType.(*attributeGroup).Fields...)
	}

	return res
}

func (t *Decoder) generateFromSequence(sequence *_type.Sequence) []*Field {
	var res []*Field
	for _, childElem := range sequence.Element {
		field := t.generateFromElement(childElem, false)
		if field != nil {
			res = append(res, field)
		}
	}

	return res
}

func (t *Decoder) generateFromSimpleContent(simpleContent *_type.Content) []*Field {
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

func (t *Decoder) generateFromComplexContent(complexContent *_type.Content, baseTypeName string) ([]*Field, string) {
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

func (t *Decoder) parseAttributeGroupTypes(attrGr *_type.AttributeGroup) {
	curType := &attributeGroup{Name: attrGr.Name, Namespace: t.curTargetNamespace}

	for _, attr := range attrGr.Attribute {
		field := t.generateFromAttribute(attr)
		curType.Fields = append(curType.Fields, field)
	}

	t.attributeGroupCache.put(t.curTargetNamespace, curType)
}

func (t *Decoder) generateFromSimpleType(simpleType *_type.SimpleType, fieldName string) {
	if simpleType == nil {
		return
	}

	typeType := simpleType.Restriction.Base
	if typeType == "" && (simpleType.Union != nil || simpleType.List != nil) {
		typeType = "string"
	}

	curType := &SimpleType{
		BaseTypeName: t.parseFullName(typeType)}

	curType.Name = simpleType.Name
	if curType.Name == "" {
		curType.Name = fieldName
	}
	t.addType(curType)
}

func (t *Decoder) resolveTypes() {
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

func (t *Decoder) resolve(typeName *QName) NamedType {
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

func (t *Decoder) resolveBaseTypes() {
	for _, cType := range t.typesList {
		realType, ok := cType.(*ComplexType)
		if !ok {
			continue
		}

		t.resolveBaseTypesImpl(realType)
	}
}

func (t *Decoder) resolveBaseTypesImpl(cType *ComplexType) {
	if cType.BaseType == nil {
		return
	}

	baseType := cType.BaseType.(*ComplexType)
	t.resolveBaseTypesImpl(baseType)
	cType.BaseType = nil

	cType.Fields = append(baseType.Fields, cType.Fields...)
}

func (t *Decoder) prepareGoNames() {
	usedNames := make(map[string]bool)
	for _, cType := range t.typesList {
		switch realType := cType.(type) {
		case *SimpleType:
			baseName := realType.Name
			realName := baseName
			index := 1
			for _, ok := usedNames[realName]; ok; index++ {
				realName = baseName + "_" + strconv.Itoa(index)
				ok = usedNames[realName]
			}
			realType.GoName = strings.Title(realName)
			usedNames[realName] = true
		case *ComplexType:
			baseName := realType.Name
			realName := baseName
			index := 1
			for _, ok := usedNames[realName]; ok; index++ {
				realName = baseName + "_" + strconv.Itoa(index)
				ok = usedNames[realName]
			}
			realType.GoName = strings.Title(realName)
			usedNames[realName] = true
		}
	}
}

//func (t *Decoder) resolveTypeImpl(qName QName) string {
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
