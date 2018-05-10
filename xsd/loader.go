package xsd

import (
	"os"
	"encoding/xml"
	"path"
	"github.com/dmitryvakulenko/gosoapgen/xsd/type"
)

type Parser struct {
	alreadyParsed map[string]bool
	decoder decoder
}

func NewLoader() Parser {
	return Parser{
		alreadyParsed: make(map[string]bool),
		decoder: newDecoder()}
}

func (p *Parser) Parse(fileName string) {
	p.parseImpl(fileName, "")
}

func (p *Parser) GetTypes() []NamedType {
	return p.decoder.GetTypes()
}

func (p *Parser) GetNamespaces() []string {
	return p.decoder.GetNamespaces()
}

func (p *Parser) parseImpl(fileName, ns string) {
	s := loadSchema(fileName)

	baseDir := path.Dir(fileName) + "/"
	for _, imp := range s.Import {
		fullName := path.Clean(baseDir + imp.SchemaLocation)
		if _, parsed := p.alreadyParsed[fullName]; !parsed {
			p.parseImpl(fullName, "")
			p.alreadyParsed[fullName] = true
		}
	}

	for _, imp := range s.Include {
		fullName := path.Clean(baseDir + imp.SchemaLocation)
		if _, parsed := p.alreadyParsed[fullName]; !parsed {
			p.parseImpl(fullName, s.TargetNamespace)
			p.alreadyParsed[fullName] = true
		}
	}

	if ns != "" {
		p.decoder.decode(s, ns)
	} else {
		p.decoder.decode(s, s.TargetNamespace)
	}
}

func loadSchema(fileName string) *_type.Schema {
	reader, err := os.Open(fileName)
	defer reader.Close()

	if err != nil {
		panic(err)
	}

	s := _type.Schema{}
	err = xml.NewDecoder(reader).Decode(&s)
	if err != nil {
		panic(err)
	}

	return &s
}
