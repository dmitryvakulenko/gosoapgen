package _type

import (
	"io/ioutil"
	"encoding/xml"
	"os"
	"path"
)

func ParseSchema(fileName string) []*Schema {
	reader, err := os.Open(fileName)
	defer reader.Close()

	if err != nil {
		panic(err)
	}

	var res []*Schema

	content, _ := ioutil.ReadAll(reader)
	s := &Schema{}
	xml.Unmarshal(content, s)
	res = append(res, s)

	for _, s := range s.Import {
		importFile := path.Dir(fileName) + "/" + s.SchemaLocation
		res = append(res, ParseSchema(importFile)...)
	}

	for _, s := range s.Include {
		importFile := path.Dir(fileName) + "/" + s.SchemaLocation
		res = append(res, ParseSchema(importFile)...)
	}

	return res
}

