package main

import (
	"path"
	"io/ioutil"
)

type XsdLoader struct {
	alreadyLoaded map[string]bool
	baseDir string
}

func newXsdLoader(baseDir string) *XsdLoader {
	return &XsdLoader{
		alreadyLoaded: make(map[string]bool),
		baseDir: baseDir}
}

func (l *XsdLoader) Load(xsdFilePath string) ([]byte, bool) {
	filePath := path.Clean(l.baseDir + "/" + xsdFilePath)
	if _, exists := l.alreadyLoaded[filePath]; exists {
		return make([]byte, 0), true
	}

	res, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return res, false
}