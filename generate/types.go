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

		fieldType := mapStandardType(field.Type.Name)
		if fieldType == "" {
			fieldType = firstUp(field.Type.Name)
		}

		writer.Write([]byte("*" + fieldType + " `xml:\"" + ns + " " + field.Name))
		if field.IsAttr {
			writer.Write([]byte(",attr"))

		}

		if field.Type.IsSimpleContent && field.Name == "Value" {
			writer.Write([]byte(",chardata"))
		} else if field.MinOccurs == 0 {
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

func mapStandardType(xmlType string) string {
	switch xmlType {
	case "int", "integer", "positiveInteger", "nonNegativeInteger", "ID":
		return "int"
	case "decimal":
		return "float64"
	case "boolean":
		return "bool"
	case "date", "dateTime", "time":
		return "time.Time"
	case "string", "NMTOKEN", "anyURI", "language", "base64Binary", "duration", "IDREF", "IDREFS", "gYear", "gMonth", "gDay", "gYearMonth":
		return "string"
	default:
		return ""
	}
}