package yamlutil

import (
	"testing"

	"github.com/goccy/go-yaml"
	"gotest.tools/v3/assert"
)

func TestApply(t *testing.T) {
	const orig = `# begin YAML
---
# begin foo
foo:
  # here is a string
  existingString: "apple"
  # some multi-line comments
  # blah
  existingNumber: 42 # here is a number (this in-line comment will be lost)
# end foo
# begin bar
bar:
  existingString: "alpha"
# end bar
baz: {}
# end YAML`
	ops := []Op{
		{
			Type:  OpSet,
			Path:  "$.foo.existingString",
			Value: "\"banana\"",
		},
		{
			Type:  OpSet,
			Path:  "$.foo.newString",
			Value: "\"chocolate\"",
		},
		{
			Type:  OpSet,
			Path:  "$.foo.existingNumber",
			Value: "43",
		},
		{
			Type:  OpSet,
			Path:  "$.foo.newNumber",
			Value: "44 # new number",
		},
		{
			Type:  OpSet,
			Path:  "$.bar.existingString",
			Value: "\"beta\"",
		},
		{
			Type:  OpSet,
			Path:  "$.foo.newMapping.number",
			Value: "45",
		},
		{
			Type:  OpSet,
			Path:  "$.foo.newMapping.float",
			Value: "45.1",
		},

		// FIXME
		{
			Type:  OpSet,
			Path:  "$.newMapping.number",
			Value: "1234",
		},
	}

	const expected = `# begin YAML
---
# begin foo
foo:
  # here is a string
  existingString: "banana"
  # some multi-line comments
  # blah
  existingNumber: 43
  newString: "chocolate"
  newNumber: 44
  newMapping: {number: 45, float: 45.1}
# end foo
# begin bar
bar:
  existingString: "beta"
# end bar
baz: {}
newMapping: {number: 1234}
# end YAML`

	t.Log("==> orig <==")
	t.Log(orig)
	assert.NilError(t, Editable([]byte(orig)))

	t.Log("==> ops <==")
	t.Log(ops)

	got, err := Apply([]byte(orig), ops...)
	assert.NilError(t, err)
	gotS := string(got)
	t.Log("==> got <==")
	t.Log(gotS)

	gotJ, err := yaml.YAMLToJSON(got)
	assert.NilError(t, err)
	t.Log("==> got (JSON representation) <==")
	t.Log(string(gotJ))

	assert.Equal(t, expected, gotS)
}

func TestApplyMultiDoc(t *testing.T) {
	const orig = `
---
foo:
  num: 42
---
foo:
  num: 43`
	ops := []Op{
		{
			Type:  OpSet,
			Path:  "$.foo.num",
			Value: "100",
		},
	}
	_, err := Apply([]byte(orig), ops...)
	assert.ErrorContains(t, err, "multi-document YAML is unsupported yet: the YAML contains 2 documents")
}
