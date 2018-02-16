package xsd

import (
	"encoding/xml"
	"io"
)

var (
	types = make(WsdlTypes, 0)
	//lastIndex = 1
	currentStruct *Struct = nil
	currentField *StructField = nil
)

func GetTypes() WsdlTypes {
	return types
}

func Parse(reader io.Reader) {
	decoder := xml.NewDecoder(reader)
	parseImpl(decoder)
}

func parseImpl(decoder *xml.Decoder) {
	t, err := decoder.Token()

	if err == io.EOF {
		return
	}

	switch t := t.(type) {
	case xml.StartElement:
		parseElements(&t)
	case xml.EndElement:
		closeElements(&t)
	}

	parseImpl(decoder)
}

func parseElements(elem *xml.StartElement) {
	switch elem.Name.Local {
	//case "schema":
	//	s := Schema{}
	//	decoder.DecodeElement(s, elem)
	case "element":
		//e := Element{}
		//decoder.DecodeElement(&e, elem)
		if currentStruct == nil {
			currentStruct = newStruct(getAttribute(elem.Attr, "name"))
			currentStruct.Namespace = elem.Name.Space
		} else {
			//currentField = newField(e.Name, e.Type)
		}
	//case "complexType":
	//	c := ComplexType{}
	//	decoder.DecodeElement(c, elem)
	}
}

func closeElements(elem *xml.EndElement) {
	switch elem.Name.Local {
	case "element":
		if currentField != nil {
			currentStruct.Fields = append(currentStruct.Fields, currentField)
			currentField = nil
		} else {
			types = append(types, currentStruct)
			currentStruct = nil
		}
	}
}


func getAttribute(attr []xml.Attr, name string) string {
	for _, v := range attr {
		if v.Name.Local == name {
			return v.Value
		}
	}

	return ""
}