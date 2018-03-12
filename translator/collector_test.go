package translator

import "testing"

func TestGetTypeFromEmptyCollection(t *testing.T) {
	collector := newTypesCollection()

	curType, ok := collector.find("aa", "bbbb")
	if ok {
		t.Errorf("Should be error")
	}

	if curType != nil {
		t.Errorf("Type should be nil")
	}
}

func TestStoringType(t *testing.T) {
	collector := newTypesCollection()

	ns := "namespace"
	typeName := "ComplexType"
	addedType := &ComplexType{Name: typeName, Namespace: ns}
	collector.put(addedType)
	res, ok := collector.find(ns, typeName)
	if !ok {
		t.Errorf("Type should exists!")
	}

	resType, ok := res.(*ComplexType)
	if resType != addedType {
		t.Errorf("Types addresses are not equal %p != %p", resType, addedType)
	}
}