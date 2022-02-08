package types

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type IntegerInstance struct {
	BuiltinInstance
	Value int64
}

func NewIntegerInstance(value int64) IntegerInstance {
	return IntegerInstance{
		BuiltinInstance: BuiltinInstance{
			ClassInstance{
				class:      Integer,
				attributes: map[string]common.Value{},
				address:    "",
			},
		},
		Value: value,
	}
}

func (t IntegerInstance) String(common.State) (string, error) {
	return fmt.Sprintf("%d", t.Value), nil
}

func (t IntegerInstance) Representation(state common.State) (string, error) {
	return t.String(state)
}

func (t IntegerInstance) AsBool(common.State) (bool, error) {
	return t.Value != 0, nil
}

func toInteger(_ common.State, args ...common.Value) (common.Value, error) {
	if len(args) == 0 {
		return NewIntegerInstance(0), nil
	}

	if len(args) != 1 {
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"'цілий()' приймає лише один аргумент (отримано %d)", len(args),
			),
		)
	}

	switch vt := args[0].(type) {
	case RealInstance:
		return NewIntegerInstance(int64(vt.Value)), nil
	case IntegerInstance:
		return vt, nil
	case StringInstance:
		intVal, err := strconv.ParseInt(vt.Value, 10, 64)
		if err != nil {
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"некоректний літерал для функції 'цілий()' з основою 10: '%s'", vt.Value,
				),
			)
		}

		return NewIntegerInstance(intVal), nil
	case BoolInstance:
		if vt.Value {
			return NewIntegerInstance(1), nil
		}

		return NewIntegerInstance(0), nil
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"'%s' неможливо інтерпретувати як ціле число", args[0].GetTypeName(),
			),
		)
	}
}

func evalUnaryOperatorWithIntegers(_ common.State, operator common.Operator, value common.Value) (common.Value, error) {
	if self, ok := value.(IntegerInstance); ok {
		switch operator {
		case common.UnaryPlus:
			return self, nil
		case common.UnaryMinus:
			return NewIntegerInstance(-self.Value), nil
		case common.UnaryBitwiseNotOp:
			return NewIntegerInstance(^self.Value), nil
		default:
			return nil, util.InternalOperatorError(operator)
		}
	}

	return nil, util.BadOperandForUnaryOperatorError(operator)
}

func compareIntegers(_ common.State, op common.Operator, self common.Value, other common.Value) (int, error) {
	left, ok := self.(IntegerInstance)
	if !ok {
		return 0, util.IncorrectUseOfFunctionError("compareIntegers")
	}

	switch right := other.(type) {
	case NilInstance:
	case BoolInstance:
		rightVal := boolToInt64(right.Value)
		if left.Value == rightVal {
			return 0, nil
		}

		if left.Value < rightVal {
			return -1, nil
		}

		return 1, nil
	case IntegerInstance:
		if left.Value == right.Value {
			return 0, nil
		}

		if left.Value < right.Value {
			return -1, nil
		}

		return 1, nil
	case RealInstance:
		leftVal := float64(left.Value)
		if leftVal == right.Value {
			return 0, nil
		}

		if leftVal < right.Value {
			return -1, nil
		}

		return 1, nil
	default:
		return 0, util.OperatorNotSupportedError(op, left, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func intOperator(
	operator common.Operator,
	handler func(common.State, IntegerInstance, common.Value) (common.Value, error),
) common.Value {
	return NewFunctionInstance(
		operator.Name(),
		[]FunctionParameter{
			{
				Type:       Integer,
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
			left, ok := (*args)[0].(IntegerInstance)
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

func intOperator_Pow(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case RealInstance:
		return NewRealInstance(math.Pow(float64(left.Value), other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(int64(math.Pow(float64(left.Value), float64(other.Value)))), nil
	case BoolInstance:
		return NewIntegerInstance(int64(math.Pow(float64(left.Value), boolToFloat64(other.Value)))), nil
	default:
		return nil, nil
	}
}

func intOperator_Mul(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(left.Value * boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(left.Value * other.Value), nil
	case RealInstance:
		return NewRealInstance(float64(left.Value) * other.Value), nil
	case StringInstance:
		count := int(left.Value)
		if count <= 0 {
			return NewStringInstance(""), nil
		}

		return NewStringInstance(strings.Repeat(other.Value, count)), nil
	case ListInstance:
		count := int(left.Value)
		list := NewListInstance()
		if count > 0 {
			for c := 0; c < count; c++ {
				list.Values = append(list.Values, other.Values...)
			}
		}

		return list, nil
	default:
		return nil, nil
	}
}

func intOperator_Div(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		if other.Value {
			return NewRealInstance(float64(left.Value)), nil
		}
	case IntegerInstance:
		if other.Value != 0 {
			return NewRealInstance(float64(left.Value) / float64(other.Value)), nil
		}
	case RealInstance:
		if other.Value != 0.0 {
			return NewRealInstance(float64(left.Value) / other.Value), nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення на нуль")
}

func intOperator_Modulo(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		if other.Value {
			return NewIntegerInstance(left.Value % boolToInt64(other.Value)), nil
		}
	case IntegerInstance:
		if other.Value != 0 {
			return NewIntegerInstance(left.Value % other.Value), nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення за модулем на нуль")
}

func intOperator_Add(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(left.Value + boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(left.Value + other.Value), nil
	case RealInstance:
		return NewRealInstance(float64(left.Value) + other.Value), nil
	default:
		return nil, nil
	}
}

func intOperator_Sub(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(left.Value - boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(left.Value - other.Value), nil
	case RealInstance:
		return NewRealInstance(float64(left.Value) - other.Value), nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseLeftShift(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(left.Value << boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(left.Value << other.Value), nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseRightShift(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(left.Value >> boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(left.Value >> other.Value), nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseAnd(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(left.Value & boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(left.Value & other.Value), nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseXor(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(left.Value ^ boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(left.Value ^ other.Value), nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseOr(_ common.State, left IntegerInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(left.Value | boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(left.Value | other.Value), nil
	default:
		return nil, nil
	}
}

func newIntegerClass() *Class {
	return &Class{
		Name:    common.IntegerTypeName,
		IsFinal: true,
		Bases:   []*Class{},
		Parent:  BuiltinPackage,
		AttrInitializer: func(attrs *map[string]common.Value) {
			*attrs = MergeAttributes(
				map[string]common.Value{
					// TODO: add doc
					common.ConstructorName: makeVariadicConstructor(Integer, toInteger, ""),

					common.PowOp.Name():    intOperator(common.PowOp, intOperator_Pow),
					common.MulOp.Name():    intOperator(common.MulOp, intOperator_Mul),
					common.DivOp.Name():    intOperator(common.DivOp, intOperator_Div),
					common.ModuloOp.Name(): intOperator(common.ModuloOp, intOperator_Modulo),
					common.AddOp.Name():    intOperator(common.AddOp, intOperator_Add),
					common.SubOp.Name():    intOperator(common.SubOp, intOperator_Sub),
					common.BitwiseLeftShiftOp.Name(): intOperator(
						common.BitwiseLeftShiftOp,
						intOperator_BitwiseLeftShift,
					),
					common.BitwiseRightShiftOp.Name(): intOperator(
						common.BitwiseRightShiftOp,
						intOperator_BitwiseRightShift,
					),
					common.BitwiseAndOp.Name(): intOperator(common.BitwiseAndOp, intOperator_BitwiseAnd),
					common.BitwiseXorOp.Name(): intOperator(common.BitwiseXorOp, intOperator_BitwiseXor),
					common.BitwiseOrOp.Name():  intOperator(common.BitwiseOrOp, intOperator_BitwiseOr),
				},
				MakeUnaryOperators(Integer, Integer, evalUnaryOperatorWithIntegers),
				MakeLogicalOperators(Integer),
				MakeComparisonOperators(Integer, compareIntegers),
				MakeCommonOperators(Integer),
			)
		},
		GetEmptyInstance: func() (common.Value, error) {
			return NewIntegerInstance(0), nil
		},
	}
}
