package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadResource(mmPath, resourceFile string) (*ResourceConfig, error) {
	data, err := os.ReadFile(mmPath + "/mmv1/products/" + resourceFile)
	if err != nil {
		return nil, err
	}
	var cfg ResourceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cfg.Properties = ParsePropertiesFromYAML(string(data))
	return &cfg, nil
}
