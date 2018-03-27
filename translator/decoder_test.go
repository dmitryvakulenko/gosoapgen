package translator

import (
	"encoding/xml"
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"os"
	"testing"
)

func TestEmptySchema(t *testing.T) {
	typesList := parseTypesFrom(t.Name(), "")
	if len(typesList) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestSimpleTypes(t *testing.T) {
	typesList := parseTypesFrom(t.Name(), "")

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
	typesList := parseTypesFrom(t.Name(), "")

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
	typesList := parseTypesFrom(t.Name(), "")

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
		t.Fatalf("TypeName should has 4 fields, %d getting", len(cType.Fields))
	}

}

func TestComplexTypeWithAttributes(t *testing.T) {
	typesList := parseTypesFrom(t.Name(), "")

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
	typesList := parseTypesFrom(t.Name(), "")

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
	typesList := parseTypesFrom(t.Name(), "")

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

//func TestSimpleContent(t *testing.T) {
//	ns := "namespace"
//	typesList := parseTypesFrom("simpleContent.xsd", ns)
//
//	if len(typesList) != 1 {
//		t.Fatalf("Wrong types amount. 2 expected, %d got", len(typesList))
//	}
//
//	cType := typesList[0]
//	name := "CompanyNameType"
//	if cType.GoName != name {
//		t.Fatalf("TypeName name should be %q, got %q instead", name, cType.GoName)
//	}
//
//	if len(cType.Fields) != 3 {
//		t.Errorf("CompanyNameType should has 3 field, %d instead", len(cType.Fields))
//	}
//
//	field := cType.Fields[0]
//
//	name = "Value"
//	if field.GoName != name {
//		t.Fatalf("Field name should be %q, got %q instead", name, field.GoName)
//	}
//
//	fType := "string"
//	if field.Type != fType {
//		t.Fatalf("Field type should be %q, got %q instead", fType, field.Type)
//	}
//}
//
//func TestComplexContent(t *testing.T) {
//	ns := "namespace"
//	typesList := parseTypesFrom("complexContent.xsd", ns)
//
//	if len(typesList) != 2 {
//		t.Fatalf("Wrong types amount. 4 expected, %d got", len(typesList))
//	}
//
//	cType := typesList[1]
//	typeName := "AddressWithModeType"
//	if cType.GoName != typeName {
//		t.Fatalf("TypeName name should be %q, got %q instead", typeName, cType.GoName)
//	}
//
//	if len(cType.Fields) != 3 {
//		t.Fatalf("Fields amount should be 3, got %q instead", len(cType.Fields))
//	}
//
//	field := cType.Fields[2]
//	fieldName := "Mode"
//	if field.GoName != fieldName {
//		t.Fatalf("TypeName name should be %q, got %q instead", fieldName, field.GoName)
//	}
//}
//
//func TestNoDuplicateTypes(t *testing.T) {
//	decoder := newDecoder()
//	s := loadSchemaFrom("complexType.xsd")
//	decoder.decode(&s, s.TargetNamespace)
//	s = loadSchemaFrom("complexType.xsd")
//	decoder.decode(&s, s.TargetNamespace)
//
//	typesList := decoder.GetTypes()
//	if len(typesList) != 1 {
//		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
//	}
//
//	namespaces := decoder.GetNamespaces()
//	if len(namespaces) != 1 {
//		t.Fatalf("Wrong namespaces amount. 1 expected, %d got", len(namespaces))
//	}
//}
//
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
//func TestParseElementRef(t *testing.T) {
//	ns := "namespace"
//	typesList := parseTypesFrom("elementRef.xsd", ns)
//
//	if len(typesList) != 1 {
//		t.Fatalf("Wrong types amount. 1 expected, %d got", len(typesList))
//	}
//}
//
//func TestParseComplexElement(t *testing.T) {
//	ns := "namespace"
//	typesList := parseTypesFrom("complexElement.xsd", ns)
//
//	if len(typesList) != 3 {
//		t.Fatalf("Wrong types amount. 3 expected, %d got", len(typesList))
//	}
//
//	field := typesList[2].Fields[0]
//	if field.Type != "string" {
//		t.Errorf("Field type should be string, %q got", field.Type)
//	}
//}


func parseTypesFrom(name, namespace string) []NamedType {
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
	reader, err := os.Open("./translator/schema_test/" + name + ".xsd")
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