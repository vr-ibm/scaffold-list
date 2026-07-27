package config

import (
	"fmt"
	"os"
	"testing"
)

func TestParseProperties(t *testing.T) {
	data, _ := os.ReadFile(os.ExpandEnv("$HOME/Documents/github/magic-modules/mmv1/products/dns/ManagedZone.yaml"))
	props := ParsePropertiesFromYAML(string(data))
	if len(props) == 0 {
		t.Fatal("expected properties, got none")
	}
	fmt.Printf("parsed %d top-level properties:\n", len(props))
	for _, p := range props {
		fmt.Printf("  name=%-20s type=%-15s required=%v nested=%d\n", p.Name, p.Type, p.Required, len(p.Properties))
	}
}
