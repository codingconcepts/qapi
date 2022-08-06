package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddVariables(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		variables map[string]string
		expected  string
	}{
		{
			name:      "no placeholders or variables",
			input:     "test",
			variables: map[string]string{},
			expected:  "test",
		},
		{
			name:      "placeholders without variables",
			input:     "{{test}}",
			variables: map[string]string{},
			expected:  "{{test}}",
		},
		{
			name:  "variables without placeholders",
			input: "test",
			variables: map[string]string{
				"test": "hello",
			},
			expected: "test",
		},
		{
			name:  "single placeholder single variable",
			input: "{{test}}",
			variables: map[string]string{
				"test": "hello",
			},
			expected: "hello",
		},
		{
			name:  "multiple placeholder single variable",
			input: "{{test}} {{test}}",
			variables: map[string]string{
				"test": "hello",
			},
			expected: "hello hello",
		},
		{
			name:  "multiple placeholder multiple variable",
			input: "{{test}} {{test2}}",
			variables: map[string]string{
				"test":  "hello",
				"test2": "goodbye",
			},
			expected: "hello goodbye",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := AddVariables(c.variables, c.input)
			assert.Equal(t, c.expected, actual)
		})
	}
}
