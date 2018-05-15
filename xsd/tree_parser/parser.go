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
	case "complexType":
		p.parseComplexType(elem)
	case "restriction":
		p.parseRestriction(elem)
	case "element":
		p.parseElement(elem)
	case "sequence":
		p.parseSequence(elem)
	}
}

func (p *parser) parseEndElement(elem *xml.EndElement) {
	switch elem.Name.Local {
	case "schema":
		p.nsStack.Pop()
	case "simpleType":
		p.typesStack.Pop()
	case "complexType":
		p.endComplexType()
	case "sequence":
		p.endSequence()
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
	nameAttr := findAttributeByName(elem.Attr, "name")
	if nameAttr == nil {
		// анонимный тип, может быть только встроен в element
	} else {
		curType := p.typeStarted(nameAttr.Value)
		curType.IsSimple = true
	}
}

func (p *parser) GetTypes() []*Type {
	return p.typesList.GetAllTypes()
}

func (p *parser) parseElement(elem *xml.StartElement) {
	p.typeStarted(findAttributeByName(elem.Attr, "name").Value)
}

func (p *parser) parseRestriction(element *xml.StartElement) {
	c := p.typesStack.GetLast()
	if c.IsSimple {
		baseTypeName := findAttributeByName(element.Attr, "base")
		c.BaseTypeName = p.createQName(baseTypeName.Value)
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

// создаёт реальный, именованный тип и добавляет его во все списки
func (p *parser) typeStarted(name string) *Type {
	t := &Type{Name: name, Namespace: p.nsStack.GetLast()}
	p.typesStack.Push(t)
	p.typesList.Put(t)
	return t
}

// создаёт анонимный тип, в список типов не помещает
// анонимные типы выступают просто контейнерами для других
func (p *parser) anonTypeStarted() *Type {
	t := &Type{}
	p.typesStack.Push(t)
	return t
}


func (p *parser) parseComplexType(elem *xml.StartElement) {
	nameAttr := findAttributeByName(elem.Attr, "name")
	if nameAttr == nil {
		p.anonTypeStarted()
	} else {
		// обычный тип с именем
		p.typeStarted(nameAttr.Value)
	}
}
func (p *parser) parseSequence(element *xml.StartElement) {
	p.anonTypeStarted()
}


func (p *parser) endSequence() {
	t := p.typesStack.Pop()
	context := p.typesStack.GetLast()
	context.Fields = t.Fields
}


func (p *parser) endComplexType() {
	t := p.typesStack.Pop()
	if len(t.Fields) != 0 {
		t.IsSimple = false
	}

	if t.Name != "" {
		p.typesList.Put(t)
 	} else {
 		context := p.typesStack.GetLast()
 		context.Fields = t.Fields
	}
}