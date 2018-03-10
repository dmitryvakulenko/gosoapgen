package translator

import "log"

type typesCollection map[string]interface{}
type namespacedTypes map[string]*typesCollection

type xsdTypes struct {
	typesList namespacedTypes
}

type Namespaceable interface {
	GetNamespace() string
	GetName() string
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

func (t *xsdTypes) put(addedType Namespaceable) {
	var (
		namespace = addedType.GetNamespace()
		typeName = addedType.GetName()
	)

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


func (t *xsdTypes) getAllTypes() []interface{} {
	var res []interface{}

	for _, nsList := range t.typesList {
		for _, curType := range *nsList {
			res = append(res, curType)
		}
	}

	return res
}