package dom_parser

import "testing"

func TestEmptySchema(t *testing.T) {
	typesList := loadXsd(t.Name())

	if len(typesList) != 0 {
		t.Errorf("Should be no types")
	}
}


func loadXsd(testName string) []*Type {
	p := NewParser()
	p.LoadFile("./test_data/" + testName + ".xsd")
	return p.GetTypes()
}
