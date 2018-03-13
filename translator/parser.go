package translator

import (
	"os"
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"encoding/xml"
)

//func DeepParse(fileName string) {
//	res :=
//}

func readFile(fileName string) *xsd.Schema {
	reader, err := os.Open(fileName)
	defer reader.Close()

	if err != nil {
		panic(err)
	}

	s := xsd.Schema{}
	err = xml.NewDecoder(reader).Decode(&s)
	if err != nil {
		panic(err)
	}

	return &s
}