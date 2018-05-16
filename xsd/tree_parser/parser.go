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
	loader  Loader
	elStack *elementsStack
	nsStack *stringsStack
	curNs   map[string]string

	// результирующий список типов
	typesList *NamespacedTypes
}

func NewParser(l Loader) *parser {
	return &parser{
		loader:    l,
		elStack:   &elementsStack{},
		nsStack:   &stringsStack{},
		curNs:     make(map[string]string),
		typesList: NewTypesCollection()}
}

func (p *parser) Parse(inputFile string) {
	reader, _ := p.loader.Load(inputFile)
	decoder := xml.NewDecoder(reader)
	p.parseImpl(decoder)
}

func (p *parser) parseImpl(decoder *xml.Decoder) {
	for token, err := decoder.Token(); err != io.EOF; token, err = decoder.Token() {
		switch elem := token.(type) {
		case xml.StartElement:
			p.parseStartElement(&elem)
		case xml.EndElement:
			p.parseEndElement(&elem)
		}
	}
}

func (p *parser) parseStartElement(elem *xml.StartElement) {
	switch elem.Name.Local {
	case "schema":
		p.parseSchema(elem)
	case "element", "simpleType", "restriction":
		p.elementStarted(elem)
	//case "complexType":
	//	p.parseComplexType(elem)
	//case "restriction":
	//	p.parseRestriction(elem)
	//case "element":
	//	p.startElement(elem)
	//case "sequence":
	//	p.parseSequence(elem)
	}
}

// Начало элемента, любого, а не только element
func (p *parser) elementStarted(e *xml.StartElement) {
	p.elStack.Push(newElement(e))
}

func (p *parser) parseEndElement(elem *xml.EndElement) {
	switch elem.Name.Local {
	case "schema":
		p.nsStack.Pop()
	case "simpleType":
		p.endSimpleType()
	case "restriction":
		p.endRestriction()
	//case "complexType":
	//	p.endComplexType()
	//case "sequence":
	//	p.endSequence()
	case "element":
		p.endElement()
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

//func (p *parser) startSimpleType(e *xml.StartElement) {
//	nameAttr := findAttributeByName(e.Attr, "name")
//	if nameAttr == nil {
//		// анонимный тип, может быть только встроен в element
//	} else {
//		curType := p.createType(e.Name.Local, nameAttr.Value)
//		curType.IsSimple = true
//	}
//}

func (p *parser) GetTypes() []*Type {
	return p.typesList.GetAllTypes()
}

//func (p *parser) startElement(e *xml.StartElement) {
//	t := p.createType(e.Name.Local, findAttributeByName(e.Attr, "name").Value)
//	typeAttr := findAttributeByName(e.Attr, "name")
//	if typeAttr != nil {
//		t.BaseTypeName = p.createQName(typeAttr.Value)
//	}
//}
//
func (p *parser) endElement() {
	e := p.elStack.Pop()
	nameAttr := findAttributeByName(e.startElem.Attr, "name")
	if nameAttr != nil {
		t := p.createType(nameAttr.Value)
		t.BaseTypeName = e.typeName
		t.IsSimple = e.isSimple
	} else {

	}
}

//func (p *parser) parseRestriction(element *xml.StartElement) {
//	c := p.elStack.GetLast()
//	if c.IsSimple {
//		baseTypeName := findAttributeByName(element.Attr, "base")
//		c.BaseTypeName = p.createQName(baseTypeName.Value)
//	}
//}

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

func (p *parser) createType(name string) *Type {
	t := &Type{Name: name, Namespace: p.nsStack.GetLast()}
	p.typesList.Put(t)
	return t
}

// создаёт анонимный тип, в список типов не помещает
// анонимные типы выступают просто контейнерами для других
//func (p *parser) anonTypeStarted(sourceElement string) *Type {
//	t := &Type{Element: sourceElement}
//	p.elStack.Push(t)
//	return t
//}
//
//func (p *parser) parseSequence(e *xml.StartElement) {
//	p.anonTypeStarted(e.Name.Local)
//}
//
//func (p *parser) endSequence() {
//	t := p.elStack.Pop()
//	context := p.elStack.GetLast()
//	context.Fields = t.Fields
//}
//
//func (p *parser) parseComplexType(e *xml.StartElement) {
//	nameAttr := findAttributeByName(e.Attr, "name")
//	if nameAttr == nil {
//		p.anonTypeStarted(e.Name.Local)
//	} else {
//		// обычный тип с именем
//		p.createType(e.Name.Local, nameAttr.Value)
//	}
//}
//
//func (p *parser) endComplexType() {
//	t := p.elStack.Pop()
//	if len(t.Fields) != 0 {
//		t.IsSimple = false
//	}
//
//	if t.Name != "" {
//		p.typesList.Put(t)
//	} else {
//		context := p.elStack.GetLast()
//		context.Fields = t.Fields
//	}
//}

func (p *parser) endSimpleType() {
	e := p.elStack.Pop()
	eName := findAttributeByName(e.startElem.Attr, "name")
	if eName != nil {
		// отдельный глобальный тип
		t := p.createType(eName.Value)
		t.BaseTypeName = e.typeName
		t.IsSimple = true
	} else {
		// анонимный тип, встраиваем в контейнер
		context := p.elStack.GetLast()
		context.typeName = e.typeName
		context.isSimple = true
	}
}

func (p *parser) endRestriction() {
	e := p.elStack.Pop()
	context := p.elStack.GetLast()
	baseType := findAttributeByName(e.startElem.Attr, "base")
	context.typeName = p.createQName(baseType.Value)
}
