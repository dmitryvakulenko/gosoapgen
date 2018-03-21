package generator

import (
	"github.com/dmitryvakulenko/gosoapgen/translator"
	"strings"
	"strconv"
	"github.com/dmitryvakulenko/gosoapgen/wsdl"
	"text/template"
	"io"
)

var innerTypes = []string{
	"int",
	"float64",
	"bool",
	"time.Time",
	"string"}

const funcTemplate = `
func (*NewClient) {{.Name}}() {
}
`

func All(parser translator.Parser, operations []*wsdl.Operation, writer io.Writer) {
	var (
		processedTypes = make(map[string]bool)
		nsAliases = make(map[string]string)
		typeNamespace = make(map[string]string)
	)

	writer.Write([]byte("var namespaceMap = map[string]string{"))
	for idx, ns := range parser.GetNamespaces() {
		alias := "ns" + strconv.Itoa(idx)
		writer.Write([]byte("\n\"" + alias + "\": \"" + ns + "\","))
		nsAliases[ns] = alias
	}
	writer.Write([]byte("}\n\n"))
	
	for _, curType := range parser.GetTypes() {
		if _, ok := processedTypes[curType.Name]; ok {
			continue
		}

		goTypeName := strings.Title(curType.Name)
		typeNamespace[goTypeName] = curType.Namespace
		processedTypes[curType.Name] = true
		writer.Write([]byte("type " + goTypeName + " struct {\n"))
		for _, f := range curType.Fields {
			alias := nsAliases[f.Namespace]
			writer.Write([]byte(strings.Title(f.Name) + " " + strings.Title(f.Type) + " `xml:\"" + alias + " " + f.XmlExpr + "\"`\n"))
		}
		writer.Write([]byte("}\n\n"))
	}

	writer.Write([]byte("var typeNamespace = map[string]string{"))
	for typeName, ns := range typeNamespace {
		writer.Write([]byte("\n\"" + typeName + "\": \"" + ns + "\","))
	}
	writer.Write([]byte("}\n\n"))

	writeOperations(operations, writer)
}

func writeOperations(operations []*wsdl.Operation, writer io.Writer) {
	tmpl, err := template.New("function").Parse(funcTemplate)
	if err != nil {
		panic(err)
	}

	for _, op := range operations {
		tmpl.Execute(writer, op)
	}
}
