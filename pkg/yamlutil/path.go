package yamlutil

import (
	"strings"
)

// SplitPath splits a yaml path into the parent and the base.
func SplitPath(s string) (string, string) {
	idx := strings.LastIndex(s, ".")
	if idx == -1 {
		return "", s
	}
	return s[:idx], s[idx+1:]
}
