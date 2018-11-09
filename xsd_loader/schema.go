package xsd_loader

import "encoding/xml"

type (
	Schema struct {
		Elements        []*Element
		Types           []*Type
		Attributes      []*Attribute
		AttributeGroups []*AttributeGroup
	}

	Element struct {
		Name     xml.Name
		Type     *Type
		typeName xml.Name
	}

	Type struct {
		Name         xml.Name
		BaseType     *Type
		baseTypeName xml.Name
	}

	Attribute struct {
	}

	AttributeGroup struct {
	}
)

func (s *Schema) addElement(e *Element) {
	s.Elements = append(s.Elements, e)
}

func (s *Schema) addType(t *Type) {
	s.Types = append(s.Types, t)
}
