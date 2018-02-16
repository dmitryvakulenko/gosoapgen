package xsd

import (
	"encoding/xml"
	"io"
)

const (
	STAGE_NONE = iota
	STAGE_STRUCT
	STAGE_FIELD
)

func CreateParser(reader io.Reader) Parser {
	return Parser{
		types:   make(WsdlTypes, 0),
		decoder: xml.NewDecoder(reader),
		stage: STAGE_NONE}
}

type Parser struct {
	types         WsdlTypes
	currentStruct *Type
	currentField  *Type
	lastElement   *Type
	decoder       *xml.Decoder
	stage         int
}

func (p *Parser) GetTypes() WsdlTypes {
	return p.types
}

func (p *Parser) Parse() {
	t, err := p.decoder.Token()

	if err == io.EOF {
		return
	}

	switch t := t.(type) {
	case xml.StartElement:
		p.parseElements(&t)
	case xml.EndElement:
		p.closeElements(&t)
	}

	p.Parse()
}

func (p *Parser) parseElements(elem *xml.StartElement) {
	switch elem.Name.Local {
	//case "schema":
	//	s := Schema{}
	//	decoder.DecodeElement(s, elem)
	case "element":
		p.lastElement = &Type{}
		p.lastElement.Name = getAttribute(elem.Attr, "name")
		p.lastElement.Namespace = elem.Name.Space
		p.lastElement.TypeName = getAttribute(elem.Attr, "type")

	case "complexType":
		p.currentStruct = p.lastElement
	}
}

func (p *Parser) closeElements(elem *xml.EndElement) {
	switch elem.Name.Local {
	case "element":
		if p.currentField != nil {
			p.currentStruct.appendField(p.currentField)
		}
		p.currentField = nil
	case "complexType":
		p.types = append(p.types, p.currentStruct)
		p.currentStruct = nil
	}
}

func getAttribute(attr []xml.Attr, name string) string {
	for _, v := range attr {
		if v.Name.Local == name {
			return v.Value
		}
	}

	return ""
}
