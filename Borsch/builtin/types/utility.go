package types

import (
	"errors"
	"go/types"
)

func getIndex(index, length int64) (int64, error) {
	if index >= 0 && index < length {
		return index, nil
	} else if index < 0 && index >= -length {
		return length + index, nil
	}

	return 0, errors.New("індекс за межами послідовності")
}

func normalizeBound(bound, length int64) int64 {
	if bound < 0 {
		return length + bound
	}

	return bound
}

func boolToInt64(v bool) int64 {
	if v {
		return 1
	}

	return 0
}

func boolToFloat64(v bool) float64 {
	if v {
		return 1.0
	}

	return 0.0
}

func MergeAttributes(a map[string]types.Object, b ...map[string]types.Object) map[string]types.Object {
	for _, m := range b {
		for key, val := range m {
			a[key] = val
		}
	}

	return a
}
