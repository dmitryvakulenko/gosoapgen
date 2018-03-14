package generator

import "github.com/dmitryvakulenko/gosoapgen/translator"

func All(typesList []interface{}) string {
	res := ""
	for _, curType := range typesList {
		switch v := curType.(type) {
		case *translator.ComplexType:
			res += "type " + v.Name + " struct {\n"
			for _, f := range v.Fields {
				res += f.Name + " " + f.Type + " `xml:\"" + f.XmlExpr + "\"`\n"
			}
			res += "}\n\n"
		case *translator.SimpleType:
			res += "type " + v.Name + " " + v.Type + "\n\n"
		}

	}

	return res
}
