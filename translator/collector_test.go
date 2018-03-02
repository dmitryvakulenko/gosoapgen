package translator

import "testing"

func TestGetTypeFromEmptyCollection(t *testing.T) {
	collector := new(xsdTypes)

	curType, ok := collector.findType("aa", "bbbb")
	if ok {
		t.Errorf("Should be error")
	}

	if curType != nil {
		t.Errorf("Type should be nil")
	}
}