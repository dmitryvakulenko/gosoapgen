package translator

import "log"

type typesCollection map[string]NamedType
type namespacedTypes map[string]*typesCollection

func newTypesCollection() *namespacedTypes {
	res := make(namespacedTypes)
	return &res
}

func (t *namespacedTypes) find(namespace, typeName string) (NamedType, bool) {
	var (
		ns *typesCollection
		ok bool
		curType NamedType
	)
	if ns, ok = (*t)[namespace]; !ok {
		return nil, false
	}

	if curType, ok = (*ns)[typeName]; !ok {
		return nil, false
	}

	return curType, true
}

func (t *namespacedTypes) put(namespace string, addedType NamedType) {
	if _, ok := (*t)[namespace]; !ok {
		newCollection := make(typesCollection)
		(*t)[namespace] = &newCollection
	}

	typeName := addedType.GetName()
	ns := (*t)[namespace]
	if _, ok := (*ns)[typeName]; ok {
		log.Panicf("Namespace %q already contain type %q", namespace, typeName)
	}

	(*ns)[typeName] = addedType
}

func (t *namespacedTypes) has(namespace, typeName string) bool {
	if _, ok := (*t)[namespace]; !ok {
		return false
	}

	ns := (*t)[namespace]
	if _, ok := (*ns)[typeName]; !ok {
		return false
	}

	return true
}

func (t *namespacedTypes) getAllTypes() []NamedType {
	var res []NamedType

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