package translator

import (
	"testing"
	"github.com/dmitryvakulenko/gosoapgen/xsd"
)

func TestGetNoTypes(t *testing.T) {
	schemas := []*xsd.Schema{{}}
	res := GenerateTypes(schemas)

	if len(res) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestGenerateSimpleTypes(t *testing.T) {
	s := xsd.Schema{}
	s.SimpleType = append(s.SimpleType, &xsd.SimpleType{Name: "AlphaNumericString_Length1To3", Restriction: xsd.Restriction{Base: "xs:string"} })

	schemas := []*xsd.Schema{&s}
	res := GenerateTypes(schemas)
	if len(res) != 1 {
		t.Fatalf("Should be 1 type, %d getting", len(res))
	}

	if res[0].Name != "AlphaNumericString_Length1To3" {
		t.Errorf("Type name should be 'AlphaNumericString_Length1To3', %s getting", res[0].Name)
	}

	if res[0].Type != "string" {
		t.Errorf("Type should be 'string', %s getting", res[0].Type)
	}
}

func TestSimpleStructureHasSeveralFields(t *testing.T) {
	elem := xsd.Element{Name: "Session"}
	elem.ComplexType = &xsd.ComplexType{}
	elem.ComplexType.Sequence = &xsd.Sequence{}
	elem.ComplexType.Sequence.Element = append(elem.ComplexType.Sequence.Element, &xsd.Element{Name: "SessionId", Type: "xs:string"})
	elem.ComplexType.Sequence.Element = append(elem.ComplexType.Sequence.Element, &xsd.Element{Name: "sequenceNumber", Type: "xs:string"})
	elem.ComplexType.Sequence.Element = append(elem.ComplexType.Sequence.Element, &xsd.Element{Name: "SecurityToken", Type: "xs:string"})

	s := xsd.Schema{}
	s.Element = append(s.Element, &elem)

	schemas := []*xsd.Schema{&s}
	res := GenerateTypes(schemas)
	if len(res) != 1 {
		t.Fatalf("Should be 1 type, %d getting", len(res))
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


func TestGenerateSchemaComplexTypes(t *testing.T) {
	elem := xsd.ComplexType{Name: "Session"}
	elem.Sequence = &xsd.Sequence{}
	elem.Sequence.Element = append(elem.Sequence.Element, &xsd.Element{Name: "SessionId", Type: "xs:string"})
	elem.Sequence.Element = append(elem.Sequence.Element, &xsd.Element{Name: "sequenceNumber", Type: "xs:string"})
	elem.Sequence.Element = append(elem.Sequence.Element, &xsd.Element{Name: "SecurityToken", Type: "xs:string"})

	s := xsd.Schema{}
	s.ComplexType = append(s.ComplexType, &elem)

	schemas := []*xsd.Schema{&s}
	res := GenerateTypes(schemas)
	if len(res) != 1 {
		t.Fatalf("Should be 2 type, %d getting", len(res))
	}

	if res[0].Name != "Session" {
		t.Fatalf("Type name should be 'Session', '%s' getting", res[0].Name)
	}

	if len(res[0].Fields) != 3 {
		t.Fatalf("Type should has 3 fields, %d getting", len(res[0].Fields))
	}

}


func TestComplexTypeWithAttributes(t *testing.T) {
	elem := xsd.Element{Name: "Session"}
	elem.ComplexType = &xsd.ComplexType{}
	elem.ComplexType.Attribute = append(elem.ComplexType.Attribute, &xsd.Attribute{Name: "TransactionStatusCode", Type: "xs:NMTOKEN"})

	s := xsd.Schema{}
	s.Element = append(s.Element, &elem)

	schemas := []*xsd.Schema{&s}
	res := GenerateTypes(schemas)

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


func TestInnerComplexTypes(t *testing.T) {
	innerElem := xsd.Element{}
	innerElem.Name = "innerElement"
	innerElem.ComplexType = &xsd.ComplexType{}
	innerElem.ComplexType.Sequence = &xsd.Sequence{}
	innerElem.ComplexType.Sequence.Element = append(innerElem.ComplexType.Sequence.Element, &xsd.Element{Name: "innerField", Type: "xs:string"})

	elem := xsd.Element{Name: "Session"}
	elem.ComplexType = &xsd.ComplexType{}
	elem.ComplexType.Sequence = &xsd.Sequence{}
	elem.ComplexType.Sequence.Element = append(elem.ComplexType.Sequence.Element, &innerElem)

	s := xsd.Schema{}
	s.Element = append(s.Element, &elem)

	schemas := []*xsd.Schema{&s}
	res := GenerateTypes(schemas)

	if len(res) != 2 {
		t.Fatalf("Types amount should be 2, %d instead", len(res))
	}

	if len(res[0].Fields) != 1 {
		t.Fatalf("Should be 1 fields, %d getting", len(res[0].Fields))
	}

	field := res[0].Fields[0]
	if field.Name != "InnerElement" {
		t.Errorf("Field name should be 'InnerElement' %s instead", field.Name)
	}

	if field.Type != "InnerElement" {
		t.Errorf("Field type should be 'InnerElement' %s instead", field.Type)
	}

	if field.XmlExpr != "innerElement" {
		t.Errorf("Field xml name should be 'innerElement' %s instead", field.XmlExpr)
	}

	if res[1].Name != "InnerElement" {
		t.Errorf("Second type name shoud be 'InnerElement' %s instead", res[1].Name)
	}

	if len(res[1].Fields) != 1 {
		t.Fatalf("Second type fields amount should be 2, got %d instead", len(res[1].Fields))
	}

	if res[1].Fields[0].Name != "InnerField" {
		t.Errorf("Second type name shoud be 'InnerField' %s instead", res[1].Fields[0].Name)
	}

	if res[1].Fields[0].XmlExpr != "innerField" {
		t.Errorf("Second type name shoud be 'innerField' %s instead", res[1].Fields[0].XmlExpr)
	}
}

func TestAttributeGroup(t *testing.T) {
	group := xsd.AttributeGroup{}
	group.Name = "attributeGroup"
	group.Attribute = append(group.Attribute, &xsd.Attribute{Name: "innerAttribute", Type: "xs:string"})

	elem := xsd.ComplexType{Name: "Session"}
	elem.AttributeGroup = append(elem.AttributeGroup, &group)

	s := xsd.Schema{}
	s.ComplexType = append(s.ComplexType, &elem)

	schemas := []*xsd.Schema{&s}
	res := GenerateTypes(schemas)

	if len(res) != 1 {
		t.Fatalf("Types amount should be 1, %d instead", len(res))
	}

	if len(res[0].Embed) != 1 {
		t.Fatalf("Embed types amount should be 1")
	}
}