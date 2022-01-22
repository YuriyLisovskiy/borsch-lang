package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
)

func GetTypeOfInstance(object common.Type) (common.Type, error) {
	if instance, ok := object.(ObjectInstance); ok {
		return instance.GetPrototype(), nil
	}

	panic("unreachable")
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
					Type:       TypeClass,
					Name:       "я",
					IsVariadic: false,
					IsNullable: false,
				},
				{
					Type:       Any,
					Name:       "обєкт",
					IsVariadic: false,
					IsNullable: false,
				},
			},
			func(_ common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				return getTypeFunc((*args)[1])
			},
			[]FunctionReturnType{
				{
					Type:       Any,
					IsNullable: true,
				},
			},
			true,
			nil,
			"", // TODO: add doc
		),
	}
	return NewBuiltinClass(
		common.TypeTypeName,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			panic("unreachable")
		},
	)
}
