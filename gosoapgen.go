package main

import (
	"os"
	"fmt"
	"encoding/xml"
	"gosoapgen/wsdl"
	"path"
	"gosoapgen/xsd"
	"gosoapgen/translator"
	"gosoapgen/generating"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("Enter wsdl")
		return
	}

	wsdlName := os.Args[1]
	stat, err := os.Stat(wsdlName)
	if os.IsNotExist(err) {
		fmt.Printf("File %q not exists", wsdlName)
		return
	}

    xmlFile, err := os.Open(wsdlName)
	xmlData := make([]byte, stat.Size())
	_, _ = xmlFile.Read(xmlData)

	def := wsdl.Definitions{}
	err = xml.Unmarshal(xmlData, &def)
	xmlFile.Close()

	basePath := path.Dir(wsdlName)
	schemas := []*xsd.Schema{}
	for _, attr := range def.Types {
		s := xsd.ParseSchema(basePath + "/" + attr.SchemaLocation)
		schemas = append(schemas, s...)
	}

	structs := translator.GenerateTypes(schemas)
	res := generating.Types(structs)

	file, err := os.Open("./res.go")
	file.Write([]byte(res))
	file.Close()
}
