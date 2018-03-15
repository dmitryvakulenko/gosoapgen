package tree

import (
	"os"
	"encoding/xml"
	"io"
	"strings"
)

type Builder struct {
	typesList []*ComplexType
	curType   *ComplexType
	curField  *Field
	curNamespace string
}

func NewBuilder() Builder {
	return Builder{}
}

func (b *Builder) getTypes() []*ComplexType {
	return b.typesList
}

func (b *Builder) Build(uri string) {
	reader := makeReader(uri)
	defer reader.Close()

	decoder := xml.NewDecoder(reader)
	token, err := decoder.Token()
	for err != io.EOF {
		b.parseToken(token)
		token, err = decoder.Token()
	}
}

func (b *Builder) parseToken(token xml.Token) {
	switch t := token.(type) {
	case xml.StartElement:
		b.startElement(t)
	case xml.EndElement:
		b.endElement(t)
	case xml.CharData:

	}
}

func (b *Builder) startElement(element xml.StartElement) {
	switch element.Name.Local {
	case "schema":
		b.curNamespace = findAttributeValue(element.Attr, "targetNamespace")
	case "element":
		if b.curType == nil {
			b.curType = NewComplexType(findAttributeValue(element.Attr, "name"), b.curNamespace)
		} else {
			b.curField = b.createField(element)
		}
	case "attribute":
		b.curField = b.createField(element)
	}
}

func (b *Builder) endElement(element xml.EndElement) {
	switch element.Name.Local {
	case "element":
		if b.curField != nil {
			b.curType.Fields = append(b.curType.Fields, b.curField)
			b.curField = nil
		} else if b.curType != nil {
			b.typesList = append(b.typesList, b.curType)
			b.curType = nil
		}
	case "attribute":
		if b.curField != nil {
			b.curType.Fields = append(b.curType.Fields, b.curField)
			b.curField = nil
		}
	}
}


func (b *Builder) createType(element xml.StartElement) {
	b.curType = NewComplexType(findAttributeValue(element.Attr, "name"), b.curNamespace)
}


func (b *Builder) createField(element xml.StartElement) *Field {
	name := findAttributeValue(element.Attr, "name")
	return NewField(
		name,
		findAttributeValue(element.Attr, "type"),
		name + ",attr")
}

func makeReader(uri string) io.ReadCloser {
	reader, err := os.Open(uri)

	if err != nil {
		panic(err)
	}

	return reader
}

func findAttributeValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if a.Name.Local == name {
			return a.Value
		}
	}

	return ""
}

func parseType(xmlType string) string {
	parts := strings.Split(xmlType, ":")

	var fieldType string
	if len(parts) == 2 {
		fieldType = parts[1]
	} else {
		fieldType = parts[0]
	}

	switch fieldType {
	case "integer", "positiveInteger", "nonNegativeInteger":
		return "int"
	case "decimal":
		return "float64"
	case "boolean":
		return "bool"
	case "date":
		return "time.time"
	case "string", "NMTOKEN", "anyURI":
		return "string"
	default:
		return fieldType
	}
}
