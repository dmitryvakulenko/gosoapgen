package main

import (
	"os"
	"fmt"
	"encoding/xml"
	"gosoapgen/wsdl"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("Enter wsdl")
		return
	}

	wsdlName := os.Args[1]
	stat, err := os.Stat(wsdlName)
	if os.IsNotExist(err) {
		fmt.Printf("File %s not exists", wsdlName)
		return
	}

    xmlFile, err := os.Open(wsdlName)
	xmlData := make([]byte, stat.Size())
	_, _ = xmlFile.Read(xmlData)

	def := wsdl.Definitions{}
	err = xml.Unmarshal(xmlData, &def)
	for _, attr := range def.PortType[0].Operation {
		fmt.Printf("%s: %s\n", attr.Name.Value, attr.Input.Message)
	}
}
