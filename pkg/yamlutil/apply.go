package yamlutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

// OpType is the operator type.
type OpType string

const (
	// OpInvalid is invalid.
	OpInvalid = OpType("")
	// OpSet sets a property.
	OpSet = OpType("set")
	// TODO: OpRemove, OpAppend, ...
)

// Op is an op.
type Op struct {
	Type  OpType
	Path  string // like $.ssh.localPort
	Value string
}

// Apply applies ops.
func Apply(b []byte, ops ...Op) ([]byte, error) {
	header, body, footer, err := splitCommentOnlyHeaderFooter(b)
	if err != nil {
		return nil, fmt.Errorf("failed to split comment-only headers and footers: %w", err)
	}
	body, err = applyOps(body, ops...)
	if err != nil {
		return nil, err
	}
	if body[len(body)-1] != byte('\n') {
		body = append(body, byte('\n'))
	}
	return append(append(header, body...), footer...), nil
}

func applyOps(b []byte, ops ...Op) ([]byte, error) {
	astFile, err := parser.ParseBytes(b, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the YAML: %w", err)
	}
	for i, op := range ops {
		if err = Apply1(astFile, op); err != nil {
			return nil, fmt.Errorf("failed to apply op %d/%d (%q): %w", i+1, len(ops), op.Path, err)
		}
	}
	return []byte(astFile.String()), nil
}

// Apply1 applies a single op.
func Apply1(astFile *ast.File, op Op) error {
	if len(astFile.Docs) != 1 {
		return fmt.Errorf("multi-document YAML is unsupported yet: the YAML contains %d documents", len(astFile.Docs))
	}
	if op.Type != OpSet {
		return fmt.Errorf("unexpected op: %v", op.Type)
	}
	path, err := yaml.PathString(op.Path)
	if err != nil {
		return fmt.Errorf("failed to parse the YAML path %q: %w", op.Path, err)
	}
	if _, err := parser.ParseBytes([]byte(op.Value), parser.ParseComments); err != nil {
		return fmt.Errorf("failed to parse the value %q: %w", op.Value, err)
	}
	if _, err := path.FilterFile(astFile); err == nil {
		if err := path.ReplaceWithReader(astFile, strings.NewReader(op.Value)); err != nil {
			return fmt.Errorf("failed to call ReplaceWithReader for the YAML path %q: %w", path, err)
		}
		return nil
	} else if !yaml.IsNotFoundNodeError(err) {
		return fmt.Errorf("failed to get a node at %q: %w", path, err)
	}
	return mergeMappingNode(astFile, path, op.Value)
}

func mergeMappingNode(astFile *ast.File, path *yaml.Path, value string) error {
	parent, base := SplitPath(path.String())
	if parent == "" || base == "" {
		return fmt.Errorf("unexpected YAML path %q (parent=%q, base=%q)", path, parent, base)
	}
	snippet := "{\n  " + base + ": " + value + "\n}"
	if parent == "$" {
		snippetAstFile, err := parser.ParseBytes([]byte(snippet), parser.ParseComments)
		if err != nil {
			return err
		}
		return ast.Merge(astFile.Docs[0], snippetAstFile.Docs[0])
	}
	parentPath, err := yaml.PathString(parent)
	if err != nil {
		return fmt.Errorf("failed to parse %q (parent of %q): %w", parent, path, err)
	}
	if err = ensureMappingNode(astFile, parentPath); err != nil {
		return fmt.Errorf("failed to ensure *ast.MappingNode on %q (parent of %q): %w", parent, path, err)
	}
	if err = parentPath.MergeFromReader(astFile, strings.NewReader(snippet)); err != nil {
		return fmt.Errorf("failed to merge %q into Mapping %q (parent of %q): %w", snippet, parent, path, err)
	}
	return nil
}

func getMappingNode(astFile *ast.File, path *yaml.Path) (*ast.MappingNode, error) {
	node, err := path.FilterFile(astFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get a node at %q: %w", path, err)
	}
	mappingNode, ok := node.(*ast.MappingNode)
	if !ok {
		return nil, fmt.Errorf("expected %q to be *ast.MappingNode, got %T", path, node)
	}
	return mappingNode, nil
}

func ensureMappingNode(astFile *ast.File, path *yaml.Path) error {
	if path.String() == "$" {
		return errors.New("ensureMappingNode: unexpected YAML path \"$\"")
	}
	if _, err := getMappingNode(astFile, path); err == nil {
		return nil
	} else if !yaml.IsNotFoundNodeError(err) {
		return fmt.Errorf("failed to query %q: %w", path, err)
	}
	return mergeMappingNode(astFile, path, "{}")
}

// Editable returns nil error if the YAML is safely editable.
func Editable(b []byte) error {
	_, body, _, err := splitCommentOnlyHeaderFooter(b)
	if err != nil {
		return fmt.Errorf("failed to split comment-only headers and footers: %w", err)
	}

	dummyOp := Op{
		Type:  OpSet,
		Path:  "$.yamlctl.internal.test-editable",
		Value: "null",
	}

	applied, err := applyOps(body, dummyOp)
	if err != nil {
		return fmt.Errorf("failed to apply %+v: %w", dummyOp, err)
	}

	expected := string(body)
	if body[len(body)-1] != byte('\n') {
		expected += "\n"
	}
	expected += "yamlctl: {internal: {test-editable: null}}"
	if string(applied) != expected {
		diff := cmp.Diff(string(applied), expected)
		logrus.WithField("op", dummyOp).Debug("Diff: " + diff)
		return errors.New("the YAML is not safely editable (reason is printed in the debug log)")
	}
	return nil
}
