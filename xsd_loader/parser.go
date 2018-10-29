package xsd_loader

import "github.com/dmitryvakulenko/gosoapgen/xsd_loader/tree"

func NewParser(l *tree.Loader) *Parser {
	return &Parser{loader: l}
}

type Parser struct {
	loader    *tree.Loader
	curSchema *Schema
}

func (p *Parser) Parse(fileName string) *Schema {
	p.curSchema = &Schema{}
	rawXml := p.loader.Load(fileName)
	p.schema(rawXml)
	p.linkTypes()
	return p.curSchema
}

func (p *Parser) schema(schema *tree.Schema) {
	for _, e := range schema.Children() {
		switch e.Name() {
		case "element":
			p.element(e)
		case "simpleType":
			p.xsdType(e)
		}
	}
}

func (p *Parser) element(node *tree.Node) {
	e := &Element{}
	e.Name = node.AttributeValue("name")
	e.typeName = node.AttributeValue("type")

	p.curSchema.addElement(e)
}

func (p *Parser) xsdType(node *tree.Node) {
	t := &Type{}
	t.Name = node.AttributeValue("name")
	p.curSchema.addType(t)
}

func (p *Parser) linkTypes() {
	for _, e := range p.curSchema.Elements {
		if e.typeName != "" {
			e.Type = p.findType(e.typeName)
		}
	}
}

func (p *Parser) findType(name string) *Type {
	for _, t := range p.curSchema.Types {
		if t.Name == name {
			return t
		}
	}

	panic("Can't find type " + name)
}