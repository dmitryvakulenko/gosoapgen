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

	if !tp.IsSimpleContent {
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
	if !cType.IsSimpleContent {
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
	if cType.IsSimpleContent {
		t.Fatalf("Type should be complex type")
	}

	typeName := "Session"
	if cType.Name != typeName {
		t.Errorf("TypeName elemName should be %q, got %q", typeName, cType.Name)
	}

	ns := "http://xml.amadeus.com/2010/06/Session_v3"
	if cType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, cType.Namespace)
	}

	if cType.BaseType != nil {
		t.Errorf("Type should has no base type")
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
	if cType.IsSimpleContent {
		t.Fatalf("Type should be complex type")
	}

	typeName := "AMA_SecurityHostedUser"
	if cType.Name != typeName {
		t.Errorf("TypeName elemName should be %q, got %q", typeName, cType.Name)
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

	if field.Type.Name != "string" {
		t.Fatalf("Field type name shoud be 'string', %q got", field.Type.Name)
	}
}

func TestComplexTypeWithAttributes(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	if cType.IsSimpleContent {
		t.Fatalf("Type should be complex type")
	}

	typeName := "Session"
	if cType.Name != typeName {
		t.Errorf("TypeName elemName should be %q, got %q", typeName, cType.Name)
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
	if firstType.IsSimpleContent {
		t.Fatalf("Type should be complex type")
	}

	secType := typesList[1]
	if secType.IsSimpleContent {
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
	if cType.IsSimpleContent {
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

	if len(typesList[0].Fields) != 1 {
		t.Fatalf("Fields amount should be 1, %d got", len(typesList[0].Fields))
	}

	f := typesList[0].Fields[0]
	if f.Name != "TPA_Extensions" {
		t.Errorf("Field name shoud be 'TPA_Extensions', %q got", f.Name)
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

func TestUnion(t *testing.T) {
	typesList := parseTypesFrom(t.Name())
	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	if typesList[0].Name != "OwnerSimpleType" {
		t.Errorf("Type name shoud be 'OwnerSimpleType', %q given", typesList[0].Name)
	}

	if typesList[0].BaseTypeName.Name != "string" {
		t.Errorf("Type base type shoud be 'string', %q given", typesList[0].BaseTypeName.Name)
	}
}

func TestSimpleContent(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 3 {
		t.Fatalf("Wrong types amount. 3 expected, %d got", len(typesList))
	}

	cType := typesList[1]
	if !cType.IsSimpleContent {
		t.Errorf("Type should be simple content")
	}

	name := "CompanyNameType"
	if cType.Name != name {
		t.Fatalf("TypeName name should be %q, got %q instead", name, cType.Name)
	}

	if len(cType.Fields) != 3 {
		t.Fatalf("CompanyNameType should has 3 field, %d instead", len(cType.Fields))
	}

	field := cType.Fields[2]

	name = "Value"
	if field.Name != name {
		t.Errorf("Field name should be %q, got %q instead", name, field.Name)
	}

	fType := "string"
	if field.Type.Name != fType {
		t.Errorf("Field type should be %q, got %q instead", fType, field.Type.Name)
	}
}

func TestComplexContent(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 4 {
		t.Fatalf("Wrong types amount. 4 expected, %d got", len(typesList))
	}

	cType := typesList[3]
	if cType.IsSimpleContent {
		t.Errorf("Type should be complex type")
	}

	typeName := "AddressWithModeType"
	if cType.Name != typeName {
		t.Errorf("TypeName name should be %q, got %q instead", typeName, cType.Name)
	}

	if cType.BaseType == nil {
		t.Fatalf("Type should has base type")
	}

	if len(cType.Fields) != 1 {
		t.Fatalf("Fields amount should be 1, got %d instead", len(cType.Fields))
	}

	field := cType.Fields[0]
	fieldName := "Mode"
	if field.Name != fieldName {
		t.Fatalf("TypeName name should be %q, got %q instead", fieldName, field.Name)
	}
}

func TestComplexTypeSimpleContent(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}
}

func TestComplexTypeWithSimpleContent(t *testing.T) {
    typesList := parseTypesFrom(t.Name())

    if len(typesList) != 1 {
        t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
    }
    ct := typesList[0]

    if ct.BaseType != nil {
        t.Fatalf("Base type should be nil")
    }


    if len(ct.Fields) != 3 {
        t.Fatalf("Wrong type fields amount. 3 expected, %d got", len(typesList))
    }
}

func TestChoice(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 2 expected, %d got", len(typesList))
	}

	cType := typesList[1]
	if len(cType.Fields) != 2 {
		t.Errorf("Wrong type fields amount. 2 expected, %d got", len(cType.Fields))
	}
}

func TestElementRefWithType(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 3 {
		t.Fatalf("Wrong types amount. 3 expected, %d got", len(typesList))
	}

	tp := typesList[1]
	if tp.BaseType == nil {
		t.Fatalf("Type should has base type")
	}

	if tp.BaseType.Name != "PointOfSaleType" {
		t.Errorf("Base type name should be 'PointOfSaleType', %q got", tp.BaseType.Name)
	}
}

func TestComplexChoice(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	if len(cType.Fields) != 3 {
		t.Errorf("Wrong type fields amount. 3 expected, %d got", len(cType.Fields))
	}
}

func TestSimpleTypesFolding(t *testing.T) {
    typesList := parseTypesFrom(t.Name())

    if len(typesList) != 3 {
        t.Fatalf("Wrong types amount. 3 expected, %d got", len(typesList))
    }

    ct := typesList[0]
    if ct.Name != "CountryCode" {
        t.Errorf("Type name shoud be 'OwnerSimpleType', %q given", ct.Name)
    }

    if len(ct.Fields) != 0 {
        t.Errorf("Wrong type fields amount. 0 expected, %d got", len(ct.Fields))
    }

    if ct.BaseType == nil {
        t.Fatalf("Type should has base type")
    }

    fields := ct.BaseType.Fields
    if len(fields) != 2 {
        t.Errorf("Wrong type fields amount. 2 expected, %d got", len(fields))
    }

    if fields[1].Name != "Value" {
        t.Errorf("Last field should be Value")
    }

    if fields[1].Type.Name != "string" {
        t.Errorf("Last field type should be string")
    }
}

func TestTwoLevelSimpleContent(t *testing.T) {
    typesList := parseTypesFrom(t.Name())

    if len(typesList) != 3 {
        t.Fatalf("Wrong types amount. 3 expected, %d got", len(typesList))
    }

    ct := typesList[0]
    if len(ct.Fields) != 3 {
        t.Fatalf("Wrong type fields amount. 3 expected, %d got", len(ct.Fields))
    }
}

func parseTypesFrom(name string) []*Type {
	parser := NewParser(&SimpleLoader{})
	parser.Load(name + ".xsd")

	return parser.GetTypes()
}

type SimpleLoader struct{}

func (l *SimpleLoader) Load(path string) (io.ReadCloser, error) {
	file, _ := os.Open("./test_data/" + path)
	return file, nil
}

func (l *SimpleLoader) IsAlreadyLoadedError(e error) bool {
	return false
}
