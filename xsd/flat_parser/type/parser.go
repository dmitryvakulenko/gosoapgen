package _type

import (
	"encoding/xml"
	"path"
)

// Интерфейс загрузки xsd
// должен отслеживать уже загруженные файлы
// и правильно отрабатывать относительные пути
type Loader interface {
	// Если второй параметр true, это означает
	// что такой файл уже был загружен
	Load(path string) ([]byte, error)
	IsAlreadyLoadedError(error) bool
}


// Parser for xsd schemas
type Parser struct {
	loader Loader
}

func NewParser(l Loader) *Parser {
	return &Parser{
		loader: l}
}


// Returns all loaded schemas (included and imported)
func (p *Parser) Parse(schemaFileName string, ns string) []*Schema {
	xsdData, err := p.loader.Load(schemaFileName)
	if p.loader.IsAlreadyLoadedError(err) {
		return make([]*Schema, 0)
	}

	curSchema := p.unmarshalXsd(xsdData)
	if ns != "" {
		curSchema.TargetNamespace = ns
	}
	res := p.parseImports(curSchema, ns, path.Dir(schemaFileName))
	// Здесь важен порядок - основная схема должна идти последней
	res = append(res, curSchema)

	return res
}

// Загрузить схему с помощью Loader-а и распасить в Schema
func (p *Parser) unmarshalXsd(data []byte) *Schema {
	s := Schema{}
	err := xml.Unmarshal(data, &s)
	if err != nil {
		panic(err)
	}

	return &s
}

// Parsing imports and includes
func (p *Parser) parseImports(s *Schema, ns string, baseFilePath string) []*Schema {
	targetNamespace := ns
	if targetNamespace == "" {
		targetNamespace = s.TargetNamespace
	}

	res := make([]*Schema, 0)
	for _, imp := range s.Import {
		curNs := targetNamespace
		if imp.Namespace != "" {
			curNs = imp.Namespace
		}
		res = append(res, p.Parse(baseFilePath + "/" + imp.SchemaLocation, curNs)...)
	}

	for _, imp := range s.Include {
		res = append(res, p.Parse(baseFilePath + "/" + imp.SchemaLocation, targetNamespace)...)
	}

	return res
}


