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
		res.parseAttributeGroup(attrGr)
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
	allTypes := t.typesList.getAllTypes()
	res := make([]interface{}, 0)

	for _, v := range allTypes {
		if _, ok := v.(*xsd.AttributeGroup); ok {
			continue
		}
		res = append(res, v)
	}

	return res
}


func (t *SchemaTypes) addType(newType Namespaceable) {
	t.typesList.put(newType)
}

func (t *SchemaTypes) findAttributeGroup(fullTypeName string) (interface{}, bool) {
	ns, name := t.parseFullName(fullTypeName)
	return t.attributeGroup.find(ns, name)
}

func (t *SchemaTypes) findType(fullTypeName string) (interface{}, bool)  {
	ns, name := t.parseFullName(fullTypeName)
	return t.typesList.find(ns, name)
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
		for _, childElem := range complexType.Sequence.Element {
			field := t.generateFromElement(childElem)
			if field != nil {
				curStruct.Fields = append(curStruct.Fields, field)
			}
		}
	}

	for _, childElem := range complexType.Attribute {
		field := t.generateFromAttribute(childElem)
		curStruct.Fields = append(curStruct.Fields, field)
	}

	for _, gr := range complexType.AttributeGroup {
		group, _ := t.findAttributeGroup(gr.Ref)
		curStruct.Fields = append(curStruct.Fields, group.(*attributeGroup).Fields...)
	}

	if complexType.SimpleContent != nil {
		curStruct.Fields = append(curStruct.Fields, t.generateFromSimpleContent(complexType.SimpleContent)...)
	}
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
				panic("No attribute group " + v.Ref + " found", )
			}
			res = append(res, groupI.(*attributeGroup).Fields...)
		}
	}

	if simpleContent.Restriction != nil {
		valField.Type = parseType(simpleContent.Restriction.Base)
	}

	return res
}

func (t *SchemaTypes) parseAttributeGroup(attrGr *xsd.AttributeGroup) {
	curType := &attributeGroup{Name: attrGr.Name, Namespace: t.targetNamespace}

	for _, attr := range attrGr.Attribute {
		field := t.generateFromAttribute(attr)
		curType.Fields = append(curType.Fields, field)
	}

	t.attributeGroup.put(curType)
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
