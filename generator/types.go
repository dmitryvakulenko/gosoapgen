package generator

import (
	"github.com/dmitryvakulenko/gosoapgen/translator"
	"strings"
)

var innerTypes = []string{
	"int",
	"float64",
	"bool",
	"time.Time",
	"string"}

func All(typesList []*translator.ComplexType) string {
	res := ""
	for _, curType := range typesList {
		res += "type " + firstUp(curType.Name) + " struct {\n"
		for _, f := range curType.Fields {
			res += firstUp(f.Name) + " " + firstUp(f.Type) + " `xml:\"" + f.XmlExpr + "\"`\n"
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