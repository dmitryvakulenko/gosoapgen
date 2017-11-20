package wsdl

import "encoding/xml"

type PortType struct {
	Name xml.Attr `xml:"name,attr"`
	Operation []Operation `xml:"operation"`
}
