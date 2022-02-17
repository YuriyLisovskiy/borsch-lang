package types

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
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

func MergeAttributes(a map[string]common.Value, b ...map[string]common.Value) map[string]common.Value {
	for _, m := range b {
		for key, val := range m {
			a[key] = val
		}
	}

	return a
}
