package tree

import "encoding/xml"

type Element struct {
	Name    xml.Name `xml:",any"`
	Element []*Element
	Attr    []*xml.Attr `xml:",any,attr"`
}
