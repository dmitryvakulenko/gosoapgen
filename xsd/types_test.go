package xsd

import (
	"testing"
)

func TestSingleElementParsing(t *testing.T) {
	s := parseSchema("./types_test/1.xsd")

	if len(s.Element) != 1 {
		t.Fatalf("Should be 1 type, %d instead", len(s.Element))
	}

	if s.Element[0].Name != "Session" {
		t.Errorf("TypeName name should be '%s', got '%s' instead", "Session", s.Element[0].Name)
	}

	if len(s.Element[0].ComplexType.Sequence.Element) != 0 {
		t.Errorf("Fields sould be empty")
	}
}

func TestParsingComplexTypeWithAttributes(t *testing.T) {
	s := parseSchema("./types_test/2.xsd")

	if len(s.Element) != 1 {
		t.Fatalf("Should be 1 type, %d instead", len(s.Element))
	}

	if s.Element[0].Name != "Session" {
		t.Errorf("TypeName name should be '%s', got '%s' instead", "Session", s.Element[0].Name)
	}

	if len(s.Element[0].ComplexType.Sequence.Element) != 3 {
		t.Fatalf("Fields amount sould be 3, %d instead", len(s.Element[0].ComplexType.Sequence.Element))
	}

	if s.Element[0].ComplexType.Sequence.Element[2].Name != "SecurityToken" {
		t.Errorf("Field name should be 'SecurityToken', %s instead", s.Element[0].ComplexType.Sequence.Element[2].Name)
	}

	if len(s.Element[0].ComplexType.Attribute) != 1 {
		t.Fatalf("Attributes amount sould be 1, %d instead", len(s.Element[0].ComplexType.Attribute))
	}

	if s.Element[0].ComplexType.Attribute[0].Name != "TransactionStatusCode" {
		t.Fatalf("Attribute name should be TransactionStatusCode, %d instead", len(s.Element[0].ComplexType.Attribute[0].Name))
	}
}


func TestParsingAdditionTypes(t *testing.T) {
	s := parseSchema("./types_test/3.xsd")

	if len(s.ComplexType) != 2 {
		t.Fatalf("Comples types amount sould be 2, %d instead", len(s.ComplexType))
	}

	if s.ComplexType[1].Name != "AvailabilityOptionsType" {
		t.Fatalf("Complex type name should be 'AvailabilityOptionsType', %s instead", s.ComplexType[1].Name)
	}

	if len(s.SimpleType) != 2 {
		t.Fatalf("Simple types amount sould be 2, %d instead", len(s.SimpleType))
	}
}

func TestImportInclude(t *testing.T) {
	s := ParseSchema("./types_test/4.xsd")

	if len(s) != 3 {
		t.Fatalf("Should be imported 3 schemas, %d instead", len(s))
	}
}

func parseSchema(fileName string) *Schema {
	return ParseSchema(fileName)[0]
}