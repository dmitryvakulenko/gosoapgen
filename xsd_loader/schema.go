package xsd_loader

import "encoding/xml"

type (
	Type interface {
		GetName() xml.Name
	}

	Schema struct {
		Elements        []*Element
		Types           []Type
		Attributes      []*Attribute
		AttributeGroups []*AttributeGroup
	}

	Element struct {
		Name     xml.Name
		Type     Type
		typeName xml.Name
	}

	SimpleType struct {
		Name         xml.Name
		BaseType     *SimpleType
		baseTypeName xml.Name
	}

	ComplexType struct {
		Name xml.Name
	}

	Attribute struct {
	}

	AttributeGroup struct {
	}
)

func (s *Schema) addElement(e *Element) {
	s.Elements = append(s.Elements, e)
}

func (s *Schema) addType(t Type) {
	s.Types = append(s.Types, t)
}

func (s *SimpleType) GetName() xml.Name {
	return s.Name
}

func (s *ComplexType) GetName() xml.Name {
	return s.Name
}
