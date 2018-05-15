package generate

import (
	"strings"
	"github.com/dmitryvakulenko/gosoapgen/wsdl"
	"text/template"
	"io"
	"github.com/dmitryvakulenko/gosoapgen/xsd/flat_parser"
	"github.com/metaleap/go-xsd"
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

func Types(parser flat_parser.Decoder, writer io.Writer) {
	for _, t := range parser.GetTypes() {
		switch curType := t.(type) {
		case *flat_parser.ComplexType:
			writer.Write([]byte("type " + curType.GoName + " struct {\n"))
			//typeNamespace[curType.GoName] = curType.Namespace
			//alias := nsAliases[curType.Namespace]
			//writer.Write([]byte("XMLName string `xml:\"" + alias + ":" + curType.GoName + "\"`\n"))
			for _, f := range curType.Fields {
				writer.Write([]byte(firstUp(f.Name) + " "))
				if f.MaxOccurs != 0 {
					writer.Write([]byte("[]"))
				}

				_, isSimpleType := f.Type.(*flat_parser.SimpleType)
				if !isInnerType(f.Type.GetName()) && !isSimpleType {
					writer.Write([]byte("*"))
				}

				//alias := ""
				//if !f.IsAttr {
				//	alias = nsAliases[curType.Namespace] + ":"
				//}
				writer.Write([]byte(firstUp(f.Type.GetGoName()) + " `xml:\"" + curType.Namespace + " " + f.Name))
				if f.IsAttr {
					writer.Write([]byte(",attr"))

				}
				if f.MinOccurs == 0 {
					writer.Write([]byte(",omitempty"))
				}
				writer.Write([]byte("\"`\n"))
			}
			writer.Write([]byte("}\n\n"))
		case *flat_parser.SimpleType:
			writer.Write([]byte("type " + curType.GoName + " " + curType.BaseType.GetName() + "\n\n"))
		}
	}
}

//func Client(wsdl *wsdl.Definitions, writer io.Writer) {
	//var (
	//	nsAliases = make(map[string]string)
	//	typeNamespace = make(map[string]string)
	//)

	//writer.Write([]byte("var namespaceMap = map[string]string{"))
	//for idx, ns := range parser.GetNamespaces() {
	//	alias := "ns" + strconv.Itoa(idx)
	//	writer.Write([]byte("\n\"" + ns + "\": \"" + alias + "\","))
	//	nsAliases[ns] = alias
	//}
	//writer.Write([]byte("}\n\n"))
	


	//writer.Write([]byte("var typeNamespace = map[string]string{"))
	//for typeName, ns := range typeNamespace {
	//	writer.Write([]byte("\n\"" + typeName + "\": \"" + ns + "\","))
	//}
	//writer.Write([]byte("}\n\n"))

	//writeOperations(wsdl, writer)
//}

func Client(wsdl *wsdl.Definitions, writer io.Writer) {
	writer.Write([]byte(`

type Transporter interface {
	Send(string, interface{}, interface{}) []byte
	CreateHeader(string) interface{}
}

type SoapClient struct {
	transporter Transporter
}

func NewSoapClient(t Transporter) SoapClient {
	return SoapClient{transporter: t}
}
`))

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