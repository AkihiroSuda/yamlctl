package yamlutil

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

// Query returns the filtered YAML.
func Query(b []byte, pathStr string) ([]byte, error) {
	astFile, err := parser.ParseBytes(b, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the YAML: %w", err)
	}
	path, err := yaml.PathString(pathStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the YAML path %q: %w", pathStr, err)
	}
	n, err := path.FilterFile(astFile)
	if err != nil {
		return nil, fmt.Errorf("failed to query %q: %w", path, err)
	}
	return n.MarshalYAML()
}
