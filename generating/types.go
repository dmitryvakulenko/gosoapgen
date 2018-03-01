package generating

import "gosoapgen/translator"

func Types(structs []*translator.Struct) string {
	res := ""
	for _, s := range structs {
		res += "type " + s.Name + " struct {\n"
		for _, f := range s.Fields {
			res += f.Name + " " + f.Type + "\n"
		}
		res += "}\n\n"
	}

	return res
}
