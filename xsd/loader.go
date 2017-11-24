package xsd

import (
	"os"
	"io/ioutil"
	"encoding/xml"
)

type Sequence struct {
	Element []Element `xml:"element"`
}

type Attribute struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Element struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	SimpleType []SimpleType `xml:"simpleType"`
	ComplexType []ComplexType `xml:"complexType"`
}

type Restriction struct {
	Base string `xml:"base,attr"`
}

type SimpleType struct {
	Name string `xml:"name,attr"`
	Restriction Restriction `xml:"restriction"`
}

type ComplexType struct {
	Name string `xml:"name,attr"`
	Sequence Sequence `xml:"sequence"`
	Attribute []Attribute `xml:"attribute"`
}

type Schema struct {
	XMLName xml.Name `xml:"schema"`
	Xmlns []xml.Attr `xml:",any,attr"`
	SimpleType []SimpleType `xml:"simpleType"`
	ComplexType []ComplexType `xml:"complexType"`
	Element []Element `xml:"element"`
}

func LoadSchema(fileName string) (*Schema, error) {
	xmlFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	content, _ := ioutil.ReadAll(xmlFile)
	schema := Schema{}
	xml.Unmarshal(content, &schema)

	return &schema, nil
}