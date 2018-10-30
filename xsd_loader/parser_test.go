package xsd_loader

import (
	"github.com/dmitryvakulenko/gosoapgen/xsd_loader/tree"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestEmptySchema(t *testing.T) {
	schema := parseTypesFrom(t.Name())

	assert.Empty(t, schema.Elements)
	assert.Empty(t, schema.Types)
	assert.Empty(t, schema.Attributes)
	assert.Empty(t, schema.AttributeGroups)
}


func TestSimpleTypes(t *testing.T) {
	schema := parseTypesFrom(t.Name())
	assert.Len(t, schema.Elements, 1)
	assert.Len(t, schema.Types, 1)
	assert.Equal(t, schema.Elements[0].Type, schema.Types[0])
	assert.Equal(t, schema.Types[0].BaseType.Name, )
}

func TestSimpleElements(t *testing.T) {

}

func parseTypesFrom(name string) *Schema {
	parser := NewParser(tree.NewLoader(&SimpleResolver{}))
	return parser.Parse(name)
}

type SimpleResolver struct{}

func (l *SimpleResolver) Load(path string) (io.ReadCloser, error) {
	file, _ := os.Open("./testdata/" + path + ".xsd")
	return file, nil
}

func (l *SimpleResolver) IsAlreadyLoadedError(e error) bool {
	return false
}
