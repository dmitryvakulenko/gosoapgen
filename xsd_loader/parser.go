package xsd_loader

func NewParser(loader Loader) *Parser {
	return &Parser{loader: loader}
}

type Parser struct {
	loader Loader
}


func (p *Parser) Parse(fileName string) *Schema {
	return &Schema{}
}