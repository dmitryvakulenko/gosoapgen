package main

import (
	"os"
	"fmt"
	"path"
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"github.com/dmitryvakulenko/gosoapgen/generate"
	"flag"
	"io"
	"github.com/dmitryvakulenko/gosoapgen/internal/pkg/xsdloader"
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
	xsdProcessing(inName, outFile)
}


func xsdProcessing(xsdName string, out io.Writer) {
	basePath := path.Dir(xsdName)
	parser := xsd.NewDecoder(xsdloader.NewXsdLoader(basePath))
	parser.Decode(path.Base(xsdName))

	generate.Types(parser, out)
}
