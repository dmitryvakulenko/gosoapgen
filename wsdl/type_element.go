package wsdl

import "encoding/xml"

type Element struct {
	
}

type TypeElement struct {
	XMLName xml.Name
	Name string `xml:"name,attr"`
}
