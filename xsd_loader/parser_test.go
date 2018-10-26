package xsd_loader

import (
	"io"
	"os"
)



func parseTypesFrom(name string) []*Type {
	parser := NewParser(&SimpleLoader{})
	parser.Load(name + ".xsd")

	return parser.GetTypes()
}

type SimpleLoader struct{}

func (l *SimpleLoader) Load(path string) (io.ReadCloser, error) {
	file, _ := os.Open("./test_data/" + path)
	return file, nil
}

func (l *SimpleLoader) IsAlreadyLoadedError(e error) bool {
	return false
}
