package translator

type typesCollection map[string]*interface{}
type namespacedTypes map[string]*typesCollection

type xsdTypes struct {
	typesList namespacedTypes
}

func makeTypesCollection() *xsdTypes {
	return &xsdTypes{}
}

func (t *xsdTypes) find(namespace, typeName string) (*interface{}, bool) {
	var ns *typesCollection
	var ok bool
	var curType *interface{}
	if ns, ok = t.typesList[namespace]; !ok {
		return nil, false
	}

	if curType, ok = (*ns)[typeName]; !ok {
		return nil, false
	}

	return curType, true
}
