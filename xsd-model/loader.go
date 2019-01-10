package xsd_model

import (
    "encoding/xml"
    "io"
    "container/list"
)

func newNode(e *xml.StartElement, parent interface{}) *Node {
	node := &Node{
		name:      e.Name.Local,
		startElem: e}

	if p, ok := parent.(*Node); ok {
		node.parent = p
	}

    return node
}



func newSchema(e *xml.StartElement) *Schema {
    s := &Schema{
        Node:    *newNode(e, nil),
        nsAlias: make(map[string]string)}

    s.TargetNamespace = s.AttributeValue("targetNamespace")
    for _, a := range s.AllAttributesByName("xmlns") {
        if a.Name.Local != "" {
            s.nsAlias[a.Name.Local] = a.Value
        }
    }

    return s
}

func Load(r io.ReadCloser) *Schema {
    defer r.Close()
    decoder := xml.NewDecoder(r)

    nodesStack := list.New()
    var rootNode *Schema
    for {
        token, err := decoder.Token()
        if err == io.EOF {
            break
        }

        if err != nil {
            panic(err)
        }

        switch elem := token.(type) {
        case xml.StartElement:
            var elemNode parent
            p := nodesStack.Back()
            if p == nil {
                rootNode = newSchema(&elem)
                elemNode = rootNode
            } else {
                n := newNode(&elem, p.Value)
                p.Value.(parent).addChild(n)
                elemNode = n
            }
            nodesStack.PushBack(elemNode)
        case xml.EndElement:
            nodesStack.Remove(nodesStack.Back())
        }
    }

    return rootNode
}
