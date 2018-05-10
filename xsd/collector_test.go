package xsd

import "testing"

func TestGetTypeFromEmptyCollection(t *testing.T) {
	collector := newTypesCollection()

	curType, ok := collector.find("aa", "bbbb")
	if ok {
		t.Errorf("Should be error")
	}

	if curType != nil {
		t.Errorf("TypeName should be nil")
	}
}

func TestStoringType(t *testing.T) {
	collector := newTypesCollection()

	ns := "namespace"
	typeName := "ComplexType"
	addedType := &ComplexType{GoName: typeName, Namespace: ns}
	collector.put(ns, addedType)
	res, ok := collector.find(ns, typeName)
	if !ok {
		t.Errorf("TypeName should exists!")
	}

	resType, ok := res.(*ComplexType)
	if resType != addedType {
		t.Errorf("Import addresses are not equal %p != %p", resType, addedType)
	}
}