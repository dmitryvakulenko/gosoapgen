package flat_parser

import (
	"testing"
	"io/ioutil"
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

	tp, ok := typesList[0].(*SimpleType)
	if !ok {
		t.Fatalf("Type should be *SimpleType")
	}

	name := "AlphaString_Length1To2"
	if name != tp.Name {
		t.Fatalf("Type name should be %q, got %q instead", name, tp.Name)
	}

}

func TestParseElements(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType, ok := typesList[0].(*ComplexType)
	if !ok {
		t.Fatalf("Type should be complex type")
	}

	typeName := "Session"
	if cType.GoName != typeName {
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

	if field.Type.GetName() != "string" {
		t.Errorf("Field type should be 'string' %q instead", field.Type.GetName())
	}

	field = cType.Fields[3]
	if field.Name != "TransactionStatusCode" {
		t.Errorf("Field name should be 'TransactionStatusCode' %s instead", field.Name)
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

	cType, ok := typesList[1].(*ComplexType)
	if !ok {
		t.Fatalf("Type should be complex type")
	}

	typeName := "AMA_SecurityHostedUser"
	if cType.GoName != typeName {
		t.Errorf("TypeName name should be %q, got %q", typeName, cType.GoName)
	}

	ns := "http://xml.amadeus.com/2010/06/Security_v1"
	if cType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, cType.Namespace)
	}

	if len(cType.Fields) != 4 {
		t.Fatalf("TypeName should Has 4 fields, %d getting", len(cType.Fields))
	}

}

func TestComplexTypeWithAttributes(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType, ok := typesList[0].(*ComplexType)
	if !ok {
		t.Fatalf("Type should be complex type")
	}

	typeName := "Session"
	if cType.GoName != typeName {
		t.Errorf("TypeName name should be %q, got %q", typeName, cType.GoName)
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
		t.Errorf("Field name should be 'TransactionStatusCode' %s instead", field.Name)
	}

	if field.Type.GetName() != "string" {
		t.Errorf("Field type should be 'string' %s instead", field.Type.GetName())
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

	firstType, ok := typesList[1].(*ComplexType)
	if !ok {
		t.Fatalf("Type should be complex type")
	}


	secType, ok := typesList[0].(*ComplexType)
	if !ok {
		t.Fatalf("Type should be complex type")
	}

	typeName := "PNR_AddMultiElements"
	if firstType.GoName != typeName {
		t.Errorf("TypeName name should be %q, got %q", typeName, firstType.GoName)
	}

	ns := "http://xml.amadeus.com/PNRADD_10_1_1A"
	if firstType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, firstType.Namespace)
	}

	typeName = "TravellerInfo"
	if secType.GoName != typeName {
		t.Errorf("TypeName name should be %q, got %q", typeName, secType.GoName)
	}

	typeName = "TravellerInfo"
	if secType.Name != typeName {
		t.Errorf("TypeName name should be %q, got %q", typeName, secType.Name)
	}

	if secType.Namespace != ns {
		t.Errorf("TypeName namespace should be %q, got %q", ns, secType.Namespace)
	}

	if len(firstType.Fields) != 1 {
		t.Fatalf("Should be 1 fields, %d getting", len(firstType.Fields))
	}

	field := firstType.Fields[0]
	if field.Name != "travellerInfo" {
		t.Errorf("Field name should be 'travellerInfo', %q instead", field.Name)
	}

	if field.Type.GetName() != "TravellerInfo" {
		t.Errorf("Field type should be 'travellerInfo', %q instead", field.Type.GetName())
	}

	if len(secType.Fields) != 1 {
		t.Fatalf("Second type fields amount should be 1, got %d instead", len(secType.Fields))
	}

	if secType.Fields[0].Name != "elementManagementPassenger" {
		t.Errorf("Second type name shoud be 'ElementManagementPassenger', %q instead", secType.Fields[0].Name)
	}
}

func TestAttributeGroup(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType, ok := typesList[0].(*ComplexType)
	if !ok {
		t.Fatalf("Type should be complex type")
	}

	name := "CodeType"
	if cType.GoName != "CodeType" {
		t.Fatalf("TypeName name should be %q, got %q instead", name, cType.GoName)
	}

	if len(cType.Fields) != 5 {
		t.Fatalf("Fields amount should be 5, %d instead", len(cType.Fields))
	}

	field := cType.Fields[1]
	if field.Name != "Owner" {
		t.Fatalf("Field name should be 'Owner', %q instead", field.Name)
	}

	if !field.IsAttr {
		t.Fatalf("Owner should be attribute")
	}
}

func TestSimpleContent(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 4 expected, %d got", len(typesList))
	}

	cType, ok := typesList[1].(*ComplexType)
	if !ok {
		t.Fatalf("Type should be complex type")
	}

	name := "CompanyNameType"
	if cType.GoName != name {
		t.Fatalf("TypeName name should be %q, got %q instead", name, cType.GoName)
	}

	if len(cType.Fields) != 3 {
		t.Errorf("CompanyNameType should Has 3 field, %d instead", len(cType.Fields))
	}

	field := cType.Fields[0]

	name = "Value"
	if field.Name != name {
		t.Fatalf("Field name should be %q, got %q instead", name, field.Name)
	}

	fType := "StringLength0to128"
	if field.Type.GetName() != fType {
		t.Fatalf("Field type should be %q, got %q instead", fType, field.Type.GetName())
	}
}

func TestComplexContent(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 4 {
		t.Fatalf("Wrong types amount. 4 expected, %d got", len(typesList))
	}

	cType, ok := typesList[3].(*ComplexType)
	if !ok {
		t.Fatalf("Type should be complex type")
	}

	typeName := "AddressWithModeType"
	if cType.GoName != typeName {
		t.Fatalf("TypeName name should be %q, got %q instead", typeName, cType.GoName)
	}

	if len(cType.Fields) != 3 {
		t.Fatalf("Fields amount should be 3, got %q instead", len(cType.Fields))
	}

	field := cType.Fields[2]
	fieldName := "Mode"
	if field.Name != fieldName {
		t.Fatalf("TypeName name should be %q, got %q instead", fieldName, field.Name)
	}
}

//func TestNoDuplicateTypes(t *testing.T) {
//	Decoder := NewDecoder()
//	s := loadSchemaFrom("complexType.xsd")
//	Decoder.Decode(&s, s.TargetNamespace)
//	s = loadSchemaFrom("complexType.xsd")
//	Decoder.Decode(&s, s.TargetNamespace)
//
//	typesList := Decoder.GenerateTypes()
//	if len(typesList) != 1 {
//		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
//	}
//
//	namespaces := Decoder.GetNamespaces()
//	if len(namespaces) != 1 {
//		t.Fatalf("Wrong namespaces amount. 1 expected, %d got", len(namespaces))
//	}
//}

//func TestWrongTypesOrder(t *testing.T) {
//	ns := "namespace"
//	typesList := parseTypesFrom("complexContentWrongOrder.xsd", ns)
//
//	if len(typesList) != 2 {
//		t.Fatalf("Wrong types amount. 4 expected, %d got", len(typesList))
//	}
//}
//
//
func TestParseElementRef(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 2 expected, %d got", len(typesList))
	}
}

func TestParseComplexElement(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 6 {
		t.Fatalf("Wrong types amount. 6 expected, %d got", len(typesList))
	}

	cType, ok := typesList[3].(*ComplexType)
	if !ok {
		t.Fatalf("Type should be complex type")
	}

	field := cType.Fields[0]
	if field.Type.GetName() != "TravelSector" {
		t.Errorf("Field type should be %q, %q got", "TravelSector", field.Type.GetName())
	}
}


func TestArrayAlias(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 3 {
		t.Fatalf("Wrong types amount. 3 expected, %d got", len(typesList))
	}

	firstType := typesList[0]
	if firstType.GetName() != "AddressMetadatas" {
		t.Errorf("First type name should be %q, %q got", "TravelSector", firstType.GetName())
	}

	firstField := (firstType).(*ComplexType).Fields[0]
	if firstField.Name != "AddressMetadata" {
		t.Errorf("Field name should be %q, %q got", "AddressMetadata", firstType.GetName())
	}
}

func TestChoiceParsing(t *testing.T) {
	typesList := parseTypesFrom(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 4 expected, %d got", len(typesList))
	}

	fields := typesList[0].(*ComplexType).Fields
	if len(fields) != 2 {
		t.Fatalf("Wrong fields amount. 2 expected, %d got", len(fields))
	}
}

//func TestElementsDuplication(t *testing.T) {
//	typesList := parseTypesFrom(t.Name())
//
//	if len(typesList) != 4 {
//		t.Fatalf("Wrong types amount. 4 expected, %d got", len(typesList))
//	}
//}

func parseTypesFrom(name string) []NamedType {
	decoder := NewDecoder(&Ld{})
	decoder.Decode(name + ".xsd")

	return decoder.GetTypes()
}


type Ld struct {}

func (l *Ld) Load(path string) ([]byte, error) {
	data, _ := ioutil.ReadFile("./schema_test/" + path)
	return data, nil
}

func (l *Ld) IsAlreadyLoadedError(e error) bool {
	return false
}