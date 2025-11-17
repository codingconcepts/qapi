package text

import (
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestAddVariables(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		variables map[string]any
		expected  string
	}{
		{
			name:      "no placeholders or variables",
			input:     "test",
			variables: map[string]any{},
			expected:  "test",
		},
		{
			name:      "placeholders without variables",
			input:     "{{test}}",
			variables: map[string]any{},
			expected:  "{{test}}",
		},
		{
			name:  "variables without placeholders",
			input: "test",
			variables: map[string]any{
				"test": "hello",
			},
			expected: "test",
		},
		{
			name:  "single placeholder single variable",
			input: "{{test}}",
			variables: map[string]any{
				"test": "hello",
			},
			expected: "hello",
		},
		{
			name:  "multiple placeholder single variable",
			input: "{{test}} {{test}}",
			variables: map[string]any{
				"test": "hello",
			},
			expected: "hello hello",
		},
		{
			name:  "multiple placeholder multiple variable",
			input: "{{test}} {{test2}}",
			variables: map[string]any{
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

func TestGenerateVariable(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected func(string) bool
	}{
		{
			name:  "generates a random value",
			input: "[[weekday]]",
			expected: func(s string) bool {
				possibles := lo.Map([]int{0, 1, 2, 3, 4, 5, 6}, func(wd int, _ int) string {
					return time.Weekday(wd).String()
				})

				return lo.Contains(possibles, s)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := GenerateVariable(c.input)
			assert.True(t, c.expected(actual))
		})
	}
}
