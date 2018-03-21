package wsdl

import "encoding/xml"

type SoapAddress struct {
	Location string `xml:"location,attr"`
}

type Definitions struct {
	XMLName         xml.Name
	TargetNamespace string      `xml:"targetNamespace,attr"`
	Namespace       []*xml.Attr  `xml:",any,attr"`
	Import          []*Type      `xml:"types>schema>import"`
	Message         []*Message   `xml:"message"`
	PortType        *PortType  `xml:"portType"`
	Binding         *Binding   `xml:"binding"`
	SoapAddress     SoapAddress `xml:"service>port>address"`
}
