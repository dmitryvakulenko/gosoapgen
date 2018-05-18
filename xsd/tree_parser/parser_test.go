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
		t.Errorf("Field elemName should be %q, got %q instead", name, tp.Name)
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
		t.Errorf("TypeName elemName should be %q, got %q", typeName, cType.Name)
	}

	if cType.BaseType == nil {
		t.Errorf("Type should has base type, %q exist", cType.BaseType.Name)
	}

	if cType.BaseType.Name != "decimal" {
		t.Errorf("Base type should be 'decimal', %q exist", cType.BaseType.Name)
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
		t.Errorf("TypeName elemName should be %q, got %q", typeName, cType.GoName)
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
		t.Errorf("Field elemName should be 'sequenceNumber', %q instead", field.Name)
	}

	if field.TypeName.Name != "string" {
		t.Errorf("Field type should be 'string' %q instead", field.Type.Name)
	}

	field = cType.Fields[3]
	if field.Name != "TransactionStatusCode" {
		t.Errorf("Field elemName should be 'TransactionStatusCode' %s instead", field.Name)
	}

	if !field.IsAttr {
		t.Errorf("TransactionStatusCode should be attribute")
	}
}

func TestSchemaComplexTypes(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 2 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	if cType.IsSimple {
		t.Fatalf("Type should be complex type")
	}

	typeName := "AMA_SecurityHostedUser"
	if cType.Name != typeName {
		t.Errorf("TypeName elemName should be %q, got %q", typeName, cType.GoName)
	}

	ns := "http://xml.amadeus.com/2010/06/Security_v1"
	if cType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, cType.Namespace)
	}

	if len(cType.Fields) != 4 {
		t.Errorf("TypeName should Has 4 fields, %d getting", len(cType.Fields))
	}

	field := cType.Fields[2]
	if field.Type == nil {
		t.Fatalf("Field should has type")
	}

	if field.Type.Name != "StringLength1to64" {
		t.Fatalf("Field type name shoud be 'StringLength1to64', %q got", field.Type.Name)
	}
}

func TestComplexTypeWithAttributes(t *testing.T) {
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
		t.Errorf("TypeName elemName should be %q, got %q", typeName, cType.GoName)
	}

	ns := "http://xml.amadeus.com/2010/06/Session_v3"
	if cType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, cType.Namespace)
	}

	if len(cType.Fields) != 1 {
		t.Fatalf("Should be 1 fields, %d getting", len(cType.Fields))
	}

	field := cType.Fields[0]
	if field.Name != "TransactionStatusCode" {
		t.Errorf("Field elemName should be 'TransactionStatusCode' %s instead", field.Name)
	}

	if field.TypeName.Name != "NMTOKEN" {
		t.Errorf("Field type should be 'string' %s instead", field.TypeName.Name)
	}

	if !field.IsAttr {
		t.Errorf("TransactionStatusCode should be attribute")
	}
}

func TestInnerComplexTypes(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 3 {
		t.Fatalf("Wrong types amount. 3 expected, %d got", len(typesList))
	}

	firstType := typesList[0]
	if firstType.IsSimple {
		t.Fatalf("Type should be complex type")
	}

	secType := typesList[1]
	if secType.IsSimple {
		t.Fatalf("Type should be complex type")
	}

	typeName := "PNR_AddMultiElements"
	if firstType.Name != typeName {
		t.Errorf("TypeName elemName should be %q, got %q", typeName, firstType.Name)
	}

	ns := "http://xml.amadeus.com/PNRADD_10_1_1A"
	if firstType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, firstType.Namespace)
	}

	typeName = "travellerInfo"
	if secType.Name != typeName {
		t.Errorf("TypeName elemName should be %q, got %q", typeName, secType.Name)
	}

	if secType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, secType.Namespace)
	}

	if len(firstType.Fields) != 1 {
		t.Fatalf("Should be 1 fields, %d getting", len(firstType.Fields))
	}

	field := firstType.Fields[0]
	if field.Name != "travellerInfo" {
		t.Errorf("Field elemName should be 'travellerInfo', %q instead", field.Name)
	}

	if field.Type.Name != "travellerInfo" {
		t.Errorf("Field type should be 'travellerInfo', %q instead", field.Type.Name)
	}

	if len(secType.Fields) != 1 {
		t.Fatalf("Second type fields amount should be 1, got %d instead", len(secType.Fields))
	}

	if secType.Fields[0].Name != "elementManagementPassenger" {
		t.Errorf("Second type elemName shoud be 'ElementManagementPassenger', %q instead", secType.Fields[0].Name)
	}

	thirdType := typesList[2]
	field = thirdType.Fields[0]
	if field.Name != "reference" {
		t.Errorf("Field elemName should be 'reference', %q instead", field.Name)
	}

	if field.Type.Name != "string" {
		t.Errorf("Field type should be 'string', %q instead", field.Type.Name)
	}
}

func TestAttributeGroup(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	// вообще-то, строго говоря, attributeGroup - это не совсем тип, он встраивается
	// но для простоты пусть будет так
	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	if cType.IsSimple {
		t.Fatalf("Type should be complex type")
	}

	name := "CodeType"
	if cType.Name != "CodeType" {
		t.Fatalf("TypeName elemName should be %q, got %q instead", name, cType.Name)
	}

	if len(cType.Fields) != 1 {
		t.Fatalf("Fields amount should be 5, %d instead", len(cType.Fields))
	}

	field := cType.Fields[0]
	if field.Name != "CodeGroup" {
		t.Fatalf("Field elemName should be 'CodeGroup', %q instead", field.Name)
	}

	if !field.IsAttr {
		t.Fatalf("Owner should be attribute")
	}
}

func TestParseElementRef(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 2 expected, %d got", len(typesList))
	}
}


func TestInclude(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 3 {
		t.Fatalf("Wrong types amount. 3 expected, %d got", len(typesList))
	}
}

func TestSimpleTypeAttribute(t *testing.T) {
	typesList := parseTypesFrom(t.Name())
	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}
}

func parseTypesFrom(name string) []*Type {
	parser := NewParser(&SimpleLoader{})
	parser.Parse(name + ".xsd")

	return parser.GenerateTypes()
}

type SimpleLoader struct{}

func (l *SimpleLoader) Load(path string) (io.ReadCloser, error) {
	file, _ := os.Open("./test_data/" + path)
	return file, nil
}

func (l *SimpleLoader) IsAlreadyLoadedError(e error) bool {
	return false
}
