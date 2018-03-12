package translator

import (
	"encoding/xml"
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"os"
	"testing"
)

func TestGetNoTypes(t *testing.T) {
	s := loadXsd("empty.xsd")
	res := Parse(s, s.TargetNamespace)

	if len(res.GetTypes()) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestGenerateSimpleTypes(t *testing.T) {
	s := loadXsd("simpleType.xsd")

	res := Parse(s, s.TargetNamespace)
	typesList := res.GetTypes()

	if len(typesList) != 1 {
		t.Fatalf("Wrong number of types. 1 expected, but got %d", len(typesList))
	}

	curType := typesList[0].(*SimpleType)

	name := "AlphaString_Length1To2"
	if curType.Name != name {
		t.Errorf("Type name should be %q, got %q", name, curType.Name)
	}

	if curType.Type != "string" {
		t.Errorf("Type should be 'string', %s getting", curType.Type)
	}

	ns := "http://xml.amadeus.com/PNRADD_10_1_1A"
	if curType.Namespace != ns {
		t.Errorf("Type should be %q, %q getting", ns, curType.Namespace)
	}
}

//func TestParseElementTypes(t *testing.T) {
//	s := loadXsd("element.xsd")
//	res := Parse(s, s.TargetNamespace)
//
//	typeName := "Session"
//	ns := "http://xml.amadeus.com/2010/06/Session_v3"
//	cTypeInterface, ok := res.cType.find(ns, typeName)
//
//	if !ok {
//		t.Fatalf("Type %q should exists", typeName)
//	}
//
//	cType := cTypeInterface.(*ComplexType)
//	if len(cType.Fields) != 4 {
//		t.Fatalf("Should be 4 fields, %d getting", len(cType.Fields))
//	}
//
//	field := cType.Fields[1]
//	if field.Name != "sequenceNumber" {
//		t.Errorf("Field name should be 'sequenceNumber', %q instead", field.Name)
//	}
//
//	if field.Type != "string" {
//		t.Errorf("Field type should be 'string' %s instead", field.Type)
//	}
//
//	if field.XmlExpr != "sequenceNumber" {
//		t.Errorf("Field xml expression should be 'sequenceNumber', %q instead", field.XmlExpr)
//	}
//
//	if field.Namespace != ns {
//		t.Errorf("Type should be %q, %q getting", ns, cType.Namespace)
//	}
//
//	field = cType.Fields[3]
//	if field.Name != "TransactionStatusCode" {
//		t.Errorf("Field name should be 'TransactionStatusCode' %s instead", field.Name)
//	}
//
//	if field.XmlExpr != "TransactionStatusCode,attr" {
//		t.Errorf("Field xml expr should be 'TransactionStatusCode,attr' %s instead", field.XmlExpr)
//	}
//}
//
//func TestGenerateSchemaComplexTypes(t *testing.T) {
//	s := loadXsd("complexType.xsd")
//	res := Parse(s, s.TargetNamespace)
//
//	typeName := "AMA_SecurityHostedUser"
//	cTypeInterface, ok := res.cType.find("http://xml.amadeus.com/2010/06/Security_v1", typeName)
//
//	if !ok {
//		t.Fatalf("Type %q should exists", typeName)
//	}
//
//	cType := cTypeInterface.(*ComplexType)
//
//	if cType.Name != "AMA_SecurityHostedUser" {
//		t.Fatalf("Type name should be 'Session', '%s' getting", cType.Name)
//	}
//
//	if len(cType.Fields) != 4 {
//		t.Fatalf("Type should has 4 fields, %d getting", len(cType.Fields))
//	}
//
//}
//
//func TestComplexTypeWithAttributes(t *testing.T) {
//	s := loadXsd("attribute.xsd")
//	res := Parse(s, s.TargetNamespace)
//
//	typeName := "Session"
//	cTypeInterface, ok := res.cType.find("http://xml.amadeus.com/2010/06/Session_v3", typeName)
//
//	if !ok {
//		t.Fatalf("Type %q should exists", typeName)
//	}
//
//	cType := cTypeInterface.(*ComplexType)
//
//	if len(cType.Fields) != 1 {
//		t.Fatalf("Should be 1 fields, %d getting", len(cType.Fields))
//	}
//
//	field := cType.Fields[0]
//	if field.Name != "TransactionStatusCode" {
//		t.Errorf("Field name should be 'TransactionStatusCode' %s instead", field.Name)
//	}
//
//	if field.Type != "NMTOKEN" {
//		t.Errorf("Field type should be 'NMTOKEN' %s instead", field.Type)
//	}
//
//	if field.XmlExpr != "TransactionStatusCode,attr" {
//		t.Errorf("Field xml name should be 'sequenceNumber' %s instead", field.XmlExpr)
//	}
//}
//
//func TestInnerComplexTypes(t *testing.T) {
//	s := loadXsd("innerComplexType.xsd")
//	res := Parse(s, s.TargetNamespace)
//
//	var (
//		firstTypeI, secTypeI interface{}
//		ok bool
//	)
//	if firstTypeI, ok = res.cType.find("http://xml.amadeus.com/PNRADD_10_1_1A", "PNR_AddMultiElements"); !ok {
//		t.Fatalf("Type %q should exists", "PNR_AddMultiElements")
//	}
//
//	if secTypeI, ok = res.cType.find("http://xml.amadeus.com/PNRADD_10_1_1A", "travellerInfo"); !ok {
//		t.Fatalf("Type %q should exists", "travellerInfo")
//	}
//
//	firstType := firstTypeI.(*ComplexType)
//	secType := secTypeI.(*ComplexType)
//
//	if len(firstType.Fields) != 1 {
//		t.Fatalf("Should be 1 fields, %d getting", len(firstType.Fields))
//	}
//
//	field := firstType.Fields[0]
//	if field.Name != "travellerInfo" {
//		t.Errorf("Field name should be 'travellerInfo', %q instead", field.Name)
//	}
//
//	if field.Type != "travellerInfo" {
//		t.Errorf("Field type should be 'travellerInfo', %q instead", field.Type)
//	}
//
//	if field.XmlExpr != "travellerInfo" {
//		t.Errorf("Field xml name should be 'travellerInfo' %s instead", field.XmlExpr)
//	}
//
//	if len(secType.Fields) != 1 {
//		t.Fatalf("Second type fields amount should be 1, got %d instead", len(secType.Fields))
//	}
//
//	if secType.Fields[0].Name != "elementManagementPassenger" {
//		t.Errorf("Second type name shoud be 'ElementManagementPassenger', %q instead", secType.Fields[0].Name)
//	}
//
//	if secType.Fields[0].XmlExpr != "elementManagementPassenger" {
//		t.Errorf("Second type name shoud be 'elementManagementPassenger', %q instead", secType.Fields[0].XmlExpr)
//	}
//}
//
//func TestAttributeGroup(t *testing.T) {
//	s := loadXsd("attributeGroup.xsd")
//	res := Parse(s, s.TargetNamespace)
//
//	cTypeI, ok := res.cType.find("http://xml.amadeus.com/2010/06/Types_v1", "CodeType")
//	if !ok {
//		t.Fatalf("Type %q should exists", "CodeType")
//	}
//
//	cType := cTypeI.(*ComplexType)
//
//	if cType.Name != "CodeType" {
//		t.Fatalf("Type name should be CodeType, %q instead", cType.Name)
//	}
//
//	if len(cType.Fields) != 5 {
//		t.Fatalf("Fields amount should be 5, %d instead", len(cType.Fields))
//	}
//
//	field := cType.Fields[1]
//	if field.Name != "Owner" {
//		t.Fatalf("Field name should be 'Owner', %q instead", field.Name)
//	}
//
//	if field.XmlExpr != "Owner,attr" {
//		t.Fatalf("Field xml expression should be 'Owner,attr', %q instead", field.XmlExpr)
//	}
//}
//
//
//func TestSimpleContent(t *testing.T) {
//	s := loadXsd("simpleContent.xsd")
//	ns := "namespace"
//	res := Parse(s, ns)
//
//	name := "StringLength0to128"
//	_, ok := res.sType.find(ns, name)
//	if !ok {
//		t.Fatalf("Simple type %q should exists", name)
//	}
//
//	name = "CompanyNameType"
//	cTypeI, ok := res.sType.find(ns, name)
//	if !ok {
//		t.Fatalf("Complex type %q should exists", name)
//	}
//
//	cType := cTypeI.(*ComplexType)
//	if len(cType.Fields) != 1 {
//		t.Errorf("CompanyNameType should has 1 field, %d instead", len(cType.Fields))
//	}
//}
//
//
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
