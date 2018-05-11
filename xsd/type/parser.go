package _type

import "encoding/xml"

// Интерфейс загрузки xsd
// должен отслеживать уже загруженные файлы
// и правильно отрабатывать относительные пути
type Loader interface {
	// Если второй параметр true, это означает
	// что такой файл уже был загружен
	Load(path string) ([]byte, bool)
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
	res := make([]*Schema, 0)
	xsdData, parsed := p.loader.Load(schemaFileName)
	if parsed {
		return res
	}
	res = []*Schema{p.unmarshalXsd(xsdData)}
	if ns != "" {
		res[0].TargetNamespace = ns
	}
	res = append(res, p.parseImports(res[0], ns)...)

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
func (p *Parser) parseImports(s *Schema, ns string) []*Schema {
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
		res = append(res, p.Parse(imp.SchemaLocation, curNs)...)
	}

	for _, imp := range s.Include {
		res = append(res, p.Parse(imp.SchemaLocation, targetNamespace)...)
	}

	return res
}


