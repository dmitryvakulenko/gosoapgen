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
		if curType.IsSimpleContent {
			writer.Write([]byte("type " + curType.Name + " " + curType.BaseType.Name + "\n\n"))
			continue
		}

		writer.Write([]byte("type " + curType.GoName + " struct {\n"))
		for _, f := range curType.Fields {
			writer.Write([]byte(firstUp(f.Name) + " "))
			if f.MaxOccurs != 0 {
				writer.Write([]byte("[]"))
			}

			if !isInnerType(f.Type.Name) && !f.Type.IsSimpleContent {
				writer.Write([]byte("*"))
			}

			writer.Write([]byte(firstUp(f.Type.Name) + " `xml:\"" + curType.Namespace + " " + f.Name))
			if f.IsAttr {
				writer.Write([]byte(",attr"))

			}
			if f.MinOccurs == 0 {
				writer.Write([]byte(",omitempty"))
			}
			writer.Write([]byte("\"`\n"))
		}
		writer.Write([]byte("}\n\n"))
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