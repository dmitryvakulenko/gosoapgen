package flat_parser

import "testing"

func TestGetTypeFromEmptyCollection(t *testing.T) {
	collector := NewTypesCollection()

	curType, ok := collector.Find("aa", "bbbb")
	if ok {
		t.Errorf("Should be error")
	}

	if curType != nil {
		t.Errorf("TypeName should be nil")
	}
}

func TestStoringType(t *testing.T) {
	collector := NewTypesCollection()

	ns := "namespace"
	typeName := "ComplexType"
	addedType := &ComplexType{GoName: typeName, Namespace: ns}
	collector.Put(ns, addedType)
	res, ok := collector.Find(ns, typeName)
	if !ok {
		t.Errorf("TypeName should exists!")
	}

	resType, ok := res.(*ComplexType)
	if resType != addedType {
		t.Errorf("Import addresses are not equal %p != %p", resType, addedType)
	}
}