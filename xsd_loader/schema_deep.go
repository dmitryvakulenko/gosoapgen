package xsd_loader

import (
	"encoding/xml"
	"github.com/dmitryvakulenko/gosoapgen/xsd_loader/tree"
)

type schemaDeep []*tree.Schema

func (d *schemaDeep) push(s *tree.Schema) {
	*d = append(*d, s)
}

func (d *schemaDeep) pop() {
	*d = (*d)[:len(*d)-1]
}

func (d *schemaDeep) buildFullName(n string) xml.Name {
	l := len(*d)
	if l == 0 {
		panic("List is empty")
	}

	return (*d)[l - 1].BuildFullName(n)
}