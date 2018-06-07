package xml

import (
    "encoding/xml"
)

type Node struct {
    name      string
    startElem *xml.StartElement
    children  []*Node
}

func (n *Node) addChild(e *Node) {
    n.children = append(n.children, e)
}

func (n *Node) Children() []*Node {
    return n.children
}

func (n *Node) FirstChild() *Node {
    if len(n.children) == 0 {
        return nil
    }
    return n.children[0]
}

func (n *Node) Attribute(name string) *xml.Attr {
    for _, a := range n.startElem.Attr {
        if a.Name.Local == name {
            return &a
        }
    }

    return nil
}

func (n *Node) AllAttributesByName(name string) []*xml.Attr {
    var ret []*xml.Attr
    for _, a := range n.startElem.Attr {
        if a.Name.Local == name {
            ret = append(ret, &a)
        }
    }

    return ret
}

func (n *Node) AttributeValue(name string) string {
    a := n.Attribute(name)
    if a != nil {
        return a.Value
    } else {
        return ""
    }
}