package xml

import (
    "encoding/xml"
    "io"
    "container/list"
)

func newNode(e *xml.StartElement) *Node {
    return &Node{
        name:      e.Name.Local,
        startElem: e}
}

func Load(r io.ReadCloser) *Node {
    defer r.Close()
    decoder := xml.NewDecoder(r)

    nodes := list.New()
    var rootNode *Node
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
            newNode := newNode(&elem)
            parent := nodes.Back()
            if parent != nil {
                parent.Value.(*Node).addChild(newNode)
            } else {
                rootNode = newNode
            }
            nodes.PushBack(newNode)
        case xml.EndElement:
            nodes.Remove(nodes.Back())
        }
    }

    return rootNode
}
