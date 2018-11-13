package tree_parser

import "encoding/xml"

type typesList struct {
	cache    map[xml.Name]int
	fullList []*Type
}

func newTypesList() *typesList {
	return &typesList{cache: make(map[xml.Name]int)}
}

func (t *typesList) Reset() {
	t.cache = make(map[xml.Name]int)
}

func (t *typesList) GetList() []*Type {
	return t.fullList
}

func (t *typesList) Has(name xml.Name) bool {
	_, res := t.cache[name]
	return res
}

func (t *typesList) Add(newType *Type) {
	if t.Has(newType.Name) {
		panic("Type exist " + newType.Local)
	}

	index := len(t.fullList)
	t.fullList = append(t.fullList, newType)
	t.cache[newType.Name] = index
}

func (t *typesList) Get(name xml.Name) *Type {
	index, ok := t.cache[name]

	if !ok {
		panic("Type isn't exist " + name.Local)
	}

	return t.fullList[index]
}