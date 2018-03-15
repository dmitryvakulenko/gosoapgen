package translator

import (
	"os"
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"encoding/xml"
	"path"
)

type Parser struct {
	decoder decoder
}

func NewParser() Parser {
	return Parser{
		decoder: newDecoder()}
}

func (p *Parser) Parse(fileName string) {
	p.parseImpl(fileName, "")
}

func (p *Parser) GetTypes() []*ComplexType {
	return p.decoder.GetTypes()
}

func (p *Parser) parseImpl(fileName, ns string) {
	s := loadSchema(fileName)

	baseDir := path.Dir(fileName) + "/"
	for _, imp := range s.Import {
		p.parseImpl(path.Clean(baseDir + imp.SchemaLocation), "")
	}

	for _, imp := range s.Include {
		p.parseImpl(path.Clean(baseDir + imp.SchemaLocation), s.TargetNamespace)
	}

	if ns != "" {
		p.decoder.decode(s, ns)
	} else {
		p.decoder.decode(s, s.TargetNamespace)
	}
}

func loadSchema(fileName string) *xsd.Schema {
	reader, err := os.Open(fileName)
	defer reader.Close()

	if err != nil {
		panic(err)
	}

	s := xsd.Schema{}
	err = xml.NewDecoder(reader).Decode(&s)
	if err != nil {
		panic(err)
	}

	return &s
}
