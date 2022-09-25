package yamlutil

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestSplitPath(t *testing.T) {
	type testCase struct {
		s              string
		expectedParent string
		expectedBase   string
	}

	testCases := []testCase{
		{
			s:              "$.foo.bar",
			expectedParent: "$.foo",
			expectedBase:   "bar",
		},
	}

	for _, f := range testCases {
		parent, base := SplitPath(f.s)
		t.Logf("s=%q, parent=%q, base=%q", f.s, parent, base)
		assert.Equal(t, f.expectedParent, parent)
		assert.Equal(t, f.expectedBase, base)
	}
}
