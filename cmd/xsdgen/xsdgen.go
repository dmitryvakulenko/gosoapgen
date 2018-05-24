package main

import (
	"fmt"
	"github.com/dmitryvakulenko/gosoapgen/generate"
	"flag"
	"io"
	"github.com/dmitryvakulenko/gosoapgen/internal/pkg/xsdloader"
	"github.com/dmitryvakulenko/gosoapgen/xsd/tree_parser"
	"os"
	"log"
)

func main() {
	xsds, outName, packageName := parseArguments()

	outFile, err := os.Create(outName)
	defer outFile.Close()
	if err != nil {
		log.Fatalf("Can't write result outFile %s", err)
		return
	}

	outFile.Write([]byte("package " + packageName + "\n\n"))
	xsdProcessing(xsds, outFile)
}

func parseArguments() ([]string, string, string) {
	flag.Parse()
	allArgs := flag.Args()
	amount := len(allArgs)
	if amount < 3 {
		fmt.Println("Usage: [xsd-files...] out-file out-package")
		os.Exit(0)
	}

	out := allArgs[amount-2:]

	return allArgs[:amount-2], out[0], out[1]
}

func xsdProcessing(xsds []string, out io.Writer) {
	parser := tree_parser.NewParser(xsdloader.NewXsdLoader())

	for _, xsdName := range xsds {
		parser.Load(xsdName)
	}

	typesList := parser.GetTypes()
	generate.Types(typesList, out)
}
