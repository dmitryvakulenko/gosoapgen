package generate

import (
	"strings"
	"io"
	"github.com/dmitryvakulenko/gosoapgen/xsd/tree_parser"
	"text/template"
)

var innerTypes = []string{
	"int",
	"float64",
	"bool",
	"time.Time",
	"string"}


const typeTemplate = `type {{index .TypeNames .Type}} struct {
	{{- $typeNames := .TypeNames}}
	{{range $idx, $f := .Type.Fields}}
		{{- $name := title .Name}}
		{{- $fType := mapType $f.Type.Local}}
		{{- if eq $fType ""}}
			{{- $fType = index $typeNames $f.Type}}
		{{- end}}
		{{- if lt .MinOccurs .MaxOccurs}}{{$fType = print "[]" $fType}}{{end}}
		{{- $xml := ""}}
		{{- if $f.IsAttr}}
			{{- $xml = print ",attr,omitempty"}}
		{{- else if eq $f.Name "XMLName"}}
			{{- $xml = print .Type.Space " " .Type.Local}}
		{{- else if eq $f.Name "XMLValue"}}
			{{- $xml = print ",chardata"}}
		{{- else}}
			{{- $xml = print $f.Name ",omitempty"}}
		{{- end}}
		{{- $name}} {{$fType}} ` + "`" + `xml:"{{$xml}}"` + "`" + `
	{{end}}
}

`

type tmplParams struct {
	TypeNames map[*tree_parser.Type]string
	Type      *tree_parser.Type
}

func Types(typesList []*tree_parser.Type, writer io.Writer) {
	// for anonymous types
	goNames := make(map[*tree_parser.Type]string)
	for _, curType := range typesList {
		if curType.Local != "" {
			goNames[curType] = firstUp(curType.Local)
		} else {
			// find field used this type
			var field *tree_parser.Field
			for _, t := range typesList {
				for _, f := range t.Fields {
					if f.Type == curType {
						field = f
						break
					}
				}
			}
			goNames[curType] = firstUp(field.Name)
		}
	}

	funcMap := template.FuncMap{
		"title":   strings.Title,
		"mapType": mapStandardType}
	tmpl, _ := template.New("type").Funcs(funcMap).Parse(typeTemplate)
	for _, curType := range typesList {
		p := tmplParams{Type: curType, TypeNames: goNames}
		tmpl.Execute(writer, p)
		// writer.Write([]byte("type " + curType.Name.Local + " struct {\n"))
		// for _, f := range curType.Fields {
		// 	writeField(curType, f, writer)
		// }
		// writer.Write([]byte("}\n\n"))
	}
}

func writeField(t *tree_parser.Type, field *tree_parser.Field, writer io.Writer) {
	// обработка обычного поля
	writer.Write([]byte(firstUp(field.Name) + " "))

	if field.MinOccurs < field.MaxOccurs {
		writer.Write([]byte("[]"))
	}

	fieldType := mapStandardType(field.Type.Local)
	if fieldType == "" {
		fieldType = field.Type.Local
	}

	if fieldType == "" {
		fieldType = "string"
	}

	if !isInnerType(fieldType) {
		writer.Write([]byte("*"))
	}

	writer.Write([]byte(fieldType + " `xml:\""))
	if field.IsAttr {
		writer.Write([]byte(field.Name + ",attr,omitempty"))
	} else if field.Name == "XMLName" {
		writer.Write([]byte(t.Space + " " + t.Local))
	} else if field.Name == "XMLValue" {
		writer.Write([]byte(",chardata"))
	} else {
		writer.Write([]byte(field.Name + ",omitempty"))
	}
	writer.Write([]byte("\"`\n"))
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
	case "int", "integer", "positiveInteger", "nonNegativeInteger":
		return "int"
	case "decimal":
		return "float64"
	case "boolean":
		return "bool"
	case "string", "NMTOKEN", "anyURI", "language", "base64Binary", "duration", "IDREF", "IDREFS", "gYear", "gMonth", "gDay", "gYearMonth",
		"date", "dateTime", "time", "ID":
		return "string"
	default:
		return ""
	}
}
