package generator

import (
	"github.com/dmitryvakulenko/gosoapgen/translator"
	"strings"
	"strconv"
	"github.com/dmitryvakulenko/gosoapgen/wsdl"
)

var innerTypes = []string{
	"int",
	"float64",
	"bool",
	"time.Time",
	"string"}

const funcTemplate = `
func (*NewClient) {{.Name}}
`

func All(parser translator.Parser, operations []*wsdl.Operation) string {
	var (
		processedTypes = make(map[string]bool)
		res = ""
		nsAliases = make(map[string]string)
		typeNamespace = make(map[string]string)
	)

	res = "var namespaceMap = map[string]string{"
	for idx, ns := range parser.GetNamespaces() {
		alias := "ns" + strconv.Itoa(idx)
		res += "\n\"" + alias + "\": \"" + ns + "\","
		nsAliases[ns] = alias
	}
	res += "}\n\n"

	typesDef := ""
	for _, curType := range parser.GetTypes() {
		if _, ok := processedTypes[curType.Name]; ok {
			continue
		}

		goTypeName := firstUp(curType.Name)
		typeNamespace[goTypeName] = curType.Namespace
		processedTypes[curType.Name] = true
		typesDef += "type " + goTypeName + " struct {\n"
		for _, f := range curType.Fields {
			alias := nsAliases[f.Namespace]
			typesDef += firstUp(f.Name) + " " + firstUp(f.Type) + " `xml:\"" + alias + " " + f.XmlExpr + "\"`\n"
		}
		typesDef += "}\n\n"

	}

	res = "var typeNamespace = map[string]string{"
	for typeName, ns := range typeNamespace {
		res += "\n\"" + typeName + "\": \"" + ns + "\","
	}
	res += "}\n\n"

	return res + typesDef
}

func buildOperations(operations []*wsdl.Operation) string {
	res := ""
	for _, v := range operations {
		res +=
	}
	return res
}

func firstUp(text string) string {
	for _, v := range innerTypes {
		if v == text {
			return text
		}
	}
	return strings.ToUpper(text[0:1]) + text[1:]
}