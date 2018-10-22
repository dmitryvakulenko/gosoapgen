package xsd_parser

import (
	"testing"
	"os"
	"io"
    "github.com/stretchr/testify/assert"
)

func TestEmptySchema(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestSimpleTypes(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	assert.Len(t, typesList, 2)
	assert.Equal(t, "Test", typesList[0].Name.Local)
	assert.Equal(t, typesList[1], typesList[0].BaseType)
	assert.Equal(t, "AlphaString_Length1To2", typesList[1].Name.Local)
	assert.Len(t, typesList[1].Fields, 0)
	assert.Equal(t, "string", typesList[1].BaseType.Local)
}

func TestSimpleElements(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	assert.Len(t, typesList, 1)
	assert.Equal(t, "minRange", typesList[0].Name.Local)
	assert.Equal(t, "http://xml.amadeus.com/2010/06/Types_v1", typesList[0].Name.Space)
	assert.Equal(t, "decimal", typesList[0].BaseType.BaseType.Local)
}

func TestComplexType(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	assert.Len(t, typesList, 1)
	assert.Equal(t, "Session", typesList[0].Name.Local)
	assert.Equal(t, "http://xml.amadeus.com/2010/06/Session_v3", typesList[0].Name.Space)
	assert.Len(t, typesList[0].Fields, 4)

	field := typesList[0].Fields[1]
	assert.Equal(t, "sequenceNumber", field.Name)
	assert.Equal(t, "string", field.Type.Local)

	field = typesList[0].Fields[3]
	assert.Equal(t, "TransactionStatusCode", field.Name)
	assert.True(t, field.IsAttr)
}

func TestSchemaComplexTypes(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	name := "AMA_SecurityHostedUser"
	if cType.Local != name {
		t.Errorf("TypeName elemName should be %q, got %q", name, cType.Local)
	}

	ns := "http://xml.amadeus.com/2010/06/Security_v1"
	if cType.Space != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, cType.Space)
	}

	if len(cType.Fields) != 5 {
		t.Errorf("TypeName should Has 5 fields, %d getting", len(cType.Fields))
	}

	field := cType.Fields[3]
	if field.Type == nil {
		t.Fatalf("Field should has type")
	}

	if field.Type.Local != "string" {
		t.Fatalf("Field type name shoud be 'string', %q got", field.Type.Local)
	}
}

func TestComplexTypeWithAttributes(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]

	name := "Session"
	if cType.Local != name {
		t.Errorf("TypeName elemName should be %q, got %q", name, cType.Name)
	}

	ns := "http://xml.amadeus.com/2010/06/Session_v3"
	if cType.Space != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, cType.Space)
	}

	if len(cType.Fields) != 2 {
		t.Fatalf("Should be 2 fields, %d getting", len(cType.Fields))
	}

	field := cType.Fields[1]
	if field.Name != "TransactionStatusCode" {
		t.Errorf("Field elemName should be 'TransactionStatusCode' %s instead", field.Name)
	}

	if field.Type.Local != "NMTOKEN" {
		t.Errorf("Field type should be 'string' %s instead", field.Type.Local)
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

	firstType, secType := typesList[0], typesList[1]

	name := "PNR_AddMultiElements"
	if firstType.Local != name {
		t.Errorf("TypeName elemName should be %q, got %q", name, firstType.Local)
	}

	ns := "http://xml.amadeus.com/PNRADD_10_1_1A"
	if firstType.Space != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, firstType.Space)
	}

	name = "travellerInfo"
	if secType.Local != name {
		t.Errorf("TypeName elemName should be %q, got %q", name, secType.Local)
	}

	if secType.Space != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, secType.Space)
	}

	if len(firstType.Fields) != 2 {
		t.Fatalf("Should be 2 fields, %d getting", len(firstType.Fields))
	}

	field := firstType.Fields[1]
	if field.Name != "travellerInfo" {
		t.Errorf("Field elemName should be 'travellerInfo', %q instead", field.Name)
	}

	if field.Type.Local != "travellerInfo" {
		t.Errorf("Field type should be 'travellerInfo', %q instead", field.Type.Local)
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

	if field.Type.Local != "string" {
		t.Errorf("Field type should be 'string', %q instead", field.Type.Local)
	}
}

func TestAttributeGroup(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	name := "Test"
	if cType.Local != name {
		t.Errorf("TypeName elemName should be %q, got %q instead", name, cType.Local)
	}

	if len(cType.Fields) != 6 {
		t.Fatalf("Fields amount should be 6, %d instead", len(cType.Fields))
	}

	field := cType.Fields[1]
	if field.Name != "Code" {
		t.Errorf("Field elemName should be 'Code', %q instead", field.Name)
	}

	if !field.IsAttr {
		t.Errorf("Owner should be attribute")
	}
}

func TestParseElementRef(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	if len(typesList[0].Fields) != 2 {
		t.Fatalf("Fields amount should be 2, %d got", len(typesList[0].Fields))
	}

	f := typesList[0].Fields[1]
	if f.Name != "TPA_Extensions" {
		t.Errorf("Field name shoud be 'TPA_Extensions', %q got", f.Name)
	}
}

func TestInclude(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	assert.Len(t, typesList, 2)

	tp := typesList[1]

	if tp.Local != "AirShoppingRQ" {
        t.Errorf(`Type name shoud be "AirShoppingRQ", %q`, tp.Local)
    }

    if len(tp.Fields) != 2 {
        t.Fatalf("Wrong fields amount. 2 expected, %d got", len(tp.Fields))
    }

	f := tp.Fields[1]
	if f.Name != "PointOfSale" {
        t.Errorf("Field name shoud be 'PointOfSale', %q got", f.Name)
    }
}

func TestSimpleTypeAttribute(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	if len(typesList[0].Fields) != 2 {
        t.Fatalf("Wrong fields amount. 2 expected, %d got", len(typesList[0].Fields))
    }

    f := typesList[0].Fields[1]
    if f.Name != "PieceAllowanceCombination" {
        t.Errorf(`Wrong field name. "PieceAllowanceCombination" expected, %q got`, f.Name)
    }
}

func TestUnion(t *testing.T) {
	typesList := parseTypesFrom(t.Name())
	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	if typesList[0].Local != "Test" {
		t.Errorf("Type name shoud be 'Test', %q given", typesList[0].Local)
	}

	if len(typesList[0].Fields) != 2 {
		t.Fatalf("Fields amount should be 2, %d given", len(typesList[0].Fields))
	}

    if typesList[0].Fields[1].Type.Local != "string" {
        t.Errorf("Type base type shoud be 'string', %q given", typesList[0].Fields[1].Type.Local)
    }
}

func TestSimpleContent(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	if cType.Local != "Test" {
		t.Fatalf(`TypeName name should be "Test", got %q instead`, cType.Local)
	}

	if len(cType.Fields) != 4 {
		t.Fatalf("Test should has 4 field, %d instead", len(cType.Fields))
	}

	field := cType.Fields[3]

	name := "XMLValue"
	if field.Name != name {
		t.Errorf("Field name should be %q, got %q instead", name, field.Name)
	}

	fType := "string"
	if field.Type.Local != fType {
		t.Errorf("Field type should be %q, got %q instead", fType, field.Type.Local)
	}
}

func TestComplexContent(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]

	name := "Test"
	if cType.Local != name {
		t.Errorf("TypeName name should be %q, got %q instead", name, cType.Name)
	}

	if len(cType.Fields) != 4 {
		t.Fatalf("Fields amount should be 4, got %d instead", len(cType.Fields))
	}

	field := cType.Fields[1]
	fieldName := "Format"
	if field.Name != fieldName {
		t.Fatalf("TypeName name should be %q, got %q instead", fieldName, field.Name)
	}
}

func TestComplexTypeSimpleContent(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 2 expected, %d got", len(typesList))
	}

	if len(typesList[0].Fields) != 2 {
        t.Fatalf("Wrong fields amount. 2 expected, %d got", len(typesList))
    }

    f := typesList[0].Fields[1]
    if f.Name != "City" {
        t.Errorf(`Wrong field name. Expected "City", got %q`, f.Name)
    }

    if len(typesList[1].Fields) != 2 {
        t.Fatalf("Wrong fields amount. 1 expected, %d got", len(typesList))
    }
}

func TestComplexTypeWithSimpleContent(t *testing.T) {
    typesList := parseTypesFrom(t.Name())

    if len(typesList) != 1 {
        t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
    }
    ct := typesList[0]

    if len(ct.Fields) != 4 {
        t.Fatalf("Wrong type fields amount. 4 expected, %d got", len(typesList))
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

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

    if len(typesList[0].Fields) != 2 {
        t.Fatalf("Wrong type fields amount. 2 expected, %d got", len(typesList[0].Fields))
    }

    if len(typesList[1].Fields) != 1 {
        t.Fatalf("Wrong type fields amount. 1 expected, %d got", len(typesList[0].Fields))
    }

    f := typesList[0].Fields[1]
    if f.Name != "PointOfSale" {
        t.Errorf(`Wrong field name. "PointOfSale" expected, %q got`, f.Name)
    }
}

func TestComplexChoice(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0]
	if len(cType.Fields) != 4 {
		t.Fatalf("Wrong type fields amount. 3 expected, %d got", len(cType.Fields))
	}

	if cType.Fields[3].Name != "Errors" {
        t.Errorf(`Field name should be "Errors", %q got`, cType.Fields[3].Name)
    }
}

func TestSimpleTypesFolding(t *testing.T) {
    typesList := parseTypesFrom(t.Name())

    if len(typesList) != 1 {
        t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
    }

    ct := typesList[0]
    if ct.Local != "CountryCode" {
        t.Errorf("Type name shoud be 'OwnerSimpleType', %q given", ct.Name)
    }

    fields := ct.Fields

    if len(fields) != 3 {
        t.Errorf("Wrong type fields amount. 2 expected, %d got", len(fields))
    }

    if fields[2].Name != "XMLValue" {
        t.Errorf(`Last field should be "XMLValue", %q given`, fields[2].Name)
    }

    if fields[2].Type.Local != "string" {
        t.Errorf("Last field type should be string")
    }
}

func TestTwoLevelSimpleContent(t *testing.T) {
    typesList := parseTypesFrom(t.Name())

    if len(typesList) != 1 {
        t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
    }

    ct := typesList[0]
    if len(ct.Fields) != 4 {
        t.Fatalf("Wrong fields amount. 4 expected, %d got", len(ct.Fields))
    }

    if ct.Fields[1].Name != "Language" {
        t.Errorf(`Last field should be "Language", %q given`, ct.Fields[1].Name)
    }
}

func TestFieldsDuplication(t *testing.T) {
    typesList := parseTypesFrom(t.Name())

    if len(typesList) != 1 {
        t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
    }

    ct := typesList[0]
    if len(ct.Fields) != 3 {
        t.Fatalf("Wrong fields amount. 3 expected, %d got", len(ct.Fields))
    }

    if ct.Fields[1].Name != "Name" {
        t.Errorf(`Last field should be "Name", %q given`, ct.Fields[1].Name)
    }
}

func TestSequenceChoice(t *testing.T) {
    typesList := parseTypesFrom(t.Name())

    if len(typesList) != 1 {
        t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
    }

    ct := typesList[0]
    if len(ct.Fields) != 5 {
        t.Fatalf("Wrong fields amount. 5 expected, %d got", len(ct.Fields))
    }

    if ct.Fields[4].Name != "FlightSegmentReference2" {
        t.Errorf(`Last field should be "FlightSegmentReference2", %q given`, ct.Fields[1].Name)
    }
}

func TestRenameDuplicatedTypes(t *testing.T) {
    // typesList := parseTypesFrom(t.Name())
	//
    // assert.Lenf(t, typesList, 2, "Wrong types amount")
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