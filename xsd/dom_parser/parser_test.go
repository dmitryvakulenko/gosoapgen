package dom_parser

import "testing"

func TestEmptySchema(t *testing.T) {
	typesList := loadXsd(t.Name())

	if len(typesList) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestSimpleElements(t *testing.T) {
	typesList := loadXsd(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong number of types. 1 expected, but got %d", len(typesList))
	}

	tp := typesList[0]
    if len(tp.Fields) != 2 {
        t.Fatalf("Wrong number of type fields. 2 expected, but got %d", len(tp.Fields))
    }

    if tp.Fields[0].Name != "XMLName" {
        t.Errorf(`First field name should be "XMLName", %q got`, tp.Fields[0].Name)
    }

    if tp.Fields[1].Name != "XMLValue" {
        t.Errorf(`First field name should be "XMLValue", %q got`, tp.Fields[1].Name)
    }
}

//func TestSimpleTypes(t *testing.T) {
//	typesList := loadXsd(t.Name())
//
//	if len(typesList) != 1 {
//		t.Fatalf("Wrong number of types. 1 expected, but got %d", len(typesList))
//	}
//
//	tp := typesList[0]
//
//	if !tp.IsSimpleContent {
//		t.Fatalf("Type should be complex type")
//	}
//
//	name := "AlphaString_Length1To2"
//	if name != tp.Name {
//		t.Errorf("Field elemName should be %q, got %q instead", name, tp.Name)
//	}
//
//	if tp.BaseTypeName.Name != "string" {
//		t.Errorf("Field should be string, got %q instead", tp.Name)
//	}
//}


func loadXsd(testName string) []*Type {
	p := NewParser()
	p.LoadFile("./test_data/" + testName + ".xsd")
	return p.GetTypes()
}
