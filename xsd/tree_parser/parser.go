/*
Парсер схемы xsd.
Использует ручной обход всего дерева, вместо загрузки его с помощь Unmarshal.
 */
package tree_parser

import (
	"io"
	"encoding/xml"
	"strings"
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
	curNs      map[string]string

	// результирующий список типов
	typesList *NamespacedTypes
}

func NewParser(l Loader) *parser {
	return &parser{
		loader:     l,
		typesStack: &typesStack{},
		nsStack:    &stringsStack{},
		curNs:      make(map[string]string),
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
		p.parseSchema(elem)
	case "simpleType":
		p.parseSimpleType(elem)
	case "restriction":
		p.parseRestriction(elem)
	case "element":
		p.parseElement(elem)
	}
}

func (p *parser) parseEndElement(elem *xml.EndElement) {
	switch elem.Name.Local {
	case "schema":
		p.nsStack.Pop()
	case "simpleType":
		p.typesStack.Pop()
	}
}

func (p *parser) parseSchema(elem *xml.StartElement) {
	ns := findAttributeByName(elem.Attr, "targetNamespace")
	if ns != nil {
		p.nsStack.Push(ns.Value)
	} else {
		// Используем родительский ns. А правильно ли?
		p.nsStack.Push(p.nsStack.GetLast())
	}

	p.curNs = make(map[string]string)
	for _, attr := range elem.Attr {
		if attr.Name.Space == "xmlns" && attr.Name.Local != "" {
			p.curNs[attr.Name.Local] = attr.Value
		}
	}
}

func (p *parser) parseSimpleType(elem *xml.StartElement) {
	curType := &SimpleType{}
	curType.Name = findAttributeByName(elem.Attr, "name").Value
	p.typeStarted(curType)
}

func (p *parser) typeStarted(t NamedType) {
	p.typesStack.Push(t)
	p.typesList.Put(p.nsStack.GetLast(), t)
}

func (p *parser) GetTypes() []NamedType {
	return p.typesList.GetAllTypes()
}

func (p *parser) parseElement(element *xml.StartElement) {
	//nameElem := findAttributeByName(element.Attr, "name")

}

func (p *parser) parseRestriction(element *xml.StartElement) {
	c := p.typesStack.GetLast()
	switch context := c.(type) {
	case *SimpleType:
		baseTypeName := findAttributeByName(element.Attr, "base")
		context.BaseTypeName = p.createQName(baseTypeName.Value)
	}
}

func (p *parser) createQName(name string) *QName {
	typesParts := strings.Split(name, ":")
	if len(typesParts) != 2 {
		panic("Can't parse " + name)
	}
	ns, ok := p.curNs[typesParts[0]]
	if !ok {
		panic("Unknown namespace alias " + typesParts[0])
	}

	return &QName{
		Name:      typesParts[1],
		Namespace: ns}
}

func findAttributeByName(attrsList []xml.Attr, name string) *xml.Attr {
	for _, attr := range attrsList {
		if attr.Name.Local == name {
			return &attr
		}
	}

	return nil
}
