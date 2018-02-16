package main

import (
	"os"
	"fmt"
	"encoding/xml"
	"gosoapgen/wsdl"
	"path"
	"gosoapgen/xsd"
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
	xmlFile.Close()

	basePath := path.Dir(wsdlName)
	res := ""
	for _, attr := range def.Types {
		//schema, _ := xsd.LoadSchema(basePath + "/" + attr.SchemaLocation)
		reader, _ := os.Open(basePath + "/" + attr.SchemaLocation)
		xsd.Parse(reader)
		reader.Close()
		//xsd.Parse(basePath + "/" + attr.SchemaLocation)
		//if err != nil {
		//	fmt.Printf("Error loading schema. %s", err)
		//	return
		//}
		//if outFilePath, err := schema.MakeGoPkgSrcFile(); err == nil {
		//	if raw, _ := exec.Command("gofmt", "-w=true", "-s=true", "-e=true", outFilePath).CombinedOutput(); len(raw) > 0 {
		//		log.Printf("GOFMT:\t%s\n", string(raw))
		//	}
		//}
	}
	xsd.GetTypes()
	fmt.Print(res)
}
