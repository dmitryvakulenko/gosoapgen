package main

import (
	"os"
	"fmt"
	"encoding/xml"
	"path"
	"github.com/dmitryvakulenko/gosoapgen/generate"
	"github.com/dmitryvakulenko/gosoapgen/wsdl"
	"flag"
	"io"
	"github.com/dmitryvakulenko/gosoapgen/internal/pkg/xsdloader"
    "github.com/dmitryvakulenko/gosoapgen/xsd/tree_parser"
)

func main() {
	flag.Parse()

	var inName = flag.Arg(0)
	if inName == "" {
		fmt.Print("Enter input outFile")
		return
	}

	var outName = flag.Arg(1)
	if outName == "" {
		fmt.Print("Enter output outFile")
		return
	}

	var outPackage = flag.Arg(2)
	if outPackage == "" {
		fmt.Print("Enter output package name")
		return
	}

	outFile, err := os.Create(outName)
	defer outFile.Close()
	if err != nil {
		fmt.Printf("Can't write result outFile")
		return
	}

	outFile.Write([]byte("package " + outPackage + "\n\n"))
	wsdlProcessing(inName, outFile)
}

func wsdlProcessing(wsdlName string, out io.Writer) {
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

	basePath := path.Dir(wsdlName) + "/"
	parser := tree_parser.NewParser(xsdloader.NewXsdLoader())
	for _, attr := range def.Import {
		parser.Load(basePath + attr.SchemaLocation)
	}

    typesList := parser.GetTypes()
	generate.Types(typesList, out)
	generate.Client(&def, out)
}
