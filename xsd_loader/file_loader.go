package xsd_loader

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

type FileResolver struct {
	alreadyLoaded map[string]bool
	curDir        string
}

func NewFileResolver() *FileResolver {
	return &FileResolver{
		alreadyLoaded: make(map[string]bool)}
}

func (l *FileResolver) Load(xsdFilePath string) (io.ReadCloser, error) {
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

func (l *FileResolver) buildFilePath(name string) string {
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

func (l *FileResolver) IsAlreadyLoadedError(e error) bool {
	return e == alreadyLoadedErr
}