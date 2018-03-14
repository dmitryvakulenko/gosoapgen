package main

import (
	"os"
	"fmt"
	"encoding/xml"
	"github.com/dmitryvakulenko/gosoapgen/wsdl"
	"path"
	"github.com/dmitryvakulenko/gosoapgen/translator"
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
	parser := translator.NewParser()
	for _, attr := range def.Import {
		parser.Parse(path.Clean(basePath + "/" + attr.SchemaLocation))
	}

	fmt.Printf("DONE")
}
