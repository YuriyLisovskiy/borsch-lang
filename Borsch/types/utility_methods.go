package types

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/common"

func newBinaryMethod(
	name string,
	selfType *Class,
	returnType *Class,
	doc string,
	handler func(common.State, common.Type, common.Type) (common.Type, error),
) *FunctionInstance {
	return NewFunctionInstance(
		name,
		[]FunctionParameter{
			{
				Type:       selfType,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       Any,
				Name:       "інший",
				IsVariadic: false,
				IsNullable: true,
			},
		},
		func(state common.State, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
			return handler(state, (*args)[0], (*args)[1])
		},
		[]FunctionReturnType{
			{
				Type:       returnType,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}

func newUnaryMethod(
	name string,
	selfType *Class,
	returnType *Class,
	doc string,
	handler func(common.State, common.Type) (common.Type, error),
) *FunctionInstance {
	return NewFunctionInstance(
		name,
		[]FunctionParameter{
			{
				Type:       selfType,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
			return handler(state, (*args)[0])
		},
		[]FunctionReturnType{
			{
				Type:       returnType,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}

func NewComparisonOperator(
	operator common.Operator,
	itemType *Class,
	doc string,
	comparator func(common.State, common.Type, common.Type) (int, error),
	checker func(res int) bool,
) *FunctionInstance {
	return newBinaryMethod(
		operator.Name(),
		itemType,
		Bool,
		doc,
		func(state common.State, self common.Type, other common.Type) (common.Type, error) {
			res, err := comparator(state, self, other)
			if err != nil {
				return nil, err
			}

			return NewBoolInstance(checker(res)), nil
		},
	)
}

func MakeComparisonOperators(
	itemType *Class,
	comparator func(common.State, common.Type, common.Type) (int, error),
) map[string]common.Type {
	return map[string]common.Type{
		common.EqualsOp.Name(): NewComparisonOperator(
			// TODO: add doc
			common.EqualsOp, itemType, "", comparator, func(res int) bool {
				return res == 0
			},
		),
		common.NotEqualsOp.Name(): NewComparisonOperator(
			// TODO: add doc
			common.NotEqualsOp, itemType, "", comparator, func(res int) bool {
				return res != 0
			},
		),
		common.GreaterOp.Name(): NewComparisonOperator(
			// TODO: add doc
			common.GreaterOp, itemType, "", comparator, func(res int) bool {
				return res == 1
			},
		),
		common.GreaterOrEqualsOp.Name(): NewComparisonOperator(
			// TODO: add doc
			common.GreaterOrEqualsOp, itemType, "", comparator, func(res int) bool {
				return res == 0 || res == 1
			},
		),
		common.LessOp.Name(): NewComparisonOperator(
			// TODO: add doc
			common.LessOp, itemType, "", comparator, func(res int) bool {
				return res == -1
			},
		),
		common.LessOrEqualsOp.Name(): NewComparisonOperator(
			// TODO: add doc
			common.LessOrEqualsOp, itemType, "", comparator, func(res int) bool {
				return res == 0 || res == -1
			},
		),
	}
}

func MakeLogicalOperators(itemType *Class) map[string]common.Type {
	return map[string]common.Type{
		common.NotOp.Name(): newUnaryMethod(
			// TODO: add doc
			common.NotOp.Name(),
			itemType,
			Bool,
			"",
			func(state common.State, self common.Type) (common.Type, error) {
				selfBool, err := self.AsBool(state)
				if err != nil {
					return nil, err
				}

				return NewBoolInstance(!selfBool), nil
			},
		),
		common.AndOp.Name(): newBinaryMethod(
			// TODO: add doc
			common.AndOp.Name(),
			itemType,
			Bool,
			"",
			func(state common.State, self common.Type, other common.Type) (common.Type, error) {
				selfBool, err := self.AsBool(state)
				if err != nil {
					return nil, err
				}

				otherBool, err := other.AsBool(state)
				if err != nil {
					return nil, err
				}

				return NewBoolInstance(selfBool && otherBool), nil
			},
		),
		common.OrOp.Name(): newBinaryMethod(
			// TODO: add doc
			common.OrOp.Name(),
			itemType,
			Bool,
			"",
			func(state common.State, self common.Type, other common.Type) (common.Type, error) {
				selfBool, err := self.AsBool(state)
				if err != nil {
					return nil, err
				}

				otherBool, err := other.AsBool(state)
				if err != nil {
					return nil, err
				}

				return NewBoolInstance(selfBool || otherBool), nil
			},
		),
	}
}

func MakeCommonOperators(itemType *Class) map[string]common.Type {
	return map[string]common.Type{
		// TODO: add doc
		common.BoolOperatorName: newUnaryMethod(
			common.BoolOperatorName, itemType, Bool, "",
			func(state common.State, self common.Type) (common.Type, error) {
				boolVal, err := self.AsBool(state)
				if err != nil {
					return nil, err
				}

				return NewBoolInstance(boolVal), nil
			},
		),
	}
}

func newBuiltinConstructor(
	itemType *Class,
	handler func(common.State, ...common.Type) (common.Type, error),
	doc string,
) *FunctionInstance {
	return NewFunctionInstance(
		common.ConstructorName,
		[]FunctionParameter{
			{
				Type:       itemType,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       Any,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(state common.State, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
			self, err := handler(state, (*args)[1:]...)
			if err != nil {
				return nil, err
			}

			(*args)[0] = self
			return NewNilInstance(), nil
		},
		[]FunctionReturnType{
			{
				Type:       Nil,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}

func newLengthOperator(
	itemType *Class,
	handler func(common.State, common.Type) (int64, error),
	doc string,
) *FunctionInstance {
	return NewFunctionInstance(
		common.LengthOperatorName,
		[]FunctionParameter{
			{
				Type:       itemType,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
			length, err := handler(state, (*args)[0])
			if err != nil {
				return nil, err
			}

			return NewIntegerInstance(length), nil
		},
		[]FunctionReturnType{
			{
				Type:       Integer,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}
