package yamlutil

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestSplitCommentOnlyHeaderFooter(t *testing.T) {
	type testCase struct {
		b              []byte
		expectedHeader []byte
		expectedBody   []byte
		expectedFooter []byte
	}

	testCases := []testCase{
		{
			b:              []byte(`foo: 42`),
			expectedHeader: nil,
			expectedBody:   []byte(`foo: 42`),
			expectedFooter: nil,
		},
		{
			b: []byte(`# single-line header
foo: 42
# single-line footer
`),
			expectedHeader: []byte(`# single-line header
`),
			expectedBody: []byte(`foo: 42
`),
			expectedFooter: []byte(`# single-line footer
`),
		},
		{
			b: []byte(`# multi-line header
#aaa
#  weird indent
# aaa
foo: 42
bar: 43
# footer without LF`),
			expectedHeader: []byte(`# multi-line header
#aaa
#  weird indent
# aaa
`),
			expectedBody: []byte(`foo: 42
bar: 43
`),
			expectedFooter: []byte(`# footer without LF`),
		},
	}

	for _, f := range testCases {
		header, body, footer, err := splitCommentOnlyHeaderFooter(f.b)
		assert.NilError(t, err, string(f.b))
		t.Logf("b=%q, header=%q, body=%q, footer=%q", string(f.b), string(header), string(body), string(footer))
		assert.Equal(t, string(f.expectedHeader), string(header), "header")
		assert.Equal(t, string(f.expectedBody), string(body), "body")
		assert.Equal(t, string(f.expectedFooter), string(footer), "footer")
	}
}
