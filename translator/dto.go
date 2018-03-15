package translator

type decoder struct {
	typesList           []*ComplexType
	namespacesList      map[string]bool
	curTargetNamespace  string
	typesListCache      *namespacedTypes
	attributeGroupCache *namespacedTypes
	curXmlns            map[string]string
}

func newDecoder() decoder {
	return decoder{
		typesList:           make([]*ComplexType, 0),
		namespacesList:      make(map[string]bool),
		typesListCache:      newTypesCollection(),
		attributeGroupCache: newTypesCollection(),
		curXmlns:            make(map[string]string)}
}

type ComplexType struct {
	Name      string
	Namespace string
	Fields    []*Field
	BaseType  string
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
	TypeQName string
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
