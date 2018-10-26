package tree

import (
	"container/list"
	"encoding/xml"
	"io"
)

type Loader struct {
	resolver Resolver
}

func NewLoader(r Resolver) *Loader {
	return &Loader{resolver: r}
}

func (l *Loader) Load(filename string) *Schema {
	return l.loadImpl(filename, "")
}

func (l *Loader) loadImpl(filename, targetNs string) *Schema {
	r, err := l.resolver.Load(filename)
    defer r.Close()
	if l.resolver.IsAlreadyLoadedError(err) {
		return nil
	}

    decoder := xml.NewDecoder(r)

    nodes := list.New()
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
            p := nodes.Back()
            if p == nil {
                rootNode = newSchema(&elem, targetNs)
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

	l.loadIncludeImports(rootNode)

    return rootNode
}


func (l *Loader) loadIncludeImports(s *Schema) {
	for _, i := range s.ChildrenByName("include") {
		inc := l.loadImpl(i.AttributeValue("schemaLocation"), s.TargetNamespace)
		if inc != nil {
			s.ChildSchemas = append(s.ChildSchemas, inc)
		}
	}

	for _, i := range s.ChildrenByName("import") {
		inc := p.loadSchema(i.AttributeValue("schemaLocation"), "")
		if inc != nil {
			s.ChildSchemas = append(s.ChildSchemas, inc)
		}
	}
}