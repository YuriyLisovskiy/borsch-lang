package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
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
		return 0, utilities.IncorrectUseOfFunctionError("compareTypes")
	}

	if left.EqualsTo(other) {
		return 0, nil
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func typeOperator_Call() common.Value {
	return NewFunctionInstance(
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
			return GetTypeOfInstance((*args)[1])
		},
		[]FunctionReturnType{
			{
				Type:       Any,
				IsNullable: false,
			},
		},
		true,
		nil,
		"", // TODO: add doc
	)
}

func newTypeClass() *Class {
	// TODO: add required operators and methods
	initAttributes := func(attrs *map[string]common.Value) {
		*attrs = map[string]common.Value{
			// TODO: add doc
			common.CallOperatorName: typeOperator_Call(),
			common.EqualsOp.Name(): MakeComparisonOperator(
				// TODO: add doc
				common.EqualsOp, TypeClass, "", compareTypes, func(res int) bool {
					return res == 0
				},
			),
			common.NotEqualsOp.Name(): MakeComparisonOperator(
				// TODO: add doc
				common.NotEqualsOp, TypeClass, "", compareTypes, func(res int) bool {
					return res != 0
				},
			),
		}
	}

	typeClass := &Class{
		Name:            common.TypeTypeName,
		IsFinal:         true,
		Bases:           []*Class{},
		Parent:          BuiltinPackage,
		AttrInitializer: initAttributes,
		GetEmptyInstance: func() (common.Value, error) {
			panic("unreachable")
		},
	}

	typeClass.Class = typeClass
	return typeClass
}
