package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func GetTypeOfInstance(object common.Value) (common.Value, error) {
	if instance, ok := object.(ObjectInstance); ok {
		return instance.GetClass(), nil
	}

	panic("unreachable")
}

func compareTypes(_ common.State, _ common.Operator, self, other common.Value) (int, error) {
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
	getTypeFunc := func(args ...common.Value) (common.Value, error) {
		return GetTypeOfInstance(args[0])
	}

	// TODO: add required operators and methods
	initAttributes := func(attrs *map[string]common.Value) {
		*attrs = map[string]common.Value{
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
				func(_ common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
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

	typeClass := NewClass(
		common.TypeTypeName,
		nil,
		BuiltinPackage,
		initAttributes,
		func() (common.Value, error) {
			panic("unreachable")
		},
	)
	typeClass.class = typeClass
	return typeClass
}
