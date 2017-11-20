package wsdl

import "encoding/xml"

type Type struct {
	XMLName xml.Name `xml:"import"`
	SchemaLocation string `xml:"schemaLocation,attr"`
	Namespace string `xml:"namespace,attr"`
}