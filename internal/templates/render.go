package templates

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/vr-ibm/scaffold-list/internal/config"
)

//go:embed list_resource.go.tmpl
var templateFS embed.FS

func Render(cfg *config.ResourceConfig, outputDir string) error {
	funcMap := template.FuncMap{
		"toSnake": ToSnake,
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
		// toFmtURL converts a MM-style URL (e.g. "projects/{{project}}/zones")
		// into a fmt.Sprintf format string (e.g. "projects/%s/zones").
		"toFmtURL": func(s string) string {
			s = strings.ReplaceAll(s, "{{project}}", "%s")
			s = strings.ReplaceAll(s, "{{name}}", "%s")
			return s
		},
		// toLowerCamel lowercases just the first rune (ManagedZone → managedZone).
		"toLowerCamel": func(s string) string {
			if s == "" {
				return s
			}
			r := []rune(s)
			r[0] = unicode.ToLower(r[0])
			return string(r)
		},
	}

	tmplContent, err := templateFS.ReadFile("list_resource.go.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("list_resource.go.tmpl").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	outFile := filepath.Join(outputDir, "data_source_"+ToSnake(cfg.Name)+"_list.go")
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
