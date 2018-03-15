package translator

import (
	"encoding/xml"
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"os"
	"testing"
)

func TestGetNoTypes(t *testing.T) {
	typesList := parseTypesFrom("empty.xsd", "")
	if len(typesList) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestGenerateSimpleTypes(t *testing.T) {
	typesList := parseTypesFrom("simpleType.xsd", "")

	if len(typesList) != 0 {
		t.Fatalf("Wrong number of types. 0 expected, but got %d", len(typesList))
	}
}

func TestParseElementTypes(t *testing.T) {
	typesList := parseTypesFrom("element.xsd", "")

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0].(*ComplexType)

	typeName := "Session"
	if cType.Name != typeName {
		t.Errorf("Type name should be %q, got %q", typeName, cType.Name)
	}

	ns := "http://xml.amadeus.com/2010/06/Session_v3"
	if cType.Namespace != ns {
		t.Errorf("Type namespace should be %q, got %q", ns, cType.Namespace)
	}

	if len(cType.Fields) != 4 {
		t.Fatalf("Should be 4 fields, %d getting", len(cType.Fields))
	}

	field := cType.Fields[1]
	if field.Name != "sequenceNumber" {
		t.Errorf("Field name should be 'sequenceNumber', %q instead", field.Name)
	}

	if field.Type != "string" {
		t.Errorf("Field type should be 'string' %s instead", field.Type)
	}

	if field.XmlExpr != "sequenceNumber" {
		t.Errorf("Field xml expression should be 'sequenceNumber', %q instead", field.XmlExpr)
	}

	if field.Namespace != ns {
		t.Errorf("Type should be %q, %q getting", ns, cType.Namespace)
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
	typesList := parseTypesFrom("complexType.xsd", "")

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0].(*ComplexType)

	typeName := "AMA_SecurityHostedUser"
	if cType.Name != typeName {
		t.Errorf("Type name should be %q, got %q", typeName, cType.Name)
	}

	ns := "http://xml.amadeus.com/2010/06/Security_v1"
	if cType.Namespace != ns {
		t.Errorf("Type namespace should be %q, got %q", ns, cType.Namespace)
	}

	if len(cType.Fields) != 4 {
		t.Fatalf("Type should has 4 fields, %d getting", len(cType.Fields))
	}

}

func TestComplexTypeWithAttributes(t *testing.T) {
	typesList := parseTypesFrom("attribute.xsd", "")

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0].(*ComplexType)

	typeName := "Session"
	if cType.Name != typeName {
		t.Errorf("Type name should be %q, got %q", typeName, cType.Name)
	}

	ns := "http://xml.amadeus.com/2010/06/Session_v3"
	if cType.Namespace != ns {
		t.Errorf("Type namespace should be %q, got %q", ns, cType.Namespace)
	}

	if len(cType.Fields) != 1 {
		t.Fatalf("Should be 1 fields, %d getting", len(cType.Fields))
	}

	field := cType.Fields[0]
	if field.Name != "TransactionStatusCode" {
		t.Errorf("Field name should be 'TransactionStatusCode' %s instead", field.Name)
	}

	if field.Type != "string" {
		t.Errorf("Field type should be 'string' %s instead", field.Type)
	}

	if field.XmlExpr != "TransactionStatusCode,attr" {
		t.Errorf("Field xml name should be 'sequenceNumber' %s instead", field.XmlExpr)
	}
}

func TestInnerComplexTypes(t *testing.T) {
	typesList := parseTypesFrom("innerComplexType.xsd", "")

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 2 expected, %d got", len(typesList))
	}

	firstType := typesList[0].(*ComplexType)
	secType := typesList[1].(*ComplexType)

	typeName := "PNR_AddMultiElements"
	if firstType.Name != typeName {
		t.Errorf("Type name should be %q, got %q", typeName, firstType.Name)
	}

	ns := "http://xml.amadeus.com/PNRADD_10_1_1A"
	if firstType.Namespace != ns {
		t.Errorf("Type namespace should be %q, got %q", ns, firstType.Namespace)
	}

	typeName = "travellerInfo"
	if secType.Name != typeName {
		t.Errorf("Type name should be %q, got %q", typeName, secType.Name)
	}

	if secType.Namespace != ns {
		t.Errorf("Type namespace should be %q, got %q", ns, secType.Namespace)
	}

	if len(firstType.Fields) != 1 {
		t.Fatalf("Should be 1 fields, %d getting", len(firstType.Fields))
	}

	field := firstType.Fields[0]
	if field.Name != "travellerInfo" {
		t.Errorf("Field name should be 'travellerInfo', %q instead", field.Name)
	}

	if field.Type != "travellerInfo" {
		t.Errorf("Field type should be 'travellerInfo', %q instead", field.Type)
	}

	if field.XmlExpr != "travellerInfo" {
		t.Errorf("Field xml name should be 'travellerInfo' %s instead", field.XmlExpr)
	}

	if len(secType.Fields) != 1 {
		t.Fatalf("Second type fields amount should be 1, got %d instead", len(secType.Fields))
	}

	if secType.Fields[0].Name != "elementManagementPassenger" {
		t.Errorf("Second type name shoud be 'ElementManagementPassenger', %q instead", secType.Fields[0].Name)
	}

	if secType.Fields[0].XmlExpr != "elementManagementPassenger" {
		t.Errorf("Second type xml expression shoud be 'elementManagementPassenger', %q instead", secType.Fields[0].XmlExpr)
	}
}

func TestAttributeGroup(t *testing.T) {
	typesList := parseTypesFrom("attributeGroup.xsd", "")

	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	cType := typesList[0].(*ComplexType)

	name := "CodeType"
	if cType.Name != "CodeType" {
		t.Fatalf("Type name should be %q, got %q instead", name, cType.Name)
	}

	if len(cType.Fields) != 5 {
		t.Fatalf("Fields amount should be 5, %d instead", len(cType.Fields))
	}

	field := cType.Fields[1]
	if field.Name != "Owner" {
		t.Fatalf("Field name should be 'Owner', %q instead", field.Name)
	}

	if field.XmlExpr != "Owner,attr" {
		t.Fatalf("Field xml expression should be 'Owner,attr', %q instead", field.XmlExpr)
	}
}

func TestSimpleContent(t *testing.T) {
	ns := "namespace"
	typesList := parseTypesFrom("simpleContent.xsd", ns)

	if len(typesList) != 2 {
		t.Fatalf("Wrong types amount. 2 expected, %d got", len(typesList))
	}

	cType := typesList[1].(*ComplexType)
	name := "CompanyNameType"
	if cType.Name != name {
		t.Fatalf("Type name should be %q, got %q instead", name, cType.Name)
	}

	if len(cType.Fields) != 3 {
		t.Errorf("CompanyNameType should has 3 field, %d instead", len(cType.Fields))
	}

	field := cType.Fields[0]

	name = "Value"
	if field.Name != name {
		t.Fatalf("Field name should be %q, got %q instead", name, field.Name)
	}

	fType := "StringLength0to128"
	if field.Type != fType {
		t.Fatalf("Field type should be %q, got %q instead", fType, field.Type)
	}
}

func TestComplexContent(t *testing.T) {
	ns := "namespace"
	typesList := parseTypesFrom("complexContent.xsd", ns)

	if len(typesList) != 4 {
		t.Fatalf("Wrong types amount. 4 expected, %d got", len(typesList))
	}

	cType := typesList[3].(*ComplexType)
	typeName := "AddressWithModeType"
	if cType.Name != typeName {
		t.Fatalf("Type name should be %q, got %q instead", typeName, cType.Name)
	}

	if len(cType.Fields) != 3 {
		t.Fatalf("Fields amount should be 3, got %q instead", len(cType.Fields))
	}

	field := cType.Fields[2]
	fieldName := "Mode"
	if field.Name != fieldName {
		t.Fatalf("Type name should be %q, got %q instead", fieldName, field.Name)
	}
}

func TestNoDuplicateTypes(t *testing.T) {
	decoder := newDecoder()
	s := loadSchemaFrom("complexType.xsd")
	decoder.decode(&s, s.TargetNamespace)
	s = loadSchemaFrom("complexType.xsd")
	decoder.decode(&s, s.TargetNamespace)

	typesList := decoder.GetTypes()
	if len(typesList) != 1 {
		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
	}

	namespaces := decoder.GetNamespaces()
	if len(namespaces) != 1 {
		t.Fatalf("Wrong namespaces amount. 1 expected, %d got", len(namespaces))
	}
}

func TestWrongTypesOrder(t *testing.T) {
	ns := "namespace"
	typesList := parseTypesFrom("complexContentWrongOrder.xsd", ns)

	if len(typesList) != 4 {
		t.Fatalf("Wrong types amount. 4 expected, %d got", len(typesList))
	}
}


func parseTypesFrom(name, namespace string) []interface{} {
	s := loadSchemaFrom(name)
	res := newDecoder()
	if namespace != "" {
		res.decode(&s, namespace)
	} else {
		res.decode(&s, s.TargetNamespace)
	}


	return res.GetTypes()
}


func loadSchemaFrom(name string) xsd.Schema {
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

	return s
}