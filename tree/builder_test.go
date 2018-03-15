package tree

import (
	"testing"
)

func TestNoTypes(t *testing.T) {
	builder := createAndDecode("empty.xsd")

	typesList := builder.getTypes()
	if len(typesList) != 0 {
		t.Fatalf("Types amount should be 0, got %d", len(typesList))
	}
}

func TestSimpleType(t *testing.T) {
	builder := createAndDecode("simpleType.xsd")

	typesList := builder.getTypes()
	if len(typesList) != 0 {
		t.Fatalf("Types amount should be 0, got %d", len(typesList))
	}
}

func createAndDecode(fileName string) Builder {
	b := NewBuilder()
	b.Build("./testdata/" + fileName)

	return b
}

