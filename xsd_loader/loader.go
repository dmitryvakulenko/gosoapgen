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
			l.xsdType(e)
		}
	}
}

func (l *Loader) element(node *tree.Node) {
	e := &Element{}
	e.Name = l.schemaDeep.buildFullName(node.AttributeValue("name"))
	e.typeName = l.schemaDeep.buildFullName(node.AttributeValue("type"))

	l.curSchema.addElement(e)
}

func (l *Loader) xsdType(node *tree.Node) {
	t := &Type{}
	t.Name = l.schemaDeep.buildFullName(node.AttributeValue("name"))
	if r := node.ChildByName("restriction"); r != nil {
		l.restriction(t, r)
	} else if r := node.ChildByName("extension"); r != nil {

	}

	l.curSchema.addType(t)
}

func (l *Loader) restriction(t *Type, node *tree.Node) {
	t.baseTypeName = l.schemaDeep.buildFullName(node.AttributeValue("base"))
}

func (l *Loader) linkTypes() {
	for _, e := range l.curSchema.Elements {
		if e.typeName.Local != "" {
			e.Type = l.findType(e.typeName)
		}
	}

	for _, t := range l.curSchema.Types {
		if t.baseTypeName.Local != "" {
			if t.baseTypeName.Space == xsdSpace {
				t.BaseType = &Type{Name: t.baseTypeName}
			} else {
				t.BaseType = l.findType(t.baseTypeName)
			}
		}
	}
}

func (l *Loader) findType(name xml.Name) *Type {
	for _, t := range l.curSchema.Types {
		if t.Name == name {
			return t
		}
	}

	panic("Can't find type " + name.Local)
}