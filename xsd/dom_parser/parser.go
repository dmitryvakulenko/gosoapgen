package dom_parser

import (
	dom "github.com/subchen/go-xmldom"
	"fmt"
)

func Parse(fileName string) {
	doc, err := dom.ParseFile(fileName)
	if err != nil {
		panic(err)
	}

	fmt.Println(doc.Root.Name)
}