package xsd

import (
	"testing"
	"os"
)

func TestSingleElementSchema(t *testing.T) {
	reader := getReader("./parser_test/1.xsd")
	defer reader.Close()

	p := CreateParser(reader)
	p.Parse()
	types := p.GetTypes()

	if len(types) != 1 {
		t.Errorf("Should be 1 type, %d instead", len(types))
	}

	if types[0].Name != "Session" {
		t.Errorf("Type name should be '%s', got '%s' instead", "Session", types[0].Name)
	}

	if len(types[0].Fields) != 0 {
		t.Errorf("Fields sould be empty")
	}

	ns := "http://www.w3.org/2001/XMLSchema"
	if types[0].Namespace != ns {
		t.Errorf("Namespace should be '%s', '%s' instead", ns, types[0].Namespace)
	}
}


func TestParsingComplexTypeWithAttributes(t *testing.T) {
	reader := getReader("./parser_test/2.xsd")
	defer reader.Close()

	p := CreateParser(reader)
	p.Parse()
	types := p.GetTypes()

	if len(types) != 1 {
		t.Errorf("Should be 1 type, %d instead", len(types))
	}

	if len(types[0].Fields) != 3 {
		t.Errorf("Fields amount sould be 3, %d instead", len(types[0].Fields))
	}

	if types[0].Fields[2].Name != "SecurityToken" {
		t.Errorf("Field name should be 'SecurityToken', %s instead", types[0].Fields[2].Name)
	}
}


func getReader(fileName string) *os.File {
	reader, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	return reader
}