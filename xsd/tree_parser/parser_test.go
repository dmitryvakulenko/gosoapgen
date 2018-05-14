package tree_parser

import (
	"testing"
	"io/ioutil"
)

func TestEmptySchema(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 0 {
		t.Errorf("Should be no types")
	}
}


func parseTypesFrom(name string) []*NamedType {
	parser := NewParser(&SimpleLoader{})
	parser.Parse(name + ".xsd")

	return parser.GetTypes()
}


type SimpleLoader struct {}

func (l *SimpleLoader) Load(path string) ([]byte, error) {
	data, _ := ioutil.ReadFile("./test_data/" + path)
	return data, nil
}

func (l *SimpleLoader) IsAlreadyLoadedError(e error) bool {
	return false
}

