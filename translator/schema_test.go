package translator

import (
	"testing"
	"github.com/dmitryvakulenko/gosoapgen/xsd"
	"os"
	"encoding/xml"
)

func TestGetNoTypes(t *testing.T) {
	s := loadXsd("empty.xsd")
	res := Parse(&s)

	if len(res.cType) != 0 {
		t.Errorf("Should be no types")
	}

	if len(res.sType) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestGenerateSimpleTypes(t *testing.T) {
	s := loadXsd("simpleType.xsd")

	res := Parse(&s)
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
}

func TestParseElementTypes(t *testing.T) {
	s := loadXsd("element.xsd")
	res := Parse(&s)

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
	res := Parse(&s)
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


//func TestComplexTypeWithAttributes(t *testing.T) {
//	elem := xsd.Element{Name: "Session"}
//	elem.ComplexType = &xsd.ComplexType{}
//	elem.ComplexType.Attribute = append(elem.ComplexType.Attribute, &xsd.Attribute{Name: "TransactionStatusCode", Type: "xs:NMTOKEN"})
//
//	s := xsd.Schema{}
//	s.Element = append(s.Element, &elem)
//
//	schemas := []*xsd.Schema{&s}
//	res := Parse(&s)
//
//	if len(res[0].Fields) != 1 {
//		t.Fatalf("Should be 1 fields, %d getting", len(res[0].Fields))
//	}
//
//	field := res[0].Fields[0]
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
//
//func TestInnerComplexTypes(t *testing.T) {
//	innerElem := xsd.Element{}
//	innerElem.Name = "innerElement"
//	innerElem.ComplexType = &xsd.ComplexType{}
//	innerElem.ComplexType.Sequence = &xsd.Sequence{}
//	innerElem.ComplexType.Sequence.Element = append(innerElem.ComplexType.Sequence.Element, &xsd.Element{Name: "innerField", Type: "xs:string"})
//
//	elem := xsd.Element{Name: "Session"}
//	elem.ComplexType = &xsd.ComplexType{}
//	elem.ComplexType.Sequence = &xsd.Sequence{}
//	elem.ComplexType.Sequence.Element = append(elem.ComplexType.Sequence.Element, &innerElem)
//
//	s := xsd.Schema{}
//	s.Element = append(s.Element, &elem)
//
//	schemas := []*xsd.Schema{&s}
//	res := Parse(&s)
//
//	if len(res) != 2 {
//		t.Fatalf("Types amount should be 2, %d instead", len(res))
//	}
//
//	if len(res[0].Fields) != 1 {
//		t.Fatalf("Should be 1 fields, %d getting", len(res[0].Fields))
//	}
//
//	field := res[0].Fields[0]
//	if field.Name != "InnerElement" {
//		t.Errorf("Field name should be 'InnerElement' %s instead", field.Name)
//	}
//
//	if field.Type != "InnerElement" {
//		t.Errorf("Field type should be 'InnerElement' %s instead", field.Type)
//	}
//
//	if field.XmlExpr != "innerElement" {
//		t.Errorf("Field xml name should be 'innerElement' %s instead", field.XmlExpr)
//	}
//
//	if res[1].Name != "InnerElement" {
//		t.Errorf("Second type name shoud be 'InnerElement' %s instead", res[1].Name)
//	}
//
//	if len(res[1].Fields) != 1 {
//		t.Fatalf("Second type fields amount should be 2, got %d instead", len(res[1].Fields))
//	}
//
//	if res[1].Fields[0].Name != "InnerField" {
//		t.Errorf("Second type name shoud be 'InnerField' %s instead", res[1].Fields[0].Name)
//	}
//
//	if res[1].Fields[0].XmlExpr != "innerField" {
//		t.Errorf("Second type name shoud be 'innerField' %s instead", res[1].Fields[0].XmlExpr)
//	}
//}
//
//func TestAttributeGroup(t *testing.T) {
//	group := xsd.AttributeGroup{}
//	group.Name = "attributeGroup"
//	group.Attribute = append(group.Attribute, &xsd.Attribute{Name: "innerAttribute", Type: "xs:string"})
//
//	inGr := xsd.AttributeGroup{}
//	inGr.Ref = group.Name
//
//	elem := xsd.ComplexType{Name: "Session"}
//	elem.AttributeGroup = append(elem.AttributeGroup, &inGr)
//
//	s := xsd.Schema{}
//	s.ComplexType = append(s.ComplexType, &elem)
//	s.AttributeGroup = append(s.AttributeGroup, &group)
//
//	schemas := []*xsd.Schema{&s}
//	res := Parse(&s)
//
//	if len(res) != 2 {
//		t.Fatalf("Types amount should be 2, %d instead", len(res))
//	}
//
//	if res[0].Name != "attributeGroup" {
//		t.Fatalf("No attributeGroup type")
//	}
//
//	if len(res[1].Embed) != 1 {
//		t.Fatalf("Embed types amount should be 1")
//	}
//
//	if res[1].Embed[0] != group.Name {
//		t.Fatalf("Embed types amount should be " + group.Name)
//	}
//}

func loadXsd(name string) xsd.Schema {
	reader, err := os.Open("./translator/schema_test/" + name)
	defer reader.Close()

	if err != nil {
		panic(err)
	}

	s := xsd.Schema{};
	err = xml.NewDecoder(reader).Decode(&s)
	if err != nil {
		panic(err)
	}

	return s
}