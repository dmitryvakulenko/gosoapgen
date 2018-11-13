package tree_parser

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/xml"
	xsd "github.com/dmitryvakulenko/gosoapgen/xsd-model"
	"strconv"
)

type typeHash [md5.Size]byte

type Type struct {
	xml.Name
	Fields            []*Field
	UsedByFields      []*Field
	sourceNode        *xsd.Node
	baseType          *Type
	isSimpleContent   bool
	simpleContentType *Type
	// for this type base type fields was resolved
	resolved bool
	// on this element has reference
	referenced bool
}

func (t *Type) addField(f *Field) {
	t.Fields = append(t.Fields, f)
}

func (t *Type) append(addType *Type) {
	t.Fields = append(t.Fields, addType.Fields...)
	t.isSimpleContent = addType.isSimpleContent
}

func (t *Type) hash() typeHash {
	res := []byte(t.Local)
	for _, f := range t.Fields {
		h := f.hash()
		res = append(res, h[:]...)
	}

	return md5.Sum(res)
}

func newType(n *xsd.Node, ns string) *Type {
	name := n.AttributeValue("name")
	return &Type{
		Name:       xml.Name{Local: name, Space: ns},
		sourceNode: n}
}

func newStandardType(name string) *Type {
	return &Type{Name: xml.Name{Local: name, Space: xsdSpace}, isSimpleContent: true}
}

type Field struct {
	Name      string
	Type      *Type
	MinOccurs int
	MaxOccurs int
	IsAttr    bool
	Comment   string
}

func (f *Field) hash() typeHash {
	res := []byte(f.Name)
	// Field type hasn't part it hash generation
	// h := f.Type.hash()
	// res = append(res, h[:]...)

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(f.MinOccurs))
	res = append(res, buf...)
	binary.LittleEndian.PutUint32(buf, uint32(f.MaxOccurs))
	res = append(res, buf...)

	if f.IsAttr {
		res = append(res, 1)
	}

	return md5.Sum(res)
}

func newField(n *xsd.Node, typ *Type) *Field {
	name := n.AttributeValue("name")
	if name == "" {
		name = n.AttributeValue("ref")
	}

	var min int
	switch m := n.AttributeValue("minOccurs"); m {
	case "unqualified", "":
		min = 0
	default:
		min, _ = strconv.Atoi(m)
	}

	var max int
	switch m := n.AttributeValue("maxOccurs"); m {
	case "unbounded":
		max = 1000
	case "":
		max = 0
	default:
		max, _ = strconv.Atoi(m)
	}

	return &Field{
		Name:      name,
		Type:      typ,
		MinOccurs: min,
		MaxOccurs: max}
}

func newXMLNameField() *Field {
	return &Field{
		Name: "XMLName",
		Type: newStandardType("string")}
}

func newValueField(baseType string) *Field {
	return &Field{
		Name: "XMLValue",
		Type: newStandardType(baseType)}
}
