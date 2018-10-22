package xsd_model

import (
    "testing"
    "os"
)

func TestNode(t *testing.T) {
    f, _ := os.Open("./testdata/simple.xml")
    n := Load(f)

    if n == nil {
        t.Fatalf("No tree")
    }

    if n.name != "note" {
        t.Errorf("Element name should be note")
    }

    if len(n.Children()) != 4 {
        t.Errorf("children amount should be 4< %d got", len(n.Children()))
    }
}

func TestSchema(t *testing.T) {
    f, _ := os.Open("./testdata/simple.xsd")
    n := Load(f)

    if n.TargetNamespace != "http://xml.amadeus.com/FLIREQ_07_1_1A" {
        t.Errorf("Error target namespace. Got %q", n.TargetNamespace)
    }

    ns := n.ResolveSpace("altova")
    if ns != "http://www.altova.com/xml-schema-extensions" {
        t.Errorf("Error namespace resolving. Got %q", ns)
    }
}