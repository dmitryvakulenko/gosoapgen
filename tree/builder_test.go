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

func TestElementComplexType(t *testing.T) {
	builder := createAndDecode("element.xsd")

	typesList := builder.getTypes()
	if len(typesList) != 1 {
		t.Fatalf("Types amount should be 1, got %d", len(typesList))
	}
	curType := typesList[0]

	typeName := "Session"
	if typeName != curType.Name {
		t.Errorf("Type name should be %q, got %q", typeName, curType.Name)
	}

	ns := "http://xml.amadeus.com/2010/06/Session_v3"
	if ns != curType.Namespace {
		t.Errorf("Type namespace should be %q, got %q", ns, curType.Namespace)
	}

	if len(curType.Fields) != 4 {
		t.Fatalf("Fields amount should be 4, got %d", len(curType.Fields))
	}
}

func createAndDecode(fileName string) Builder {
	b := NewBuilder()
	b.Build("./testdata/" + fileName)

	return b
}

