package text

import (
	"regexp"
	"strings"
)

var variableRegex = regexp.MustCompile(`{{\w+}}`)

// AddVariables substitutes any placeholders for which there exists an equivalent
// variable and substitutes it; returning a string with all variables replaced.
func AddVariables(variables map[string]string, s string) string {
	placeholders := variableRegex.FindAllString(s, -1)
	for _, p := range placeholders {
		if v, ok := variables[strings.Trim(p, "{}")]; ok {
			s = strings.ReplaceAll(s, p, v)
		}
	}

	return s
}
