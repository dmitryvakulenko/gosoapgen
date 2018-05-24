package xsdloader

import (
	"path"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	alreadyLoadedErr = errors.New("file already loaded")
)

type XsdLoader struct {
	alreadyLoaded map[string]bool
	curDir        string
}

func NewXsdLoader() *XsdLoader {
	return &XsdLoader{
		alreadyLoaded: make(map[string]bool)}
}

func (l *XsdLoader) Load(xsdFilePath string) (io.ReadCloser, error) {
	filePath := l.buildFilePath(xsdFilePath)
	var loadedErr error = nil
	if _, exists := l.alreadyLoaded[filePath]; exists {
		loadedErr = alreadyLoadedErr
	} else {
		l.alreadyLoaded[filePath] = true
	}

	res, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	return res, loadedErr
}

func (l *XsdLoader) buildFilePath(name string) string {
	fullName, err := filepath.Abs(name)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(fullName)
	if err == nil {
		l.curDir = path.Dir(fullName)
		return fullName
	}

	fullName = path.Clean(l.curDir + "/" + name)
	_, err = os.Stat(fullName)
	if err != nil {
		panic(err)
	}

	return fullName
}

func (l *XsdLoader) IsAlreadyLoadedError(e error) bool {
	return e == alreadyLoadedErr
}