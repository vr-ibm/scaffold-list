package config

type ResourceConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	BaseURL     string `yaml:"base_url"`
	IDFormat    string `yaml:"id_format"`
	Kind        string `yaml:"kind"`
	Properties  []Property
	APIBaseURL  string // from product.yaml GA version base_url
	PackageName string // lowercased product name (e.g. "dns")
}
