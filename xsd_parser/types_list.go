package xsd_parser

import "encoding/xml"

type typesList struct {
    cache       map[xml.Name]int
    goNameCache map[string]bool
    fullList    []*Type
}

func newTypesList() *typesList {
    return &typesList{
        cache:       make(map[xml.Name]int),
        goNameCache: make(map[string]bool)}
}

func (t *typesList) Reset() {
    t.cache = make(map[xml.Name]int)
}

func (t *typesList) Iterate() []*Type {
    return t.fullList
}

func (t *typesList) Has(name xml.Name) bool {
    _, res := t.cache[name]
    return res
}

func (t *typesList) HasGoName(name string) bool {
    _, res := t.goNameCache[name]
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

func (t *typesList) Remove(remType *Type) {
    if !t.Has(remType.Name) {
        panic("No such type " + remType.Local)
    }
    index, _ := t.cache[remType.Name]
    delete(t.cache, remType.Name)
    t.fullList = append(t.fullList[0:index], t.fullList[index+1:]...)
}
