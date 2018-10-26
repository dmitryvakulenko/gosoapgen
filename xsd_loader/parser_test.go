package xsd_loader

import (
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

func parseTypesFrom(name string) *Schema {
	parser := NewParser(&SimpleLoader{})

	return parser.Parse(name)
}

type SimpleLoader struct{}

func (l *SimpleLoader) Load(path string) (io.ReadCloser, error) {
	file, _ := os.Open("./testdata/" + path)
	return file, nil
}

func (l *SimpleLoader) IsAlreadyLoadedError(e error) bool {
	return false
}
