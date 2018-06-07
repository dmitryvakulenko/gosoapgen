package xml

import (
    "encoding/xml"
)

type Node struct {
    ElementName string
    StartElem   *xml.StartElement
    Children    []*Node
}

func (n *Node) addChild(e *Node) {
    n.Children = append(n.Children, e)
}