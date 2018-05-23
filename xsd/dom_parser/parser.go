package dom_parser

import (
	dom "github.com/subchen/go-xmldom"
)

type parser struct {
}

func NewParser() *parser {
	return &parser{}
}

func (p *parser) LoadFile(fileName string) {
	_, err := dom.ParseFile(fileName)
	if err != nil {
		panic(err)
	}

}

func (p *parser) GetTypes() []*Type {
	return []*Type{}
}
