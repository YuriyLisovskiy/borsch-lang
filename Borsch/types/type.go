package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func GetTypeOfInstance(object common.Type) (common.Type, error) {
	if instance, ok := object.(ObjectInstance); ok {
		return instance.GetPrototype(), nil
	}

	panic("unreachable")
}

func compareTypes(_ common.State, self, other common.Type) (int, error) {
	left, ok := self.(*Class)
	if !ok {
		return 0, util.IncorrectUseOfFunctionError("compareTypes")
	}

	if left.EqualsTo(other) {
		return 0, nil
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newTypeClass() *Class {
	getTypeFunc := func(args ...common.Type) (common.Type, error) {
		return GetTypeOfInstance(args[0])
	}

	// TODO: add required operators and methods
	initAttributes := func() map[string]common.Type {
		return map[string]common.Type{
			// TODO: add doc
			common.CallOperatorName: NewFunctionInstance(
				common.CallOperatorName,
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
			common.EqualsOp.Name(): NewComparisonOperator(
				// TODO: add doc
				common.EqualsOp, TypeClass, "", compareTypes, func(res int) bool {
					return res == 0
				},
			),
			common.NotEqualsOp.Name(): NewComparisonOperator(
				// TODO: add doc
				common.NotEqualsOp, TypeClass, "", compareTypes, func(res int) bool {
					return res != 0
				},
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
