package xsd

import "encoding/xml"

type Sequence struct {
	Element []Element `xml:"element"`
}

type Attribute struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Element struct {
	Name        string       `xml:"name,attr"`
	Type        string       `xml:"type,attr"`
	Ref         string       `xml:"ref,attr"`
	SimpleType  []*SimpleType `xml:"simpleType"`
	ComplexType *ComplexType  `xml:"complexType"`
}

type Restriction struct {
	Base string `xml:"base,attr"`
}

type SimpleType struct {
	Name        string      `xml:"name,attr"`
	Restriction Restriction `xml:"restriction"`
}

type ComplexType struct {
	Name      string      `xml:"name,attr"`
	Sequence  Sequence    `xml:"sequence"`
	Attribute []Attribute `xml:"attribute"`
}

type Schema struct {
	XMLName         xml.Name      `xml:"schema"`
	TargetNamespace string        `xml:"targetNamespace,attr"`
	Element         []Element     `xml:"element"`
	Attrs           []xml.Attr    `xml:",any,attr"`
	SimpleType      []*SimpleType  `xml:"simpleType"`
	ComplexType     []*ComplexType `xml:"complexType"`
	Import          []*Import      `xml:"import"`
	Include         []*Import      `xml:"include"`
}

type Import struct {
	Namespace      string `xml:"namespace,attr"`
	SchemaLocation string `xml:"schemaLocation,attr"`
}
