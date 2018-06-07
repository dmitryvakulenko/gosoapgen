package xsd_model

import (
    "encoding/xml"
)

type parent interface {
    addChild(*node)
}

type node struct {
    name      string
    startElem *xml.StartElement
    children  []*node
}

func (n *node) addChild(e *node) {
    n.children = append(n.children, e)
}

func (n *node) Children() []*node {
    return n.children
}

func (n *node) Attribute(name string) *xml.Attr {
    for _, a := range n.startElem.Attr {
        if a.Name.Local == name {
            return &a
        }
    }

    return nil
}

func (n *node) AllAttributesByName(name string) []*xml.Attr {
    var ret []*xml.Attr
    for _, a := range n.startElem.Attr {
        if a.Name.Space == name {
            ret = append(ret, &a)
        }
    }

    return ret
}

func (n *node) AttributeValue(name string) string {
    a := n.Attribute(name)
    if a != nil {
        return a.Value
    } else {
        return ""
    }
}

func (n *node) ElementsByName(name string) []*node {
    var res []*node
    for _, v := range n.children {
        if v.name == name {
            res = append(res, v)
        }
    }

    return res
}


type Schema struct {
    node
    TargetNamespace string
    nsAlias map[string]string
}


func (s *Schema) ResolveSpace(prefix string) string {
    return s.nsAlias[prefix]
}

