package translator

import "log"

type typesCollection map[string]interface{}
type namespacedTypes map[string]*typesCollection

type xsdTypes struct {
	typesList namespacedTypes
}

func makeTypesCollection() *xsdTypes {
	res := &xsdTypes{
		typesList: make(namespacedTypes)}

	return res
}

func (t *xsdTypes) find(namespace, typeName string) (interface{}, bool) {
	var ns *typesCollection
	var ok bool
	var curType interface{}
	if ns, ok = t.typesList[namespace]; !ok {
		return nil, false
	}

	if curType, ok = (*ns)[typeName]; !ok {
		return nil, false
	}

	return curType, true
}


func (t *xsdTypes) put(namespace, typeName string, addedType interface{}) {
	if _, ok := t.typesList[namespace]; !ok {
		newCollection := make(typesCollection)
		t.typesList[namespace] = &newCollection
	}

	ns := t.typesList[namespace]
	if _, ok := (*ns)[typeName]; ok {
		log.Panicf("Namespace %q already contain type %q", namespace, typeName)
	}

	(*ns)[typeName] = addedType
}