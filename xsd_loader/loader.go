package xsd_loader

import (
	"encoding/xml"
	"github.com/dmitryvakulenko/gosoapgen/xsd_loader/tree"
)

const xsdSpace = "http://www.w3.org/2001/XMLSchema"

func NewLoader(l *tree.Loader) *Loader {
	return &Loader{loader: l}
}

type Loader struct {
	loader    *tree.Loader
	curSchema *Schema

	schemaDeep schemaDeep
}

func (l *Loader) Load(fileName string) *Schema {
	l.curSchema = &Schema{}
	rawXml := l.loader.Load(fileName)
	l.schema(rawXml)
	l.linkTypes()
	return l.curSchema
}

func (l *Loader) schema(schema *tree.Schema) {
	l.schemaDeep.push(schema)
	defer l.schemaDeep.pop()

	for _, ch := range schema.ChildSchemas {
		l.schema(ch)
	}

	for _, e := range schema.Children() {
		switch e.Name() {
		case "element":
			l.element(e)
		case "simpleType":
			l.simpleType(e)
		}
	}
}

func (l *Loader) element(node *tree.Node) {
	e := &Element{}
	e.Name = l.schemaDeep.buildFullName(node.AttributeValue("name"))
	e.typeName = l.schemaDeep.buildFullName(node.AttributeValue("type"))

	if e.typeName.Local == "" {
		st := node.ChildByName("simpleType")
		ct := node.ChildByName("complexType")
		if st != nil {
			e.Type = l.simpleType(st)
		} else if ct != nil {
			e.Type = l.complexType(ct)
		}
	}

	l.curSchema.addElement(e)
}


func (l *Loader) simpleType(node *tree.Node) *SimpleType {
	t := &SimpleType{}
	t.Name = l.schemaDeep.buildFullName(node.AttributeValue("name"))
	if r := node.ChildByName("restriction"); r != nil {
		l.simpleTypeRestriction(t, r)
	} else if r := node.ChildByName("extension"); r != nil {

	}

	if t.Name.Local != "" {
		l.curSchema.addType(t)
		return nil
	} else {
		return t
	}
}

func (l *Loader) complexType(node *tree.Node) *ComplexType {
	return nil
}

func (l *Loader) simpleTypeRestriction(t *SimpleType, node *tree.Node) {
	t.baseTypeName = l.schemaDeep.buildFullName(node.AttributeValue("base"))
}

func (l *Loader) linkTypes() {
	for _, e := range l.curSchema.Elements {
		if e.typeName.Local != "" {
			e.Type = l.findType(e.typeName)
		}

		if e.Type != nil {
			l.linkBaseTypes(e.Type)
		}
	}

	for _, t := range l.curSchema.Types {
		l.linkBaseTypes(t)
	}
}

func (l *Loader) linkBaseTypes(t Type) {
	switch tp := t.(type) {
	case *SimpleType:
		base := l.findType(tp.baseTypeName).(*SimpleType)
		if base.baseTypeName.Local != "" && base.BaseType == nil {
			l.linkBaseTypes(base)
		}
		tp.BaseType = base

	case *ComplexType:
	}
}

func (l *Loader) findType(name xml.Name) Type {
	if name.Local == "" {
		return nil
	}

	if name.Space == xsdSpace {
		return &SimpleType{Name: name}
	}

	for _, t := range l.curSchema.Types {
		if t.GetName() == name {
			return t
		}
	}

	panic("Can't find type " + name.Local)
}