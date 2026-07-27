package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadResource(mmPath, resourceFile string) (*ResourceConfig, error) {
	resourcePath := mmPath + "/mmv1/products/" + resourceFile
	data, err := os.ReadFile(resourcePath)
	if err != nil {
		return nil, err
	}
	var cfg ResourceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cfg.Properties = ParsePropertiesFromYAML(string(data))

	productPath := filepath.Join(filepath.Dir(resourcePath), "product.yaml")
	if productData, err := os.ReadFile(productPath); err == nil {
		content := string(productData)
		cfg.APIBaseURL = extractGABaseURL(content)
		cfg.PackageName = strings.ToLower(extractField(content, "name"))
	}

	return &cfg, nil
}

func extractField(content, key string) string {
	prefix := key + ":"
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
		}
	}
	return ""
}

func extractGABaseURL(content string) string {
	lines := strings.Split(content, "\n")
	inGA := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "- name: ga" {
			inGA = true
			continue
		}
		if inGA && strings.HasPrefix(trimmed, "base_url:") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "base_url:"))
		}
		if inGA && strings.HasPrefix(trimmed, "- name:") {
			break // moved past ga block
		}
	}
	return ""
}
