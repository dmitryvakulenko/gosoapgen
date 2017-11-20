package wsdl

import "encoding/xml"

type Part struct {
	Element xml.Attr `xml:"element,attr"`
	Name xml.Attr `xml:"name,attr"`
}

type Message struct {
	XMLName xml.Name `xml:"message"`
	Name string `xml:"name,attr"`
	Part Part `xml:"part"`
}
