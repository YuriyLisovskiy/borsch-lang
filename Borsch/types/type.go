package types

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
)

func GetTypeOfInstance(object common.Type) (common.Type, error) {
	if instance, ok := object.(Instance); ok {
		return instance.GetClass(), nil
	}

	return nil, errors.New("unknown object")
}

func newTypeClass() *Class {
	getTypeFunc := func(args ...common.Type) (common.Type, error) {
		return GetTypeOfInstance(args[0])
	}

	// TODO: add required operators and methods
	attributes := map[string]common.Type{
		// TODO: add doc
		ops.CallOperatorName: NewFunctionInstance(
			ops.CallOperatorName,
			[]FunctionArgument{
				{
					TypeHash:   TypeClassTypeHash,
					Name:       "я",
					IsVariadic: false,
					IsNullable: false,
				},
				{
					TypeHash:   AnyTypeHash,
					Name:       "обєкт",
					IsVariadic: false,
					IsNullable: false,
				},
			},
			func(_ interface{}, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				return getTypeFunc((*args)[1])
			},
			[]FunctionReturnType{
				{
					TypeHash:   AnyTypeHash,
					IsNullable: true,
				},
			},
			true,
			nil,
			"", // TODO: add doc
		),
	}
	return NewBuiltinClass(
		TypeClassTypeHash,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			panic("unreachable")
		},
	)
}
