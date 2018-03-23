package generate

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
func (c *Client) {{.Name}}(in *{{.Input}}) *{{.Output}} {
	header := c.createHeader("{{.Action}}")
	response := c.call("{{.Name}}", header, in)
	res := {{.Output}}{}
	xml.Unmarshal(response, &res)
	return &res
}
`

func Client(parser translator.Parser, wsdl *wsdl.Definitions, writer io.Writer) {
	var (
		processedTypes = make(map[string]bool)
		nsAliases = make(map[string]string)
		typeNamespace = make(map[string]string)
	)

	writer.Write([]byte("var namespaceMap = map[string]string{"))
	for idx, ns := range parser.GetNamespaces() {
		alias := "ns" + strconv.Itoa(idx)
		writer.Write([]byte("\n\"" + ns + "\": \"" + alias + "\","))
		nsAliases[ns] = alias
	}
	writer.Write([]byte("}\n\n"))
	
	for _, curType := range parser.GetTypes() {
		if _, ok := processedTypes[curType.Name]; ok {
			continue
		}

		goTypeName := firstUp(curType.Name)
		typeNamespace[goTypeName] = curType.Namespace
		processedTypes[curType.Name] = true
		writer.Write([]byte("type " + goTypeName + " struct {\n"))
		//alias := nsAliases[curType.Namespace]
		//writer.Write([]byte("XMLName string `xml:\"" + alias + ":" + curType.Name + "\"`\n"))
		for _, f := range curType.Fields {
			writer.Write([]byte(firstUp(f.Name) + " "))
			if f.MaxOccurs != 0 {
				writer.Write([]byte("[]"))
			}
			if !isInnerType(f.Type) {
				writer.Write([]byte("*"))
			}
			alias := nsAliases[f.Namespace]
			writer.Write([]byte(firstUp(f.Type) + " `xml:\"" + alias + ":" + f.Name))
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

	writer.Write([]byte("var typeNamespace = map[string]string{"))
	for typeName, ns := range typeNamespace {
		writer.Write([]byte("\n\"" + typeName + "\": \"" + ns + "\","))
	}
	writer.Write([]byte("}\n\n"))

	writeOperations(wsdl, writer)
}

func writeOperations(wsdl *wsdl.Definitions, writer io.Writer) {
	messages := make(map[string]string)
	for _, msg := range wsdl.Message {
		messages[msg.Name] = extractName(msg.Part.Element.Value)
	}

	soapActions := make(map[string]string)
	for _, op := range wsdl.Binding.Operation {
		soapActions[op.Name] = op.SoapOperation.SoapAction
	}

	tmpl, err := template.New("function").Parse(funcTemplate)
	if err != nil {
		panic(err)
	}

	for _, op := range wsdl.PortType.Operation {
		input := messages[extractName(op.Input.Message)]
		output := messages[extractName(op.Output.Message)]
		action := soapActions[op.Name]
		tmpl.Execute(writer, struct {
				Name, Input, Output, Action string
			}{op.Name, input, output, action})
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