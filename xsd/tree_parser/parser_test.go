package tree_parser

import (
	"testing"
	"os"
	"io"
)

func TestEmptySchema(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 0 {
		t.Errorf("Should be no types")
	}
}


func TestSimpleTypes(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong number of types. 1 expected, but got %d", len(typesList))
	}

	tp, ok := typesList[0].(*SimpleType)
	if !ok {
		t.Fatalf("Type should be *SimpleType")
	}


	name := "AlphaString_Length1To2"
	if name != tp.Name {
		t.Fatalf("Type name should be %q, got %q instead", name, tp.Name)
	}

}


func parseTypesFrom(name string) []NamedType {
	parser := NewParser(&SimpleLoader{})
	parser.Parse(name + ".xsd")

	return parser.GetTypes()
}


type SimpleLoader struct {}

func (l *SimpleLoader) Load(path string) (io.ReadCloser, error) {
	file, _ := os.Open("./test_data/" + path)
	return file, nil
}

func (l *SimpleLoader) IsAlreadyLoadedError(e error) bool {
	return false
}

