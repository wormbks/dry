package str

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveTabsAndNewlines(t *testing.T) {
	output := &bytes.Buffer{}
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{
			input:    []byte("Hello\tWorld\n"),
			expected: []byte("HelloWorld"),
		},
		{
			input:    []byte("Line1\r\nLine2"),
			expected: []byte("Line1Line2"),
		},
		{
			input:    []byte("No\ttabs\ror\nnewlines"),
			expected: []byte("Notabsornewlines"),
		},
		{
			input:    []byte(""),
			expected: []byte(""),
		},
	}

	for _, tc := range testCases {
		err := RemoveTabsAndNewlines(tc.input, output)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, output.Bytes())
		output.Reset()
	}
}
