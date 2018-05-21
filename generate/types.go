package generate

import (
	"strings"
	"io"
	"github.com/dmitryvakulenko/gosoapgen/xsd/tree_parser"
)

var innerTypes = []string{
	"int",
	"float64",
	"bool",
	"time.Time",
	"string"}

const funcTemplate = `
func (c *SoapClient) {{.Name}}(body *{{.Input}}) *{{.Output}} {
	header := c.transporter.CreateHeader("{{.Action}}")
	response := c.transporter.Send("{{.Action}}", header, body)
	res := {{.Output}}{}
	xml.Unmarshal(response, &res)
	return &res
}
`

func Types(typesList []*tree_parser.Type, writer io.Writer) {
	for _, curType := range typesList {
		if curType.IsSimpleContent && len(curType.Fields) == 0 {
			writer.Write([]byte("type " + curType.Name + " " + curType.BaseType.Name + "\n\n"))
			continue
		}

		writer.Write([]byte("type " + firstUp(curType.Name) + " struct {\n"))
		for _, f := range curType.Fields {
			writeField(f, curType.Namespace, writer)
		}
		if curType.IsSimpleContent {
			writer.Write([]byte("Value string `xml:\"chardata\"`"))
		}
		writer.Write([]byte("}\n\n"))
	}
}

func writeField(field *tree_parser.Field, ns string, writer io.Writer) {
	if field.IsAttr && len(field.Type.Fields) != 0 {
		// обработка attributeGroup
		for _, f := range field.Type.Fields {
			writeField(f, field.Type.Namespace, writer)
		}
	} else {
		// обработка обычного поля
		writer.Write([]byte(firstUp(field.Name) + " "))
		if field.MaxOccurs != 0 {
			writer.Write([]byte("[]"))
		}

		if !isInnerType(field.Type.Name) && !field.Type.IsSimpleContent {
			writer.Write([]byte("*"))
		}

		writer.Write([]byte(firstUp(field.Type.Name) + " `xml:\"" + ns + " " + field.Name))
		if field.IsAttr {
			writer.Write([]byte(",attr"))

		}
		if field.MinOccurs == 0 {
			writer.Write([]byte(",omitempty"))
		}
		writer.Write([]byte("\"`\n"))
	}
}


func extractName(in string) string {
	parts := strings.Split(in, ":")
	if len(parts) == 2 {
		return parts[1]
	} else {
		return parts[0]
	}
}

func firstUp(text string) string {
	if isInnerType(text) {
		return text
	}
	return strings.Title(text)
}

func isInnerType(t string) bool {
	for _, v := range innerTypes {
		if v == t {
			return true
		}
	}
	return false
}