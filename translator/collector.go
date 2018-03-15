package translator

import "log"

type typesCollection map[string]interface{}
type namespacedTypes map[string]*typesCollection

type Namespaceable interface {
	GetNamespace() string
	GetName() string
}

func newTypesCollection() *namespacedTypes {
	res := make(namespacedTypes)
	return &res
}

func (t *namespacedTypes) find(namespace, typeName string) (interface{}, bool) {
	var ns *typesCollection
	var ok bool
	var curType interface{}
	if ns, ok = (*t)[namespace]; !ok {
		return nil, false
	}

	if curType, ok = (*ns)[typeName]; !ok {
		return nil, false
	}

	return curType, true
}

func (t *namespacedTypes) put(addedType Namespaceable) {
	var (
		namespace = addedType.GetNamespace()
		typeName = addedType.GetName()
	)

	if _, ok := (*t)[namespace]; !ok {
		newCollection := make(typesCollection)
		(*t)[namespace] = &newCollection
	}

	ns := (*t)[namespace]
	if _, ok := (*ns)[typeName]; ok {
		log.Panicf("Namespace %q already contain type %q", namespace, typeName)
	}

	(*ns)[typeName] = addedType
}

func (t *namespacedTypes) has(addedType Namespaceable) bool {
	var (
		namespace = addedType.GetNamespace()
		typeName = addedType.GetName()
	)

	if _, ok := (*t)[namespace]; !ok {
		return false
	}

	ns := (*t)[namespace]
	if _, ok := (*ns)[typeName]; !ok {
		return false
	}

	return true
}

func (t *namespacedTypes) getAllTypes() []interface{} {
	var res []interface{}

	for _, nsList := range *t {
		for _, curType := range *nsList {
			res = append(res, curType)
		}
	}

	return res
}

func (t *namespacedTypes) Merge(newTypes namespacedTypes) {
	for ns, nsList := range newTypes {
		if _, ok := (*t)[ns]; !ok {
			newCollection := make(typesCollection)
			(*t)[ns] = &newCollection
		}

		curCollection := (*t)[ns]
		for typeName, value := range *nsList {
			(*curCollection)[typeName] = value
		}
	}
}