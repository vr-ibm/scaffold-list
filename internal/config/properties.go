package config

import (
	"bufio"
	"strings"
)

type Property struct {
	Name       string
	Type       string
	Required   bool
	Properties []Property // for NestedObject
}

// parseProperties parses a slice of lines that follow a `properties:` key.
// baseIndent is the column at which `- name:` entries are expected.
func parseProperties(lines []string, baseIndent int) []Property {
	var props []Property
	i := 0
	for i < len(lines) {
		line := lines[i]
		indent := countIndent(line)
		trimmed := strings.TrimSpace(line)

		// Match `- name: <value>` at the expected indent level
		if indent == baseIndent && strings.HasPrefix(trimmed, "- name:") {
			name := strings.TrimSpace(strings.TrimPrefix(trimmed, "- name:"))
			prop := Property{Name: name}

			// Scan subsequent lines for fields at baseIndent+2
			j := i + 1
			for j < len(lines) {
				fieldLine := lines[j]
				fieldIndent := countIndent(fieldLine)
				fieldTrimmed := strings.TrimSpace(fieldLine)

				// A non-empty line back at baseIndent (or less) means a new sibling entry
				if fieldTrimmed != "" && fieldIndent <= baseIndent {
					break
				}

				if fieldIndent == baseIndent+2 {
					if strings.HasPrefix(fieldTrimmed, "type:") {
						prop.Type = strings.TrimSpace(strings.TrimPrefix(fieldTrimmed, "type:"))
					}
					if strings.HasPrefix(fieldTrimmed, "required:") {
						prop.Required = strings.TrimSpace(strings.TrimPrefix(fieldTrimmed, "required:")) == "true"
					}
					if fieldTrimmed == "properties:" && prop.Type == "NestedObject" {
						// Collect nested lines until we return to baseIndent+2 or less
						k := j + 1
						var nestedLines []string
						for k < len(lines) {
							if strings.TrimSpace(lines[k]) != "" && countIndent(lines[k]) <= baseIndent+2 {
								break
							}
							nestedLines = append(nestedLines, lines[k])
							k++
						}
						prop.Properties = parseProperties(nestedLines, baseIndent+4)
						j = k
						continue
					}
				}
				j++
			}
			props = append(props, prop)
			i = j
			continue
		}
		i++
	}
	return props
}

func countIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

// ParsePropertiesFromYAML extracts the top-level properties block from raw YAML content.
func ParsePropertiesFromYAML(content string) []Property {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	inProps := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "properties:" {
			inProps = true
			continue
		}
		if inProps && len(line) > 0 && line[0] != ' ' {
			break
		}
		if inProps {
			lines = append(lines, line)
		}
	}
	// Top-level `- name:` entries are at indent 2
	return parseProperties(lines, 2)
}
