package text

import (
	"regexp"
	"strings"
)

var variableRegex = regexp.MustCompile(`{{\w+}}`)

func AddVariables(variables map[string]string, s string) string {
	placeholders := variableRegex.FindAllString(s, -1)
	for _, p := range placeholders {
		if v, ok := variables[strings.Trim(p, "{}")]; ok {
			s = strings.ReplaceAll(s, p, v)
		}
	}

	return s
}
