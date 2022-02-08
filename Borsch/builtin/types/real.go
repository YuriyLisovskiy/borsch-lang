package types

import (
	"errors"
	"math"
	"strconv"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type RealInstance struct {
	BuiltinInstance
	Value float64
}

func NewRealInstance(value float64) RealInstance {
	return RealInstance{
		BuiltinInstance: BuiltinInstance{
			ClassInstance: ClassInstance{
				class:      Real,
				attributes: map[string]common.Value{},
				address:    "",
			},
		},
		Value: value,
	}
}

func (t RealInstance) String(common.State) (string, error) {
	return strconv.FormatFloat(t.Value, 'f', -1, 64), nil
}

func (t RealInstance) Representation(state common.State) (string, error) {
	return t.String(state)
}

func (t RealInstance) AsBool(common.State) (bool, error) {
	return t.Value != 0.0, nil
}

func compareReals(_ common.State, op common.Operator, self, other common.Value) (int, error) {
	left, ok := self.(RealInstance)
	if !ok {
		return 0, util.IncorrectUseOfFunctionError("compareReals")
	}

	switch right := other.(type) {
	case NilInstance:
	case BoolInstance:
		rightVal := boolToFloat64(right.Value)
		if left.Value == rightVal {
			return 0, nil
		}

		if left.Value < rightVal {
			return -1, nil
		}

		return 1, nil
	case IntegerInstance:
		rightVal := float64(right.Value)
		if left.Value == rightVal {
			return 0, nil
		}

		if left.Value < rightVal {
			return -1, nil
		}

		return 1, nil
	case RealInstance:
		if left.Value == right.Value {
			return 0, nil
		}

		if left.Value < right.Value {
			return -1, nil
		}

		return 1, nil
	default:
		return 0, util.OperatorNotSupportedError(op, left, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func realBinaryOperator(
	operator common.Operator,
	handler func(common.State, RealInstance, common.Value) (common.Value, error),
) common.Value {
	return NewFunctionInstance(
		operator.Name(),
		[]FunctionParameter{
			{
				Type:       Real,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       Any,
				Name:       "інший",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			left, ok := (*args)[0].(RealInstance)
			if !ok {
				return nil, util.InvalidUseOfOperator(operator, left, (*args)[1])
			}

			return handler(state, left, (*args)[1])
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

func realUnaryOperator(
	operator common.Operator,
	handler func(common.State, RealInstance) (common.Value, error),
) common.Value {
	return NewFunctionInstance(
		operator.Name(),
		[]FunctionParameter{
			{
				Type:       Real,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			left, ok := (*args)[0].(RealInstance)
			if !ok {
				return nil, util.InvalidUseOfOperator(operator, left, (*args)[1])
			}

			return handler(state, left)
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

func realOperator_Pow(_ common.State, left RealInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case RealInstance:
		return NewRealInstance(math.Pow(left.Value, other.Value)), nil
	case IntegerInstance:
		return NewRealInstance(math.Pow(left.Value, float64(other.Value))), nil
	case BoolInstance:
		return NewRealInstance(math.Pow(left.Value, boolToFloat64(other.Value))), nil
	default:
		return nil, nil
	}
}

func realOperator_UnaryPlus(_ common.State, self RealInstance) (common.Value, error) {
	return self, nil
}

func realOperator_UnaryMinus(_ common.State, self RealInstance) (common.Value, error) {
	return NewRealInstance(-self.Value), nil
}

func realOperator_Mul(_ common.State, left RealInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewRealInstance(left.Value * boolToFloat64(other.Value)), nil
	case IntegerInstance:
		return NewRealInstance(left.Value * float64(other.Value)), nil
	case RealInstance:
		return NewRealInstance(left.Value * other.Value), nil
	default:
		return nil, nil
	}
}

func realOperator_Div(_ common.State, left RealInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		if other.Value {
			return NewRealInstance(left.Value), nil
		}
	case IntegerInstance:
		if other.Value != 0 {
			return NewRealInstance(left.Value / float64(other.Value)), nil
		}
	case RealInstance:
		if other.Value != 0.0 {
			return NewRealInstance(left.Value / other.Value), nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення на нуль")
}

func realOperator_Add(_ common.State, left RealInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewRealInstance(left.Value + boolToFloat64(other.Value)), nil
	case IntegerInstance:
		return NewRealInstance(left.Value + float64(other.Value)), nil
	case RealInstance:
		return NewRealInstance(left.Value + other.Value), nil
	default:
		return nil, nil
	}
}

func realOperator_Sub(_ common.State, left RealInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewRealInstance(left.Value - boolToFloat64(other.Value)), nil
	case IntegerInstance:
		return NewRealInstance(left.Value - float64(other.Value)), nil
	case RealInstance:
		return NewRealInstance(left.Value - other.Value), nil
	default:
		return nil, nil
	}
}

func newRealClass() *Class {
	initAttributes := func(attrs *map[string]common.Value) {
		*attrs = MergeAttributes(
			map[string]common.Value{
				// TODO: add doc
				common.ConstructorName:   makeVariadicConstructor(Real, ToReal, ""),
				common.PowOp.Name():      realBinaryOperator(common.PowOp, realOperator_Pow),
				common.UnaryPlus.Name():  realUnaryOperator(common.UnaryPlus, realOperator_UnaryPlus),
				common.UnaryMinus.Name(): realUnaryOperator(common.UnaryMinus, realOperator_UnaryMinus),
				common.MulOp.Name():      realBinaryOperator(common.MulOp, realOperator_Mul),
				common.DivOp.Name():      realBinaryOperator(common.DivOp, realOperator_Div),
				common.AddOp.Name():      realBinaryOperator(common.AddOp, realOperator_Add),
				common.SubOp.Name():      realBinaryOperator(common.SubOp, realOperator_Sub),
			},
			MakeLogicalOperators(Real),
			MakeComparisonOperators(Real, compareReals),
			MakeCommonOperators(Real),
		)
	}

	return &Class{
		Name:            common.RealTypeName,
		IsFinal:         true,
		Bases:           []*Class{},
		Parent:          BuiltinPackage,
		AttrInitializer: initAttributes,
		GetEmptyInstance: func() (common.Value, error) {
			return NewRealInstance(0), nil
		},
	}
}
