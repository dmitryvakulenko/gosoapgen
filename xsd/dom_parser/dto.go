package dom_parser

import "encoding/xml"

// import (
//     dom "github.com/subchen/go-xmldom"
// )

type Type struct {
    xml.Name
    Fields    []*Field

    // Если это не simpleType, значит создано из extension/restriction
    BaseType     *Type
    BaseTypeName xml.Name
}

func (t *Type) addField(f *Field) {
    t.Fields = append(t.Fields, f)
}

type Field struct {
    Name     string
    Type     *Type
    TypeName xml.Name
    IsArray  bool
    IsAttr   bool
    Comment  string
}

func newField(name string, typ *Type) *Field {
    return &Field{
        Name: name,
        Type: typ}
}