package generator

import (
	"github.com/dmitryvakulenko/gosoapgen/translator"
	"strings"
	"strconv"
)

var innerTypes = []string{
	"int",
	"float64",
	"bool",
	"time.Time",
	"string"}

func All(parser translator.Parser) string {
	var (
		processedTypes = make(map[string]bool)
		res = ""
		nsAliases = make(map[string]string)
	)

	res = "var namespaceMap = map[string]string{"
	for idx, ns := range parser.GetNamespaces() {
		alias := "ns" + strconv.Itoa(idx)
		res += "\n\"" + alias + "\": \"" + ns + "\","
		nsAliases[ns] = alias
	}
	res += "}\n\n"

	for _, curType := range parser.GetTypes() {
		if _, ok := processedTypes[curType.Name]; ok {
			continue
		}

		processedTypes[curType.Name] = true
		res += "type " + firstUp(curType.Name) + " struct {\n"
		for _, f := range curType.Fields {
			alias := nsAliases[f.Namespace]
			res += firstUp(f.Name) + " " + firstUp(f.Type) + " `xml:\"" + alias + " " + f.XmlExpr + "\"`\n"
		}
		res += "}\n\n"

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