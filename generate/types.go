package generate

import (
	"regexp"
	"strconv"
	"strings"
	"io"
	"github.com/dmitryvakulenko/gosoapgen/xsd/tree_parser"
	"text/template"
)

var nameRegex, _ = regexp.Compile("^[A-Za-z0-9_]+$")

var innerTypes = []string{
	"int",
	"float64",
	"bool",
	"time.Time",
	"string"}

const typeTemplate = `{{$tName := index .TypeNames .Type}}
{{- $tOrigName := .Type.Name}}
{{- if validName $tName}}
type {{$tName}} struct {
	{{- $typeNames := .TypeNames}}
	{{range $idx, $f := .Type.Fields}}
		{{- $name := title .Name}}
		{{- $fType := mapType $f.Type.Local}}
		{{- if eq $fType ""}}
			{{- $fType = index $typeNames $f.Type}}
			{{- $fType = print "*" $fType}}
		{{- end}}
		{{- if lt .MinOccurs .MaxOccurs}}{{$fType = print "[]" $fType}}{{end}}
		{{- $xml := ""}}
		{{- if $f.IsAttr}}
			{{- $xml = print ",attr,omitempty"}}
		{{- else if eq $f.Name "XMLName"}}
			{{- $xml = print $tOrigName.Space " " $tOrigName.Local}}
		{{- else if eq $f.Name "XMLValue"}}
			{{- $xml = print ",chardata"}}
		{{- else}}
			{{- $xml = print $f.Name ",omitempty"}}
		{{- end}}
		{{- if validName $name}}
			{{- $name}} {{$fType}} ` + "`" + `xml:"{{$xml}}"` + "`" + `
		{{- end}}
	{{end}}
}
{{- end}}
`

type tmplParams struct {
	TypeNames map[*tree_parser.Type]string
	Type      *tree_parser.Type
}

func Types(typesList []*tree_parser.Type, writer io.Writer) {
	nameNextIndex := make(map[string]int)
	goNames := make(map[*tree_parser.Type]string)
	for _, curType := range typesList {
		var name string
		if curType.Local != "" {
			name = curType.Local
		} else {
			panic("Anonymous type found")
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
			name = field.Name
		}

		name = firstUp(name)
		if idx, ok := nameNextIndex[name]; ok {
			nameNextIndex[name] = idx + 1
			name += "_" + strconv.Itoa(idx)
		} else {
			nameNextIndex[name] = 1
		}

		goNames[curType] = name
	}

	funcMap := template.FuncMap{
		"title":   strings.Title,
		"mapType": mapStandardType,
		"validName": isValidName}
	tmpl, _ := template.New("type").Funcs(funcMap).Parse(typeTemplate)
	for _, curType := range typesList {
		p := tmplParams{Type: curType, TypeNames: goNames}
		tmpl.Execute(writer, p)
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

func isValidName(name string) bool {
	return nameRegex.MatchString(name)
}