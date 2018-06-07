package xsd_model

import (
    "encoding/xml"
    "io"
    "container/list"
)

func newNode(e *xml.StartElement) *node {
    return &node{
        name:      e.Name.Local,
        startElem: e}
}



func newSchema(e *xml.StartElement) *schema {
    s := &schema{
        node: *newNode(e),
        nsAlias: make(map[string]string)}

    s.TargetNamespace = s.AttributeValue("targetNamespace")
    for _, a := range s.AllAttributesByName("xmlns") {
        if a.Name.Local != "" {
            s.nsAlias[a.Name.Local] = a.Value
        }
    }

    return s
}

func Load(r io.ReadCloser) *schema {
    defer r.Close()
    decoder := xml.NewDecoder(r)

    nodes := list.New()
    var rootNode *schema
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
            p := nodes.Back()
            if p == nil {
                rootNode = newSchema(&elem)
                elemNode = rootNode
            } else {
                n := newNode(&elem)
                p.Value.(parent).addChild(n)
                elemNode = n
            }
            nodes.PushBack(elemNode)
        case xml.EndElement:
            nodes.Remove(nodes.Back())
        }
    }

    return rootNode
}
