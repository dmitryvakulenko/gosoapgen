package xsd_loader

type (
	Schema struct {
		Elements        []*Element
		Types           []*Type
		Attributes      []*Attribute
		AttributeGroups []*AttributeGroup
	}

	Element struct {
		Name     string
		Type     *Type
		typeName string
	}

	Type struct {
		Name         string
		BaseType     *Type
		baseTypeName string
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
