package generator

import (
	"testing"
	"gosoapgen/xsd"
)

func TestGetNoTypes(t *testing.T) {
	res := GenerateTypes([]xsd.Schema{{}})

	if len(res) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestSingleEmptyStructure(t *testing.T) {
	s := xsd.Schema{}
	s.Element = append(s.Element, xsd.Element{Name: "Session"})

	res := GenerateTypes([]xsd.Schema{s})
	if len(res) != 1 {
		t.Errorf("Should be 1 type, %d getting", len(res))
	}
}


func TestSimpleStructureHasSeveralFields(t *testing.T) {
	elem := xsd.Element{Name: "Session"}
	elem.ComplexType.Sequence.Element = append(elem.ComplexType.Sequence.Element, xsd.Element{Name: "SessionId", Type: "xs:string"})
	elem.ComplexType.Sequence.Element = append(elem.ComplexType.Sequence.Element, xsd.Element{Name: "sequenceNumber", Type: "xs:string"})
	elem.ComplexType.Sequence.Element = append(elem.ComplexType.Sequence.Element, xsd.Element{Name: "SecurityToken", Type: "xs:string"})

	s := xsd.Schema{}
	s.Element = append(s.Element, elem)

	res := GenerateTypes([]xsd.Schema{s})
	if len(res) != 1 {
		t.Errorf("Should be 1 type, %d getting", len(res))
	}

	if len(res[0].Fields) != 3 {
		t.Fatalf("Should be 3 fields, %d getting", len(res[0].Fields))
	}

	field := res[0].Fields[1]
	if field.Name != "SequenceNumber" {
		t.Errorf("Field name should be 'SequenceNumber' %s instead", field.Name)
	}

	if field.Type != "string" {
		t.Errorf("Field type should be 'string' %s instead", field.Type)
	}

	if field.XmlExpr != "sequenceNumber" {
		t.Errorf("Field xml name should be 'sequenceNumber' %s instead", field.XmlExpr)
	}
}


func TestComplexTypeWithAttributes(t *testing.T) {
	elem := xsd.Element{Name: "Session"}
	elem.ComplexType.Attribute = append(elem.ComplexType.Attribute, xsd.Attribute{Name: "TransactionStatusCode", Type: "xs:NMTOKEN"})

	s := xsd.Schema{}
	s.Element = append(s.Element, elem)

	res := GenerateTypes([]xsd.Schema{s})

	if len(res[0].Fields) != 1 {
		t.Fatalf("Should be 1 fields, %d getting", len(res[0].Fields))
	}

	field := res[0].Fields[0]
	if field.Name != "TransactionStatusCode" {
		t.Errorf("Field name should be 'TransactionStatusCode' %s instead", field.Name)
	}

	if field.Type != "NMTOKEN" {
		t.Errorf("Field type should be 'NMTOKEN' %s instead", field.Type)
	}

	if field.XmlExpr != "TransactionStatusCode,attr" {
		t.Errorf("Field xml name should be 'sequenceNumber' %s instead", field.XmlExpr)
	}
}