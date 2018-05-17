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
	loader   Loader
	elStack  *elementsStack
	nsStack  *stringsStack
	curNs    map[string]string
	rootNode *node

	// кеш
	attGroupsCache *NamespacedTypes

	// рабочий список типов
	typesListCache *NamespacedTypes

	// результирующий список типов
	resultTypesList []*Type
}

func NewParser(l Loader) *parser {
	return &parser{
		loader:  l,
		elStack: &elementsStack{},
		nsStack: &stringsStack{},
		curNs:   make(map[string]string),
		rootNode: &node{},
		//attGroupsCache: NewTypesCollection(),
		typesListCache: NewTypesCollection()}
}

func (p *parser) Parse(inputFile string) {
	reader, _ := p.loader.Load(inputFile)
	decoder := xml.NewDecoder(reader)
	p.parseImpl(decoder)
}

func (p *parser) parseImpl(decoder *xml.Decoder) {
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
	case "node", "simpleType", "complexType", "restriction", "sequence", "attribute", "attributeGroup", "element":
		p.elementStarted(elem)
	}
}

// Начало элемента, любого, а не только node
func (p *parser) elementStarted(e *xml.StartElement) {
	p.elStack.Push(newNode(e))
}

func (p *parser) resolveRefs(e *xml.StartElement) {
	//for _, t := range p.typesListCache.GetAllTypes() {
	//	if len
	//}
}

func (p *parser) parseEndElement(elem *xml.EndElement) {
	switch elem.Name.Local {
	case "schema":
		p.nsStack.Pop()
	case "simpleType":
		p.endSimpleType()
	case "restriction":
		p.endRestriction()
	case "complexType":
		p.endComplexType()
	case "sequence":
		p.endSequence()
	case "node":
		p.endElement()
	case "attribute":
		p.endAttribute()
	case "attributeGroup":
		p.endAttributeGroup()
	case "element":
		p.endElement()
	}
}

func (p *parser) parseSchema(elem *xml.StartElement) {
	ns := findAttributeByName(elem.Attr, "targetNamespace")
	if ns.Value != "" {
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

	p.elStack.Push(newNode(elem))
}

func (p *parser) GenerateTypes() []*Type {
	return p.generateTypesImpl(p.rootNode, 0)
}

// Сгенерировать список типов по построенному дереву
func (p *parser) generateTypesImpl(node *node, deep int) []*Type {
	var res []*Type
	for _, n := range node.children {
		if node.elemName == "element" {
			f := newField(n)
			node.genType.addField(f)
		} else {
			t := newType(n)
			t.BaseTypeName = n.typeName
			n.genType = t

			res = append(res, t)
			res = append(res, p.generateTypesImpl(n, deep+1)...)
		}
	}
	return res
}

func (p *parser) endElement() {
	e := p.elStack.Pop()
	e.namespace = p.nsStack.GetLast()
	nameAttr := findAttributeByName(e.startElem.Attr, "name")
	typeAttr := findAttributeByName(e.startElem.Attr, "type")

	if typeAttr != nil {
		e.typeName = p.createQName(typeAttr.Value)
	}
	//} else if p.elStack.Deep() > 2 {
	//	// если типа нет, придётся его создавать
	//	t := p.createAndAddType(nameAttr.Value, e)
	//	t.IsSimple = e.isSimple
	//	t.Fields = createFieldsFromElems(e.children)
	//
	//	e.typeName = p.createQName(t.Name)
	//}

	context := p.elStack.GetLast()

	if context.elemName == "schema" {
		// значит предок у нас - schema, т.е. это глобальный тип
		if nameAttr == nil {
			panic("Element should has elemName attribute")
		}

		//t := p.createAndAddType(nameAttr.Value, e)
		//t.BaseTypeName = e.typeName
		//t.IsSimple = e.isSimple
		//t.Fields = createFieldsFromElems(e.children)
		p.rootNode.add(e)
	} else if context.elemName == "sequence" {
		context.children = append(context.children, e)
	}
}

func createFieldsFromElems(elems []*node) []*Field {
	var res []*Field
	for _, f := range elems {
		nameAttr := findAttributeByName(f.startElem.Attr, "name")
		field := &Field{
			Name:     nameAttr.Value,
			TypeName: f.typeName,
			IsAttr:   f.isAttr}
		res = append(res, field)
	}

	return res
}

func (p *parser) createQName(qName string) *QName {
	typesParts := strings.Split(qName, ":")
	var (
		name, namespace string
		ok              bool
	)
	if len(typesParts) == 1 {
		name = typesParts[0]
		namespace = p.nsStack.GetLast()
	} else {
		name = typesParts[1]
		namespace, ok = p.curNs[typesParts[0]]
		if !ok {
			panic("Unknown namespace alias " + typesParts[0])
		}
	}

	return &QName{
		Name:      name,
		Namespace: namespace}
}

func findAttributeByName(attrsList []xml.Attr, name string) *xml.Attr {
	for _, attr := range attrsList {
		if attr.Name.Local == name {
			return &attr
		}
	}

	return nil
}

func (p *parser) createAndAddType(name string, e *node) *Type {
	t := p.createType(name, e)
	p.resultTypesList = append(p.resultTypesList, t)
	return t
}

func (p *parser) createType(name string, e *node) *Type {
	t := &Type{Name: name, Namespace: p.nsStack.GetLast()}
	p.typesListCache.Put(t)
	return t
}

func (p *parser) endSequence() {
	t := p.elStack.Pop()
	context := p.elStack.GetLast()
	context.children = t.children
}

func (p *parser) endComplexType() {
	e := p.elStack.Pop()
	context := p.elStack.GetLast()

	if context.elemName == "schema" {
		nameAttr := findAttributeByName(e.startElem.Attr, "name")
		t := p.createAndAddType(nameAttr.Value, e)
		t.BaseTypeName = e.typeName
	} else {
		context.children = e.children
	}
}

func (p *parser) endSimpleType() {
	e := p.elStack.Pop()
	e.namespace = p.nsStack.GetLast()
	e.isSimple = true
	nameAttr := findAttributeByName(e.startElem.Attr, "name")
	if nameAttr != nil {
		e.name = nameAttr.Value
	}

	context := p.elStack.GetLast()
	if context.elemName == "schema" {
		p.rootNode.add(e)
	} else {
		// анонимный тип, встраиваем в контейнер
		context := p.elStack.GetLast()
		context.typeName = e.typeName
		context.isSimple = e.isSimple
	}
}

func (p *parser) endRestriction() {
	e := p.elStack.Pop()
	context := p.elStack.GetLast()
	baseType := findAttributeByName(e.startElem.Attr, "base")
	context.typeName = p.createQName(baseType.Value)
}

func (p *parser) endAttribute() {
	e := p.elStack.Pop()
	e.isAttr = true
	typeAttr := findAttributeByName(e.startElem.Attr, "type")
	e.typeName = p.createQName(typeAttr.Value)
	context := p.elStack.GetLast()
	context.children = append(context.children, e)
}

func (p *parser) endAttributeGroup() {
	e := p.elStack.Pop()
	nameAttr := findAttributeByName(e.startElem.Attr, "name")
	refAttr := findAttributeByName(e.startElem.Attr, "ref")
	if nameAttr != nil {
		t := p.createType(nameAttr.Value, e)
		t.Namespace = p.nsStack.GetLast()
		t.IsSimple = false
		t.Fields = createFieldsFromElems(e.children)
	} else if refAttr != nil {
		context := p.elStack.GetLast()
		context.refs = append(context.refs, refAttr.Value)
	} else {
		panic("No elemName and no ref for attribute group")
	}
}
