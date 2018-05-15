/*
Парсер схемы xsd.
Использует ручной обход всего дерева, вместо загрузки его с помощь Unmarshal.
 */
package tree_parser

import (
	"io"
	"encoding/xml"
	"fmt"
)

// Интерфейс загрузки xsd
// должен отслеживать уже загруженные файлы
// и правильно отрабатывать относительные пути
type Loader interface {
	/*
	Загрузить файл по указанному пути (или url)
	Второй параметр - ошибка, которую должен уметь анализировать метод IsAlreadyLoadedError
	 */
	Load(path string) (io.ReadCloser, error)
	IsAlreadyLoadedError(error) bool
}

type parser struct {
	loader     Loader
	typesStack *typesStack
	nsStack    *stringsStack
	typesList  *NamespacedTypes
}

func NewParser(l Loader) *parser {
	return &parser{
		loader:     l,
		typesStack: &typesStack{},
		nsStack:    &stringsStack{},
		typesList:  NewTypesCollection()}
}

func (p *parser) Parse(inputFile string) {
	reader, _ := p.loader.Load(inputFile)
	decoder := xml.NewDecoder(reader)
	p.parseImpl(decoder)
}

func (p *parser) parseImpl(decoder *xml.Decoder) {
	for token, err := decoder.Token(); err != io.EOF; token, err = decoder.Token() {
		switch token.(type) {
		case xml.StartElement:
			elem := token.(xml.StartElement)
			p.parseStartElement(&elem)
		case xml.EndElement:
			elem := token.(xml.EndElement)
			p.parseEndElement(&elem)
		}
	}
}

func (p *parser) parseStartElement(elem *xml.StartElement) {
	switch elem.Name.Local {
	case "schema":
		ns := findAttributeByName(elem.Attr, "targetNamespace")
		if ns != nil {
			p.nsStack.Push(ns.Value)
		} else {
			// Используем родительский ns. А правильно ли?
			p.nsStack.Push(p.nsStack.GetLast())
		}
	case "simpleType":
		p.parseSimpleType(elem)
	}
}

func (p *parser) parseEndElement(elem *xml.EndElement) {
	fmt.Printf("End %q\n", elem.Name.Local)
}

func (p *parser) parseSimpleType(elem *xml.StartElement) {
	curType := &SimpleType{}
	p.typeStarted(curType)
}

func (p *parser) typeStarted(t NamedType) {
	p.typesStack.Push(t)
	p.typesList.Put(p.nsStack.GetLast(), t)
}

func (p *parser) GetTypes() []NamedType {
	return p.typesList.GetAllTypes()
}


func findAttributeByName(attrsList []xml.Attr, name string) *xml.Attr {
	for _, attr := range attrsList {
		if attr.Name.Local == name {
			return &attr
		}
	}

	return nil
}