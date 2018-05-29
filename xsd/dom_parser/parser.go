package dom_parser

import (
	dom "github.com/subchen/go-xmldom"
)

type parser struct {
	types []*Type
	parseStack *typesStack
}

func NewParser() *parser {
	return &parser{
		parseStack: new(typesStack)}
}

func (p *parser) LoadFile(fileName string) {
	doc, err := dom.ParseFile(fileName)
	if err != nil {
		panic(err)
	}

	for _, n := range doc.Root.Children {
		p.parseNode(n)
	}
}

func (p *parser) parseNode(n *dom.Node) {
	switch n.Name {
	//case "simpleType":
	//	p.simpleType(n)
	case "element":
		p.element(n)
	}
}

func (p *parser) simpleType(n *dom.Node) {
	nameAttr := n.GetAttribute("name")
	if nameAttr != nil {
		newType := p.createAndAddType(nameAttr.Value)
		newType.IsSimpleContent = true

		//typeAttr  := n.GetAttribute("type")
		//newType.BaseTypeName = p.createAndAddType()
	}
}

func (p *parser) createAndAddType(name string) *Type {
	newType := &Type{Name: name}
	p.types = append(p.types, newType)

	return newType
}

//func (p *parser) parseAndCreateQName(name string) QName {
//
//}

func (p *parser) GetTypes() []*Type {
	return p.types
}

func (p *parser) element(n *dom.Node) {
	nameAttr := n.GetAttribute("name")
	if nameAttr != nil {
		newType := p.createAndAddType(nameAttr.Value)
		p.parseStack.Push(newType)
		if ct := n.GetChild("complexContent"); ct != nil {
			newType.IsSimpleContent = false
		} else {
			newType.IsSimpleContent = true
		}
	}
}
