package translator

import "testing"

func TestGetTypeFromEmptyCollection(t *testing.T) {
	collector := makeTypesCollection()

	curType, ok := collector.find("aa", "bbbb")
	if ok {
		t.Errorf("Should be error")
	}

	if curType != nil {
		t.Errorf("Type should be nil")
	}
}