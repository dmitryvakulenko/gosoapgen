package generating

import "github.com/dmitryvakulenko/gosoapgen/translator"

func Types(structs []*translator.Struct) string {
	res := ""
	for _, s := range structs {
		res += "type " + s.Name + " struct {\n"
		for _, embedType := range s.Embed {
			res += embedType + "\n"
		}
		for _, f := range s.Fields {
			res += f.Name + " " + f.Type + "\n"
		}
		res += "}\n\n"
	}

	return res
}
