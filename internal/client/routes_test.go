package client

import (
	"testing"
)

func TestAppendID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"https://example.com/marks",
			"https://example.com/marks?id=1",
		},
		{
			input:    "https://example.com/marks?week=1",
			expected: "https://example.com/marks?week=1&id=1",
		},
	}

	for _, test := range tests {
		result := appendID(test.input, 1)

		if result != test.expected {
			t.Errorf("appendID(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}
