package text

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"

	"github.com/codingconcepts/qapi/random"
)

var (
	variableRegex  = regexp.MustCompile(`{{\w+}}`)
	generatorRegex = regexp.MustCompile(`\[\[\w+\]\]`)
)

// AddVariables substitutes any placeholders for which there exists an equivalent
// variable and substitutes it; returning a string with all variables replaced.
func AddVariables(variables map[string]any, s string) string {
	placeholders := variableRegex.FindAllString(s, -1)
	for _, p := range placeholders {
		if v, ok := variables[strings.Trim(p, "{}")]; ok {
			var vs string

			// If the variable contains multiple elements, pick a random one, or just
			// convert the given scalar value to a string.
			if isSliceOrArray(v) {
				vs = fmt.Sprintf("%v", tryGetRandomElement(v))
			} else {
				vs = fmt.Sprintf("%v", v)
			}

			s = strings.ReplaceAll(s, p, vs)
		}
	}

	return s
}

func isSliceOrArray(v any) bool {
	if v == nil {
		return false
	}

	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

func tryGetRandomElement(v any) any {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)

	length := val.Len()
	if length == 0 {
		return nil
	}

	return val.Index(rand.Intn(length)).Interface()
}

// GenerateVariable generates a random value for a placeholder using a gofakeit
// generator.
func GenerateVariable(s string) string {
	placeholders := generatorRegex.FindAllString(s, -1)
	for _, p := range placeholders {
		if v, ok := random.Replacements[strings.Trim(p, "[]")]; ok {
			s = strings.ReplaceAll(s, p, fmt.Sprintf("%v", v()))
		}
	}

	return s
}
