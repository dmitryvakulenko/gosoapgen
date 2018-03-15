package generator

import "github.com/dmitryvakulenko/gosoapgen/translator"

func All(typesList []*translator.ComplexType) string {
	res := ""
	for _, curType := range typesList {
		res += "type " + curType.Name + " struct {\n"
		for _, f := range curType.Fields {
			res += f.Name + " " + f.Type + " `xml:\"" + f.XmlExpr + "\"`\n"
		}
		res += "}\n\n"

	}

	return res
}
