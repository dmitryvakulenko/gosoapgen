package tree_parser

import "log"

type typesCollection map[string]*Type
type NamespacedTypes map[string]*typesCollection

func NewTypesCollection() *NamespacedTypes {
	res := make(NamespacedTypes)
	return &res
}

func (t *NamespacedTypes) Find(namespace, typeName string) (*Type, bool) {
	var (
		ns *typesCollection
		ok bool
		curType *Type
	)
	if ns, ok = (*t)[namespace]; !ok {
		return nil, false
	}

	if curType, ok = (*ns)[typeName]; !ok {
		return nil, false
	}

	return curType, true
}

func (t *NamespacedTypes) Put(namespace string, addedType *Type) {
	if _, ok := (*t)[namespace]; !ok {
		newCollection := make(typesCollection)
		(*t)[namespace] = &newCollection
	}

	typeName := addedType.Name
	ns := (*t)[namespace]
	if _, ok := (*ns)[typeName]; ok {
		log.Panicf("Namespace %q already contain type %q", namespace, typeName)
	}

	(*ns)[typeName] = addedType
}

func (t *NamespacedTypes) Has(namespace, typeName string) bool {
	if _, ok := (*t)[namespace]; !ok {
		return false
	}

	ns := (*t)[namespace]
	if _, ok := (*ns)[typeName]; !ok {
		return false
	}

	return true
}

func (t *NamespacedTypes) GetAllTypes() []*Type {
	var res []*Type

	for _, nsList := range *t {
		for _, curType := range *nsList {
			res = append(res, curType)
		}
	}

	return res
}

func (t *NamespacedTypes) Merge(newTypes NamespacedTypes) {
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