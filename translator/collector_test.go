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

func TestStoringType(t *testing.T) {
	collector := makeTypesCollection()

	ns := "namespace"
	typeName := "ComplexType"
	addedType := &ComplexType{}
	collector.put(ns, typeName, addedType)
	res, ok := collector.find(ns, typeName)
	if !ok {
		t.Errorf("Type should exists!")
	}

	resType, ok := res.(*ComplexType)
	if resType != addedType {
		t.Errorf("Types addresses are not equal %p != %p", resType, addedType)
	}
}