package translator

import (
	"testing"
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"os"
	"encoding/xml"
)

func TestGetNoTypes(t *testing.T) {
	s := loadXsd("empty.xsd")
	res := Parse(s)

	if len(res.cType) != 0 {
		t.Errorf("Should be no types")
	}

	if len(res.sType) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestGenerateSimpleTypes(t *testing.T) {
	s := loadXsd("simpleType.xsd")

	res := Parse(s)
	if len(res.cType) != 0 {
		t.Fatalf("Complex types should be empty")
	}

	if len(res.sType) != 1 {
		t.Fatalf("Schema should has one simple type")
	}

	typeName := "AlphaString_Length1To2"
	if res.sType[0].Name != typeName {
		t.Errorf("Type name should be %q, %q getting", typeName, res.sType[0].Name)
	}

	if res.sType[0].Type != "string" {
		t.Errorf("Type should be 'string', %s getting", res.sType[0].Type)
	}

	ns := "http://xml.amadeus.com/PNRADD_10_1_1A"
	if res.sType[0].Namespace != ns {
		t.Errorf("Type should be %q, %q getting", ns, res.sType[0].Namespace)
	}
}


func TestParseElementTypes(t *testing.T) {
	s := loadXsd("element.xsd")
	res := Parse(s)

	if len(res.sType) != 0 {
		t.Fatalf("Schema should not contain simple types")
	}

	if len(res.cType) != 1 {
		t.Fatalf("Should be 1 type, %d getting", len(res.sType))
	}

	cType := res.cType[0]
	if len(cType.Fields) != 4 {
		t.Fatalf("Should be 4 fields, %d getting", len(cType.Fields))
	}

	field := cType.Fields[1]
	if field.Name != "SequenceNumber" {
		t.Errorf("Field name should be 'SequenceNumber', %q instead", field.Name)
	}

	if field.Type != "string" {
		t.Errorf("Field type should be 'string' %s instead", field.Type)
	}

	if field.XmlExpr != "sequenceNumber" {
		t.Errorf("Field xml expression should be 'sequenceNumber', %q instead", field.XmlExpr)
	}

	ns := "http://xml.amadeus.com/2010/06/Session_v3"
	if field.Namespace != ns {
		t.Errorf("Type should be %q, %q getting", ns, res.sType[0].Namespace)
	}

	field = cType.Fields[3]
	if field.Name != "TransactionStatusCode" {
		t.Errorf("Field name should be 'TransactionStatusCode' %s instead", field.Name)
	}

	if field.XmlExpr != "TransactionStatusCode,attr" {
		t.Errorf("Field xml expr should be 'TransactionStatusCode,attr' %s instead", field.XmlExpr)
	}
}


func TestGenerateSchemaComplexTypes(t *testing.T) {
	s := loadXsd("complexType.xsd")
	res := Parse(s)

	if len(res.sType) != 0 {
		t.Fatalf("Schema should not contain simple types")
	}

	if len(res.cType) != 1 {
		t.Fatalf("Should be 1 complex type, %d getting", len(res.cType))
	}

	if res.cType[0].Name != "AMA_SecurityHostedUser" {
		t.Fatalf("Type name should be 'Session', '%s' getting", res.cType[0].Name)
	}

	if len(res.cType[0].Fields) != 4 {
		t.Fatalf("Type should has 4 fields, %d getting", len(res.cType[0].Fields))
	}

}


func TestComplexTypeWithAttributes(t *testing.T) {
	s := loadXsd("attribute.xsd")
	res := Parse(s)

	if len(res.cType[0].Fields) != 1 {
		t.Fatalf("Should be 1 fields, %d getting", len(res.cType[0].Fields))
	}

	field := res.cType[0].Fields[0]
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
	s := loadXsd("innerComplexType.xsd")
	res := Parse(s)

	if len(res.cType) != 2 {
		t.Fatalf("Types amount should be 2, %d instead", len(res.cType))
	}

	if len(res.cType[0].Fields) != 1 {
		t.Fatalf("Should be 1 fields, %d getting", len(res.cType[0].Fields))
	}

	field := res.cType[0].Fields[0]
	if field.Name != "TravellerInfo" {
		t.Errorf("Field name should be 'TravellerInfo', %q instead", field.Name)
	}

	if field.Type != "TravellerInfo" {
		t.Errorf("Field type should be 'TravellerInfo', %q instead", field.Type)
	}

	if field.XmlExpr != "travellerInfo" {
		t.Errorf("Field xml name should be 'travellerInfo' %s instead", field.XmlExpr)
	}

	if res.cType[1].Name != "TravellerInfo" {
		t.Errorf("Second type name shoud be 'TravellerInfo' %s instead", res.cType[1].Name)
	}


	if len(res.cType[1].Fields) != 1 {
		t.Fatalf("Second type fields amount should be 1, got %d instead", len(res.cType[1].Fields))
	}

	if res.cType[1].Fields[0].Name != "ElementManagementPassenger" {
		t.Errorf("Second type name shoud be 'ElementManagementPassenger', %q instead", res.cType[1].Fields[0].Name)
	}

	if res.cType[1].Fields[0].XmlExpr != "elementManagementPassenger" {
		t.Errorf("Second type name shoud be 'elementManagementPassenger', %q instead", res.cType[1].Fields[0].XmlExpr)
	}
}


func TestAttributeGroup(t *testing.T) {
	s := loadXsd("attributeGroup.xsd")
	res := Parse(s)

	if len(res.cType) != 1 {
		t.Fatalf("Types amount should be 1, %d instead", len(res.cType))
	}

	if res.cType[0].Name != "CodeType" {
		t.Fatalf("Type name should be CodeType, %q instead", res.cType[0].Name)
	}

	if len(res.cType[0].Fields) != 5 {
		t.Fatalf("Embed types amount should be 5, %d instead", len(res.cType[0].Fields))
	}

	field := res.cType[0].Fields[1]
	if field.Name != "Owner" {
		t.Fatalf("Field name should be 'Owner', %q instead", field.Name)
	}

	if field.XmlExpr != "Owner,attr" {
		t.Fatalf("Field xml expression should be 'Owner,attr', %q instead", field.XmlExpr)
	}
}


func TestSimpleContent(t *testing.T) {
	s := loadXsd("simpleContent.xsd")
	res := Parse(s)

	if len(res.sType) != 1 {
		t.Fatalf("Schema should has 1 simple type, %d ", len(res.sType))
	}

	if len(res.cType) != 1 {
		t.Fatalf("Schema should has 1 complex type, %d ", len(res.cType))
	}

	name := "StringLength0to128"
	if res.sType[0].Name != name {
		t.Errorf("Simple type name should be %q, %q instead", name, res.cType[0].Name)
	}

	name = "CompanyNameType"
	if res.cType[0].Name != name {
		t.Errorf("Complex type name should be %q, %q instead", name, res.cType[0].Name)
	}

	if len(res.cType[0].Fields) != 1 {
		t.Errorf("CompanyNameType should has 1 field, %d instead", len(res.cType[0].Fields))
	}
}


func loadXsd(name string) *xsd.Schema {
	reader, err := os.Open("./translator/schema_test/" + name)
	defer reader.Close()

	if err != nil {
		panic(err)
	}

	s := xsd.Schema{}
	err = xml.NewDecoder(reader).Decode(&s)
	if err != nil {
		panic(err)
	}

	return &s
}