package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
)

func GetTypeOfInstance(object common.Type) (common.Type, error) {
	if instance, ok := object.(ObjectInstance); ok {
		proto := instance.GetPrototype()
		return proto, nil
	}

	panic("unreachable")
}

func newTypeClass() *Class {
	getTypeFunc := func(args ...common.Type) (common.Type, error) {
		return GetTypeOfInstance(args[0])
	}

	// TODO: add required operators and methods
	initAttributes := func() map[string]common.Type {
		return map[string]common.Type{
			// TODO: add doc
			ops.CallOperatorName: NewFunctionInstance(
				ops.CallOperatorName,
				[]FunctionParameter{
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
				func(_ common.State, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
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
	}

	class := &Class{
		// TODO: add doc
		Object: *newClassObject(common.TypeTypeName, BuiltinPackage, initAttributes, ""),
		GetEmptyInstance: func() (common.Type, error) {
			panic("unreachable")
		},
	}

	class.prototype = class
	return class
}
