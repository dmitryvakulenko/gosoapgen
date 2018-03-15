package tree

import (
	"os"
	"encoding/xml"
	"fmt"
	"io"
)

type Builder struct {
	typesList	[]*ComplexType
}

func NewBuilder() Builder {
	return Builder{}
}

func (b *Builder) getTypes() []*ComplexType {
	return b.typesList
}

func (b *Builder) Build(uri string) {
	reader := makeReader(uri)
	defer reader.Close()

	decoder := xml.NewDecoder(reader)
	token, err := decoder.Token()
	for err != io.EOF {
		switch t := token.(type) {
		case xml.StartElement:
			fmt.Printf("Name %q\n", t.Name)
		}
		token, err = decoder.Token()
	}
}




func makeReader(uri string) io.ReadCloser {
	reader, err := os.Open(uri)

	if err != nil {
		panic(err)
	}

	return reader
}