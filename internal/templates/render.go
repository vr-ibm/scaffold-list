package templates

import (
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/vr-ibm/scaffold-list/internal/config"
)

func Render(cfg *config.ResourceConfig, outputDir string) error {
	funcMap := template.FuncMap{
		"toSnake": func(s string) string { return ToSnake(s) },
		"toSchemaType": func(t string) string {
			switch t {
			case "Boolean":
				return "schema.TypeBool"
			case "Integer":
				return "schema.TypeInt"
			case "Array", "NestedObject":
				return "schema.TypeList"
			default:
				return "schema.TypeString"
			}
		},
	}
	tmpl, err := template.New("list_resource.go.tmpl").Funcs(funcMap).ParseFiles("internal/templates/list_resource.go.tmpl")
	if err != nil {
		return err
	}
	outFile := outputDir + "/data_source_" + ToSnake(cfg.Name) + "_list.go"
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, cfg)
}

func ToSnake(s string) string {
	// normalise known acronyms to title-case so they produce a single word
	s = strings.ReplaceAll(s, "ID", "Id")
	s = strings.ReplaceAll(s, "DNS", "Dns")
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
