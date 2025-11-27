package test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertSlicesEqualValues(t *testing.T, expected, actual any) bool {
	t.Helper()

	expectedVal := reflect.ValueOf(expected)
	actualVal := reflect.ValueOf(actual)

	if expectedVal.Kind() != reflect.Slice {
		return assert.Fail(t, "expected is not a slice")
	}
	if actualVal.Kind() != reflect.Slice {
		return assert.Fail(t, "actual is not a slice")
	}

	if !assert.Equal(t, expectedVal.Len(), actualVal.Len(), "slice lengths differ") {
		return false
	}

	for i := 0; i < expectedVal.Len(); i++ {
		exp := expectedVal.Index(i).Interface()
		act := actualVal.Index(i).Interface()

		expNum, expOk := toInt64(exp)
		actNum, actOk := toInt64(act)

		if expOk && actOk {
			if !assert.Equal(t, expNum, actNum, "elements at index %d differ", i) {
				return false
			}
		} else {
			if !assert.Equal(t, exp, act, "elements at index %d differ", i) {
				return false
			}
		}
	}

	return true
}

func toInt64(v any) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int8:
		return int64(val), true
	case int16:
		return int64(val), true
	case int32:
		return int64(val), true
	case int64:
		return val, true
	case uint:
		return int64(val), true
	case uint8:
		return int64(val), true
	case uint16:
		return int64(val), true
	case uint32:
		return int64(val), true
	case uint64:
		return int64(val), true
	default:
		return 0, false
	}
}
