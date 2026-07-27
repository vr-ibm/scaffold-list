package templates

import (
	"os"
	"strings"
	"text/template"

	"github.com/vr-ibm/scaffold-list/internal/config"
)

func Render(cfg *config.ResourceConfig, outputDir string) error {
	funcMap := template.FuncMap{
		"toSnake": ToSnake,
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
	var result strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(r | 32)
	}
	return result.String()
}
