package generate

import "io"
import (
	"github.com/dmitryvakulenko/gosoapgen/wsdl"
	"text/template"
)

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
