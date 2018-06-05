/*
Парсер схемы xsd.
Использует ручной обход всего дерева, вместо загрузки его с помощь Unmarshal.
 */
package tree_parser

import (
    "io"
    "encoding/xml"
    "strings"
    "strconv"
)

var (
    stringQName = &QName{Name: "string", Namespace: "http://www.w3.org/2001/XMLSchema"}
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
    // basePath string
}

func NewParser(l Loader) *parser {
    return &parser{
        loader:   l,
        elStack:  &elementsStack{},
        nsStack:  &stringsStack{},
        curNs:    make(map[string]string),
        rootNode: &node{}}
}

func (p *parser) Load(inputFile string) {
    reader, err := p.loader.Load(inputFile)
    defer reader.Close()
    if err != nil {
        if p.loader.IsAlreadyLoadedError(err) {
            return
        } else {
            panic(err)
        }
    }

    decoder := xml.NewDecoder(reader)
    p.decodeXsd(decoder)
}

func (p *parser) decodeXsd(decoder *xml.Decoder) {
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
    case "node", "simpleType", "complexType", "extension", "restriction", "sequence", "attribute", "attributeGroup", "element", "union",
        "simpleContent", "complexContent", "choice":
        p.elementStarted(elem)
    case "include", "import":
        p.includeStarted(elem)
    }
}

// Начало элемента, любого, а не только node
func (p *parser) elementStarted(e *xml.StartElement) {
    p.elStack.Push(newNode(e))
}

func (p *parser) parseEndElement(elem *xml.EndElement) {
    switch elem.Name.Local {
    case "schema":
        p.nsStack.Pop()
    case "simpleType":
        p.endSimpleType()
    case "extension", "restriction":
        p.endExtensionRestriction()
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
    case "union":
        p.endUnion()
    case "simpleContent":
        p.endSimpleContent()
    case "complexContent":
        p.endComplexContent()
    case "choice":
        p.endChoice()
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

    // p.curNs = make(map[string]string)
    for _, attr := range elem.Attr {
        if attr.Name.Space == "xmlns" && attr.Name.Local != "" {
            p.curNs[attr.Name.Local] = attr.Value
        }
    }

    p.elStack.Push(newNode(elem))
}

func (p *parser) GetTypes() []*Type {
    l := p.parseTypesImpl(p.rootNode)
    p.linkTypes(l)
    p.renameDuplicatedTypes(l)
    p.foldSimpleTypes(l)
    return l
}

// Сгенерировать список типов по построенному дереву
func (p *parser) parseTypesImpl(node *node) []*Type {
    var res []*Type
    for _, n := range node.children {
        f := newField(n)

        if node != p.rootNode {
            node.genType.addField(f)
        }

        if node == p.rootNode || len(n.children) > 0 {
            t := newType(n)
            t.BaseTypeName = n.typeName
            n.genType, f.Type = t, t
            res = append(res, t)
        }

        res = append(res, p.parseTypesImpl(n)...)
    }

    // простое содержимое с атрибутами
    if len(node.children) > 0 && node.isSimpleContent {
        f := &Field{
            Name:     "Value",
            TypeName: node.typeName}
        node.genType.addField(f)
        // node.genType.BaseTypeName = nil
    }

    return res
}

// связать все типы
func (p *parser) linkTypes(typesList []*Type) {
    for _, t := range typesList {
        if t.BaseTypeName != nil {
            t.BaseType = p.findGlobalTypeNode(*t.BaseTypeName).genType
        }

        for _, f := range t.Fields {
            if f.Type != nil {
                continue
            }

            fTypeNode := p.findGlobalTypeNode(*f.TypeName)
            if fTypeNode.elemName == "attributeGroup" {
                f.Name = fTypeNode.name
            }
            f.Type = fTypeNode.genType
        }
    }
}

func (p *parser) renameDuplicatedTypes(typesList []*Type) {
    names := make(map[string]int)
    for _, t := range typesList {
        if _, exist := names[t.Name]; exist {
            names[t.Name]++
            t.Name = t.Name + "_" + strconv.Itoa(names[t.Name])
        } else {
            names[t.Name] = 0
        }
    }
}

func (p *parser) foldSimpleTypes(typesList []*Type) {
    for _, t := range typesList {
        for _, f := range t.Fields {
            if f.Type.IsSimpleContent && len(f.Type.Fields) == 0 {
                f.Type = getLastType(f.Type)
            }
        }
    }
}

func getLastType(t *Type) *Type {
    if t.BaseType == nil || len(t.Fields) != 0 {
        return t
    } else {
        return getLastType(t.BaseType)
    }
}

func (p *parser) findGlobalTypeNode(name QName) *node {
    if name.Namespace == "http://www.w3.org/2001/XMLSchema" {
        return &node{genType: &Type{Name: name.Name, Namespace: name.Namespace, IsSimpleContent: true}}
    }

    for _, t := range p.rootNode.children {
        if t.name == name.Name && t.namespace == name.Namespace {
            return t
        }
    }

    panic("Can't find type " + name.Name)
}

func (p *parser) endElement() {
    e := p.elStack.Pop()
    e.namespace = p.nsStack.GetLast()

    maxAttr := findAttributeByName(e.startElem.Attr, "maxOccurs")
    if maxAttr != nil {
        e.isArray = true
    }

    typeAttr := findAttributeByName(e.startElem.Attr, "type")
    if typeAttr != nil {
        e.typeName = p.createQName(typeAttr.Value)
    }

    refAttr := findAttributeByName(e.startElem.Attr, "ref")
    if refAttr != nil {
        e.typeName = p.createQName(refAttr.Value)
        e.name = e.typeName.Name
    }

    context := p.elStack.GetLast()

    if context.elemName == "schema" {
        nameAttr := findAttributeByName(e.startElem.Attr, "name")
        // значит предок у нас - schema, т.е. это глобальный тип
        if nameAttr == nil {
            panic("Element should has elemName attribute")
        }

        p.rootNode.addChild(e)
    } else if context.elemName == "sequence" || context.elemName == "choice" {
        context.children = append(context.children, e)
    }

    if e.typeName == nil && len(e.children) == 0 {
        e.typeName = stringQName
    }
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

func (p *parser) endSequence() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    context.children = append(context.children, e.children...)
}

func (p *parser) endComplexType() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()

    if context.elemName == "schema" {
        nameAttr := findAttributeByName(e.startElem.Attr, "name")
        e.name = nameAttr.Value
        e.namespace = p.nsStack.GetLast()
        p.rootNode.addChild(e)
    } else {
        context.isSimpleContent = e.isSimpleContent
        context.isAttr = e.isAttr
        context.children = append(context.children, e.children...)
        context.typeName = e.typeName
    }
}

func (p *parser) endSimpleType() {
    e := p.elStack.Pop()
    e.namespace = p.nsStack.GetLast()
    e.isSimpleContent = true
    nameAttr := findAttributeByName(e.startElem.Attr, "name")
    if nameAttr != nil {
        e.name = nameAttr.Value
    }

    context := p.elStack.GetLast()
    if context.elemName == "schema" {
        p.rootNode.addChild(e)
    } else {
        // анонимный тип, встраиваем в контейнер
        context := p.elStack.GetLast()
        context.typeName = e.typeName
        context.isSimpleContent = e.isSimpleContent
    }
}

func (p *parser) endExtensionRestriction() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    baseType := findAttributeByName(e.startElem.Attr, "base")
    context.typeName = p.createQName(baseType.Value)
    context.children = append(context.children, e.children...)
}

func (p *parser) endAttribute() {
    e := p.elStack.Pop()
    e.isAttr = true
    if e.typeName == nil {
        typeAttr := findAttributeByName(e.startElem.Attr, "type")
        if typeAttr != nil {
            e.typeName = p.createQName(typeAttr.Value)
        } else {
            e.typeName = stringQName
        }
    }

    context := p.elStack.GetLast()
    context.children = append(context.children, e)
}

func (p *parser) endAttributeGroup() {
    e := p.elStack.Pop()
    e.isAttr = true
    e.namespace = p.nsStack.GetLast()

    context := p.elStack.GetLast()
    nameAttr := findAttributeByName(e.startElem.Attr, "name")
    refAttr := findAttributeByName(e.startElem.Attr, "ref")
    if context.elemName == "schema" {
        e.name = nameAttr.Value
        p.rootNode.addChild(e)
    } else if refAttr != nil {
        e.typeName = p.createQName(refAttr.Value)
        context.addChild(e)
    } else {
        panic("No elemName and no ref for attribute group")
    }
}

func (p *parser) includeStarted(e *xml.StartElement) {
    l := findAttributeByName(e.Attr, "schemaLocation")
    p.Load(l.Value)
}

func (p *parser) endUnion() {
    p.elStack.Pop()
    context := p.elStack.GetLast()
    context.typeName = stringQName
}

func (p *parser) endSimpleContent() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    context.isSimpleContent = true
    context.typeName = e.typeName
    context.children = append(context.children, e.children...)
}

func (p *parser) endComplexContent() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    context.isSimpleContent = false
    context.typeName = e.typeName
    context.children = append(context.children, e.children...)
}

func (p *parser) endChoice() {
    e := p.elStack.Pop()
    context := p.elStack.GetLast()
    context.children = append(context.children, e.children...)
}
