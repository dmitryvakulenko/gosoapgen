package wsdl

import "encoding/xml"

type Header struct {
	Message xml.Attr `xml:"message,attr"`
	Use xml.Attr `xml:"use,attr"`
	Part xml.Attr `xml:"part,attr"`
}

type Param struct {
	XMLName xml.Name
	Message string `xml:"message,attr"`
	Header []Header `xml:"header"`
}

type SoapOperation struct {
	SoapAction string `xml:"soapAction,attr"`
}

type Operation struct {
	Name string `xml:"name,attr"`
	Input Param `xml:"input"`
	Output Param `xml:"output"`
	SoapOperation SoapOperation `xml:"operation"`
}
