package xml

import (
    "testing"
    "os"
)

func TestSimple(t *testing.T) {
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