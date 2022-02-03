package types

import (
	"errors"
	"fmt"

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

func getAttributes(attributes map[string]common.Type) (DictionaryInstance, error) {
	dict := NewDictionaryInstance()
	for key, val := range attributes {
		err := dict.SetElement(NewStringInstance(key), val)
		if err != nil {
			return DictionaryInstance{}, err
		}
	}

	return dict, nil
}

func getLength(state common.State, sequence common.Type) (int64, error) {
	switch self := sequence.(type) {
	case common.SequentialType:
		return self.Length(state), nil
	}

	return 0, errors.New(fmt.Sprint("invalid type in length operator: ", sequence.GetTypeName()))
}

func MergeAttributes(a map[string]common.Type, b ...map[string]common.Type) map[string]common.Type {
	for _, m := range b {
		for key, val := range m {
			a[key] = val
		}
	}

	return a
}
