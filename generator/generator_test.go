package generator

import (
	"testing"
	"gosoapgen/xsd"
)

func TestGetNoTypes(t *testing.T) {
	res := GenerateTypes(xsd.Schema{})

	if len(res) != 0 {
		t.Errorf("Should be no types")
	}
}

func TestSingleEmptyStructure(t *testing.T) {
	s := xsd.Schema{}
	s.Element = append(s.Element, xsd.Element{Name: "Session"})

	res := GenerateTypes(s)
	if len(res) != 1 {
		t.Errorf("Should be 1 type, %d getting", len(res))
	}
}