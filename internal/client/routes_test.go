package client

import (
	"testing"
)

func TestAppendID(t *testing.T) {
	id := 1

	tests := []struct {
		id       int
		input    string
		expected string
	}{
		{
			id,
			"https://example.com/marks",
			"https://example.com/marks?id=1",
		},
		{
			id,
			"https://example.com/marks?week=1",
			"https://example.com/marks?week=1&id=1",
		},
	}

	for _, test := range tests {
		result := appendID(test.input, test.id)

		if result != test.expected {
			t.Errorf("appendID(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}
