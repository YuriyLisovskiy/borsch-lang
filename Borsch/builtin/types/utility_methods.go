package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func MakeBinaryMethod(
	name string,
	selfType *Class,
	returnType *Class,
	doc string,
	handler func(common.State, common.Object, common.Object) (common.Object, error),
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
				Type:       AnyClass,
				Name:       "інший",
				IsVariadic: false,
				IsNullable: true,
			},
		},
		func(state common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
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

func MakeUnaryMethod(
	name string,
	selfType *Class,
	returnType *Class,
	doc string,
	handler func(common.State, common.Object) (common.Object, error),
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
		func(state common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
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

func MakeComparisonOperator(
	operator common.Operator,
	itemType *Class,
	doc string,
	comparator func(common.State, common.Operator, common.Object, common.Object) (int, error),
	checker func(res int) bool,
) *FunctionInstance {
	return MakeBinaryMethod(
		operator.Name(),
		itemType,
		BoolClass,
		doc,
		func(state common.State, self common.Object, other common.Object) (common.Object, error) {
			res, err := comparator(state, operator, self, other)
			if err != nil {
				return nil, err
			}

			return NewBoolInstance(checker(res)), nil
		},
	)
}

func MakeComparisonOperators(
	itemType *Class,
	comparator func(common.State, common.Operator, common.Object, common.Object) (int, error),
) map[string]common.Object {
	return map[string]common.Object{
		common.EqualsOp.Name(): MakeComparisonOperator(
			// TODO: add doc
			common.EqualsOp, itemType, "", comparator, func(res int) bool {
				return res == 0
			},
		),
		common.NotEqualsOp.Name(): MakeComparisonOperator(
			// TODO: add doc
			common.NotEqualsOp, itemType, "", comparator, func(res int) bool {
				return res != 0
			},
		),
		common.GreaterOp.Name(): MakeComparisonOperator(
			// TODO: add doc
			common.GreaterOp, itemType, "", comparator, func(res int) bool {
				return res == 1
			},
		),
		common.GreaterOrEqualsOp.Name(): MakeComparisonOperator(
			// TODO: add doc
			common.GreaterOrEqualsOp, itemType, "", comparator, func(res int) bool {
				return res == 0 || res == 1
			},
		),
		common.LessOp.Name(): MakeComparisonOperator(
			// TODO: add doc
			common.LessOp, itemType, "", comparator, func(res int) bool {
				return res == -1
			},
		),
		common.LessOrEqualsOp.Name(): MakeComparisonOperator(
			// TODO: add doc
			common.LessOrEqualsOp, itemType, "", comparator, func(res int) bool {
				return res == 0 || res == -1
			},
		),
	}
}

func MakeLogicalOperators(itemType *Class) map[string]common.Object {
	return map[string]common.Object{
		common.NotOp.Name(): MakeUnaryMethod(
			// TODO: add doc
			common.NotOp.Name(),
			itemType,
			BoolClass,
			"",
			func(state common.State, self common.Object) (common.Object, error) {
				selfBool, err := self.AsBool(state)
				if err != nil {
					return nil, err
				}

				return NewBoolInstance(!selfBool), nil
			},
		),
		common.AndOp.Name(): MakeBinaryMethod(
			// TODO: add doc
			common.AndOp.Name(),
			itemType,
			BoolClass,
			"",
			func(state common.State, self common.Object, other common.Object) (common.Object, error) {
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
		common.OrOp.Name(): MakeBinaryMethod(
			// TODO: add doc
			common.OrOp.Name(),
			itemType,
			BoolClass,
			"",
			func(state common.State, self common.Object, other common.Object) (common.Object, error) {
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

func MakeCommonOperators(itemType *Class) map[string]common.Object {
	return map[string]common.Object{
		// TODO: add doc
		common.BoolOperatorName: MakeUnaryMethod(
			common.BoolOperatorName, itemType, BoolClass, "",
			func(state common.State, self common.Object) (common.Object, error) {
				boolVal, err := self.AsBool(state)
				if err != nil {
					return nil, err
				}

				return NewBoolInstance(boolVal), nil
			},
		),
	}
}

func MakeUnaryOperators(
	selfClass, returnClass *Class,
	handler func(common.State, common.Operator, common.Object) (common.Object, error),
) map[string]common.Object {
	return map[string]common.Object{
		common.UnaryPlus.Name(): MakeUnaryMethod(
			common.UnaryPlus.Name(),
			selfClass,
			returnClass,
			"",
			func(state common.State, value common.Object) (common.Object, error) {
				return handler(state, common.UnaryPlus, value)
			},
		),
		common.UnaryMinus.Name(): MakeUnaryMethod(
			common.UnaryMinus.Name(),
			selfClass,
			returnClass,
			"",
			func(state common.State, value common.Object) (common.Object, error) {
				return handler(state, common.UnaryMinus, value)
			},
		),
		common.UnaryBitwiseNotOp.Name(): MakeUnaryMethod(
			common.UnaryBitwiseNotOp.Name(),
			selfClass,
			returnClass,
			"",
			func(state common.State, value common.Object) (common.Object, error) {
				return handler(state, common.UnaryBitwiseNotOp, value)
			},
		),
	}
}

func makeVariadicConstructor(
	itemType *Class,
	converter func(common.State, ...common.Object) (common.Object, error),
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
				Type:       AnyClass,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(state common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
			self, err := converter(state, (*args)[1:]...)
			if err != nil {
				return nil, err
			}

			(*args)[0] = self
			return NewNilInstance(), nil
		},
		[]FunctionReturnType{
			{
				Type:       NilClass,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}

func makeLengthOperator(
	itemType *Class,
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
		func(state common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
			sequence := (*args)[0]
			switch self := sequence.(type) {
			case common.SequentialType:
				return NewIntegerInstance(self.Length(state)), nil
			}

			return nil, errors.New(fmt.Sprint("invalid type in length operator: ", sequence.GetTypeName()))
		},
		[]FunctionReturnType{
			{
				Type:       IntClass,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}

func makeDefaultConstructor(cls *Class, doc string) *FunctionInstance {
	if cls == nil {
		panic("makeDefaultConstructor: cls is nil")
	}

	return NewFunctionInstance(
		common.ConstructorName,
		[]FunctionParameter{
			{
				Type:       cls,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
			instance, err := cls.GetEmptyInstance()
			if err != nil {
				return nil, err
			}

			(*args)[0] = instance
			return NewNilInstance(), nil
		},
		[]FunctionReturnType{
			{
				Type:       NilClass,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}
