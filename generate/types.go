package generate

import (
    "strings"
    "io"
    "github.com/dmitryvakulenko/gosoapgen/xsd/tree_parser"
)

var innerTypes = []string{
    "int",
    "float64",
    "bool",
    "time.Time",
    "string"}

func Types(typesList []*tree_parser.Type, writer io.Writer) {
    for _, curType := range typesList {
        writer.Write([]byte("type " + firstUp(curType.Local) + " struct {\n"))
        for _, f := range curType.Fields {
            writeField(curType, f, writer)
        }
        writer.Write([]byte("}\n\n"))
    }
}

func writeField(t *tree_parser.Type, field *tree_parser.Field, writer io.Writer) {
        // обработка обычного поля
        writer.Write([]byte(firstUp(field.Name) + " "))

        if field.MinOccurs < field.MaxOccurs {
            writer.Write([]byte("[]"))
        }

        fieldType := mapStandardType(field.Type.Local)
        if fieldType == "" {
            fieldType = firstUp(field.Type.Local)
        }

        if !isInnerType(fieldType) {
            writer.Write([]byte("*"))
        }

        writer.Write([]byte(fieldType + " `xml:\""))
        if field.IsAttr {
            writer.Write([]byte(field.Name + ",attr,omitempty"))
        } else if  field.Name == "XMLName" {
            writer.Write([]byte(t.Space + " " + t.Local))
        } else if field.Name == "XMLValue" {
            writer.Write([]byte(",chardata"))
        } else {
            writer.Write([]byte(field.Name + ",omitempty"))
        }
        writer.Write([]byte("\"`\n"))
}

func extractName(in string) string {
    parts := strings.Split(in, ":")
    if len(parts) == 2 {
        return parts[1]
    } else {
        return parts[0]
    }
}

func firstUp(text string) string {
    if isInnerType(text) {
        return text
    }
    return strings.Title(text)
}

func isInnerType(t string) bool {
    for _, v := range innerTypes {
        if v == t {
            return true
        }
    }
    return false
}

func mapStandardType(xmlType string) string {
    switch xmlType {
    case "int", "integer", "positiveInteger", "nonNegativeInteger":
        return "int"
    case "decimal":
        return "float64"
    case "boolean":
        return "bool"
    case "string", "NMTOKEN", "anyURI", "language", "base64Binary", "duration", "IDREF", "IDREFS", "gYear", "gMonth", "gDay", "gYearMonth",
        "date", "dateTime", "time", "ID":
        return "string"
    default:
        return xmlType
    }
}
