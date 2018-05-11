package main

import (
	"path"
	"io/ioutil"
	"errors"
)

var (
	alreadyLoadedErr = errors.New("file already loaded")
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

func (l *XsdLoader) Load(xsdFilePath string) ([]byte, error) {
	filePath := path.Clean(l.baseDir + "/" + xsdFilePath)
	var loadedErr error = nil
	if _, exists := l.alreadyLoaded[filePath]; exists {
		loadedErr = alreadyLoadedErr
	} else {
		l.alreadyLoaded[filePath] = true
	}

	res, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return res, loadedErr
}

func (l *XsdLoader) IsAlreadyLoadedError(e error) bool {
	return e == alreadyLoadedErr
}