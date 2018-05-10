package _type

import "encoding/xml"

type Sequence struct {
	Element []*Element `xml:"element"`
}

type Attribute struct {
	Name       string      `xml:"name,attr"`
	Type       string      `xml:"type,attr"`
	SimpleType *SimpleType `xml:"simpleType"`
}

type Element struct {
	Name        string       `xml:"name,attr"`
	Type        string       `xml:"type,attr"`
	Ref         string       `xml:"ref,attr"`
	SimpleType  *SimpleType  `xml:"simpleType"`
	ComplexType *ComplexType `xml:"complexType"`
	MinOccurs   string       `xml:"minOccurs,attr"`
	MaxOccurs   string       `xml:"maxOccurs,attr"`
}

type SimpleType struct {
	Name        string      `xml:"name,attr"`
	Restriction Restriction `xml:"restriction"`
	Union       *Union      `xml:"union"`
	List        *List       `xml:"list"`
}

type ComplexType struct {
	Name           string            `xml:"name,attr"`
	Sequence       *Sequence         `xml:"sequence"`
	Attribute      []*Attribute      `xml:"attribute"`
	AttributeGroup []*AttributeGroup `xml:"attributeGroup"`
	SimpleContent  *Content          `xml:"simpleContent"`
	ComplexContent *Content          `xml:"complexContent"`
}

type AttributeGroup struct {
	Name           string            `xml:"name,attr"`
	Attribute      []*Attribute      `xml:"attribute"`
	AttributeGroup []*AttributeGroup `xml:"attributeGroup"`
	Ref            string            `xml:"ref,attr"`
}

type Schema struct {
	XMLName         xml.Name          `xml:"schema"`
	TargetNamespace string            `xml:"targetNamespace,attr"`
	Element         []*Element        `xml:"element"`
	Attrs           []*xml.Attr       `xml:",any,attr"`
	SimpleType      []*SimpleType     `xml:"simpleType"`
	ComplexType     []*ComplexType    `xml:"complexType"`
	Import          []*Import         `xml:"import"`
	Include         []*Import         `xml:"include"`
	AttributeGroup  []*AttributeGroup `xml:"attributeGroup"`
}

type Import struct {
	Namespace      string `xml:"namespace,attr"`
	SchemaLocation string `xml:"schemaLocation,attr"`
}

type Content struct {
	Restriction *Restriction `xml:"restriction"`
	Extension   *Extension   `xml:"extension"`
}

type Restriction struct {
	Base           string            `xml:"base,attr"`
	Attribute      []*Attribute      `xml:"attribute"`
	AttributeGroup []*AttributeGroup `xml:"attributeGroup"`
}

type Extension struct {
	Base           string            `xml:"base,attr"`
	Sequence       *Sequence         `xml:"sequence"`
	Attribute      []*Attribute      `xml:"attribute"`
	AttributeGroup []*AttributeGroup `xml:"attributeGroup"`
}

type Union struct {
	MemberTypes string `xml:"memberTypes,attr"`
}

type List struct {
	ItemType string `xml:"itemType,attr"`
}
