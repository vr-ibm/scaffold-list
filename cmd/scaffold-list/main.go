package main

import (
	"flag"
	"log"

	"github.com/vr-ibm/scaffold-list/internal/config"
	"github.com/vr-ibm/scaffold-list/internal/templates"
)

func main() {
	resource := flag.String("resource", "", "Resource name (e.g. dns/ManagedZone.yaml)")
	mmPath := flag.String("magic-modules", "~/Documents/github/magic-modules", "Path to local magic-modules clone")
	output := flag.String("output", ".", "Output directory for generated files")
	flag.Parse()

	cfg, err := config.LoadResource(*mmPath, *resource)
	if err != nil {
		log.Fatalf("failed to load resource: %v", err)
	}
	if err := templates.Render(cfg, *output); err != nil {
		log.Fatalf("failed to render template: %v", err)
	}
	log.Printf("generated: %s/data_source_%s_list.go", *output, templates.ToSnake(cfg.Name))
}
