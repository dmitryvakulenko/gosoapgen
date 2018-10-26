package tree

import (
    "encoding/xml"
)

type parent interface {
    addChild(*Node)
}

type Node struct {
    name      string
    startElem *xml.StartElement
    children  []*Node
}

func (n *Node) Name() string {
    return n.name
}

func (n *Node) addChild(e *Node) {
    n.children = append(n.children, e)
}

func (n *Node) Children() []*Node {
    return n.children
}

func (n *Node) Attribute(name string) *xml.Attr {
    for _, a := range n.startElem.Attr {
        if a.Name.Local == name {
            return &a
        }
    }

    return nil
}

func (n *Node) AllAttributesByName(name string) []xml.Attr {
    var ret []xml.Attr
    for _, a := range n.startElem.Attr {
        if a.Name.Space == name {
            ret = append(ret, a)
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

func (n *Node) ChildByName(name string) *Node {
    for _, v := range n.children {
        if v.name == name {
            return v
        }
    }

    return nil
}

func (n *Node) ChildrenByName(name string) []*Node {
    var res []*Node
    for _, v := range n.children {
        if v.name == name {
            res = append(res, v)
        }
    }

    return res
}

type Schema struct {
    Node
    TargetNamespace string
    ChildSchemas    []*Schema
    nsAlias         map[string]string
}

func (s *Schema) ResolveSpace(prefix string) string {
    return s.nsAlias[prefix]
}

func (s *Schema) FindRootType(typeName xml.Name) *Node {
    if s.TargetNamespace == typeName.Space {
        for _, n := range s.children {
            if n.AttributeValue("name") == typeName.Local {
                return n
            }
        }
    }

    for _, sc := range s.ChildSchemas {
        n := sc.FindRootType(typeName)
        if n != nil {
            return n
        }
    }

    return nil
}

func newNode(e *xml.StartElement) *Node {
	return &Node{
		name:      e.Name.Local,
		startElem: e}
}

func newSchema(e *xml.StartElement, targetNamespace string) *Schema {
	s := &Schema{
		Node:    *newNode(e),
		nsAlias: make(map[string]string)}

	if targetNamespace != "" {
		s.TargetNamespace = targetNamespace
	} else {
		s.TargetNamespace = s.AttributeValue("targetNamespace")
	}

	for _, a := range s.AllAttributesByName("xmlns") {
		if a.Name.Local != "" {
			s.nsAlias[a.Name.Local] = a.Value
		}
	}

	return s
}
