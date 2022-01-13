package types

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
)

func getType(object Type) (Type, error) {
	if instance, ok := object.(Instance); ok {
		return instance.GetClass(), nil
	}

	return nil, errors.New("unknown object")
}

func newTypeClass() *Class {
	getTypeFunc := func(args ...Type) (Type, error) {
		return getType(args[0])
	}

	// TODO: add required operators and methods
	attributes := map[string]Type{
		// TODO: add doc
		ops.ConstructorName: newBuiltinConstructor(StringTypeHash, getTypeFunc, ""),
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
			func(args *[]Type, _ *map[string]Type) (Type, error) {
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
		func() (Type, error) {
			return nil, nil
		},
	)
}
