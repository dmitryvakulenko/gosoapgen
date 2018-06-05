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



func TestSimpleTypes(t *testing.T) {
	typesList := loadXsd(t.Name())

	if len(typesList) != 1 {
		t.Fatalf("Wrong number of types. 1 expected, but got %d", len(typesList))
	}

	tp := typesList[0]
	name := "Test"
	if name != tp.Local {
		t.Errorf("Type name should be %q, got %q instead", name, tp.Name)
	}

	ns := "http://xml.amadeus.com/PNRADD_10_1_1A"
    if ns != tp.Space {
        t.Errorf("Type space should be %q, got %q instead", ns, tp.Space)
    }

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


func loadXsd(testName string) []*Type {
	p := NewParser()
	p.LoadFile("./test_data/" + testName + ".xsd")
	return p.GetTypes()
}
