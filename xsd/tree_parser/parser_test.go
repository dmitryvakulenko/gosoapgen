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

	tp := typesList[0]

	if !tp.IsSimple {
		t.Fatalf("Type should be complex type")
	}

	name := "AlphaString_Length1To2"
	if name != tp.Name {
		t.Errorf("Field name should be %q, got %q instead", name, tp.Name)
	}

	if tp.BaseTypeName.Name != "string" {
		t.Errorf("Field should be string, got %q instead", tp.Name)
	}
}

func TestSimpleElements(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	if !cType.IsSimple {
		t.Fatalf("Type should be simple")
	}

	typeName := "minRange"
	if cType.Name != typeName {
		t.Errorf("TypeName name should be %q, got %q", typeName, cType.GoName)
	}


	ns := "http://xml.amadeus.com/2010/06/Types_v1"
	if cType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, cType.Namespace)
	}

	if cType.BaseTypeName.Name != "decimal" {
		t.Errorf("Type should be decimal, got %q", cType.BaseTypeName.Name)
	}
}

func TestComplexType(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	if cType.IsSimple {
		t.Fatalf("Type should be complex type")
	}

	typeName := "Session"
	if cType.Name != typeName {
		t.Errorf("TypeName name should be %q, got %q", typeName, cType.GoName)
	}

	ns := "http://xml.amadeus.com/2010/06/Session_v3"
	if cType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, cType.Namespace)
	}

	if len(cType.Fields) != 4 {
		t.Fatalf("Should be 4 fields, %d getting", len(cType.Fields))
	}

	field := cType.Fields[1]
	if field.Name != "sequenceNumber" {
		t.Errorf("Field name should be 'sequenceNumber', %q instead", field.Name)
	}

	if field.TypeName.Name != "string" {
		t.Errorf("Field type should be 'string' %q instead", field.Type.Name)
	}

	field = cType.Fields[3]
	if field.Name != "TransactionStatusCode" {
		t.Errorf("Field name should be 'TransactionStatusCode' %s instead", field.Name)
	}

	if !field.IsAttr {
		t.Errorf("TransactionStatusCode should be attribute")
	}
}


func parseTypesFrom(name string) []*Type {
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

