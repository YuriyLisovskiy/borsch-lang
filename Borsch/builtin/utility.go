package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func Length(sequence types.Type) (types.Type, error) {
	switch arg := sequence.(type) {
	case types.SequentialType:
		return types.NewIntegerInstance(arg.Length()), nil
	case types.DictionaryInstance:
		return types.NewIntegerInstance(arg.Length()), nil
	}

	return nil, util.RuntimeError(fmt.Sprintf(
		"об'єкт типу '%s' не має довжини", sequence.GetTypeName(),
	))
}

func AppendToList(list types.ListInstance, values ...types.Type) (types.Type, error) {
	for _, value := range values {
		list.Values = append(list.Values, value)
	}

	return list, nil
}

func RemoveFromDictionary(dict types.DictionaryInstance, key types.Type) (types.Type, error) {
	err := dict.RemoveElement(key)
	if err != nil {
		return nil, util.RuntimeError(err.Error())
	}

	return dict, nil
}
