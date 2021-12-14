package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func Length(sequence types.ValueType) (types.ValueType, error) {
	switch arg := sequence.(type) {
	case types.SequentialType:
		return types.IntegerType{Value: arg.Length()}, nil
	case types.DictionaryType:
		return types.IntegerType{Value: arg.Length()}, nil
	}

	return nil, util.RuntimeError(fmt.Sprintf(
		"об'єкт типу '%s' не має довжини", sequence.TypeName(),
	))
}

func AppendToList(list types.ListType, values ...types.ValueType) (types.ValueType, error) {
	for _, value := range values {
		list.Values = append(list.Values, value)
	}

	return list, nil
}

func RemoveFromDictionary(dict types.DictionaryType, key types.ValueType) (types.ValueType, error) {
	err := dict.RemoveElement(key)
	if err != nil {
		return nil, util.RuntimeError(err.Error())
	}

	return dict, nil
}
