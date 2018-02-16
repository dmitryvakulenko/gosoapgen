package xsd

import (
	"encoding/xml"
	"io"
)

func CreateParser(reader io.Reader) Parser {
	return Parser{
		make(WsdlTypes, 0),
		nil,
		nil,
		xml.NewDecoder(reader)}
}

type Parser struct {
	types         WsdlTypes
	currentStruct *Struct
	currentField  *StructField
	decoder       *xml.Decoder
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
		//e := Element{}
		//decoder.DecodeElement(&e, elem)
		if p.currentStruct == nil {
			p.currentStruct = newStruct(getAttribute(elem.Attr, "name"))
			p.currentStruct.Namespace = elem.Name.Space
		} else {
			p.currentField = newField(getAttribute(elem.Attr, "name"), getAttribute(elem.Attr, "type"))
			p.currentField.Namespace = elem.Name.Space
		}
		//case "complexType":
		//	c := ComplexType{}
		//	decoder.DecodeElement(c, elem)
	}
}

func (p *Parser) closeElements(elem *xml.EndElement) {
	switch elem.Name.Local {
	case "element":
		if p.currentField != nil {
			p.currentStruct.appendField(p.currentField)
			p.currentField = nil
		} else {
			p.types = append(p.types, p.currentStruct)
			p.currentStruct = nil
		}
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
