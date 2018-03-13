package translator

type SchemaTypes struct {
	include             []string
	targetNamespace     string
	typesList           []interface{}
	typesListCache      *namespacedTypes
	attributeGroupCache *namespacedTypes
	curXmlns            map[string]string
}

func newSchemaTypes() SchemaTypes {
	return SchemaTypes{
		typesList:           make([]interface{}, 0),
		typesListCache:      newTypesCollection(),
		attributeGroupCache: newTypesCollection(),
		curXmlns:            make(map[string]string)}
}

type ComplexType struct {
	Name      string
	Namespace string
	Fields    []*Field
}

func (t ComplexType) GetNamespace() string {
	return t.Namespace
}

func (t ComplexType) GetName() string {
	return t.Name
}

type SimpleType struct {
	Name      string
	Type      string
	Namespace string
}

func (t SimpleType) GetNamespace() string {
	return t.Namespace
}

func (t SimpleType) GetName() string {
	return t.Name
}

type Field struct {
	Name      string
	Type      string
	XmlExpr   string
	Comment   string
	Namespace string
}

type attributeGroup struct {
	Name      string
	Namespace string
	Fields    []*Field
}

func (t attributeGroup) GetNamespace() string {
	return t.Namespace
}

func (t attributeGroup) GetName() string {
	return t.Name
}
