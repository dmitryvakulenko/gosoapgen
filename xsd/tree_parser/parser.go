/*
Парсер схемы xsd.
Использует ручной обход всего дерева, вместо загрузки его с помощь Unmarshal.
 */
package tree_parser

// Интерфейс загрузки xsd
// должен отслеживать уже загруженные файлы
// и правильно отрабатывать относительные пути
type Loader interface {
	// Если второй параметр true, это означает
	// что такой файл уже был загружен
	Load(path string) ([]byte, error)
	IsAlreadyLoadedError(error) bool
}

type parser struct {
	loader Loader
}

func NewParser(l Loader) *parser {
	return &parser{
		loader: l}
}

func (p *parser) Parse(fileName string) {

}

func (p *parser) GetTypes() []*NamedType {
	return make([]*NamedType, 0)
}

