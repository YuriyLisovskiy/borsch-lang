package types

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type Int int64

func (t Int) GetClass() *Class {
	return Integer
}

func (t Int) GetTypeName() string {
	return t.GetClass().GetName()
}

func (t Int) GetOperator(name string) (common.Value, error) {
	if attr, err := t.GetClass().getAttribute(name); err == nil {
		return attr, nil
	}

	return nil, utilities.OperatorNotFoundError(t.GetTypeName(), name)
}

func (t Int) GetAttribute(name string) (common.Value, error) {
	if attr, err := t.GetClass().getAttribute(name); err == nil {
		return attr, nil
	}

	return nil, utilities.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t Int) SetAttribute(name string, _ common.Value) error {
	return utilities.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t Int) HasAttribute(name string) bool {
	return t.GetClass().HasAttribute(name)
}

func (t Int) String(common.State) (string, error) {
	return fmt.Sprintf("%d", t), nil
}

func (t Int) Representation(state common.State) (string, error) {
	return t.String(state)
}

func (t Int) AsBool(common.State) (bool, error) {
	return t != 0, nil
}

func toInteger(_ common.State, args ...common.Value) (common.Value, error) {
	if len(args) == 0 {
		return Int(0), nil
	}

	if len(args) != 1 {
		return nil, errors.New(
			fmt.Sprintf(
				"'цілий()' приймає лише один аргумент (отримано %d)", len(args),
			),
		)
	}

	switch vt := args[0].(type) {
	case RealInstance:
		return Int(vt.Value), nil
	case Int:
		return vt, nil
	case StringInstance:
		intVal, err := strconv.ParseInt(vt.Value, 10, 64)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf(
					"некоректний літерал для функції 'цілий()' з основою 10: '%s'", vt.Value,
				),
			)
		}

		return Int(intVal), nil
	case Bool:
		if vt {
			return Int(1), nil
		}

		return Int(0), nil
	default:
		return nil, errors.New(
			fmt.Sprintf(
				"'%s' неможливо інтерпретувати як ціле число", args[0].GetTypeName(),
			),
		)
	}
}

func evalUnaryOperatorWithIntegers(_ common.State, operator common.Operator, value common.Value) (common.Value, error) {
	if self, ok := value.(Int); ok {
		switch operator {
		case common.UnaryPlus:
			return self, nil
		case common.UnaryMinus:
			return -self, nil
		case common.UnaryBitwiseNotOp:
			return ^self, nil
		default:
			return nil, utilities.InternalOperatorError(operator)
		}
	}

	return nil, utilities.BadOperandForUnaryOperatorError(operator)
}

func compareIntegers(_ common.State, op common.Operator, self common.Value, other common.Value) (int, error) {
	left, ok := self.(Int)
	if !ok {
		return 0, utilities.IncorrectUseOfFunctionError("compareIntegers")
	}

	switch right := other.(type) {
	case NilInstance:
	case Bool:
		rightVal := Int(boolToInt64(right))
		if left == rightVal {
			return 0, nil
		}

		if left < rightVal {
			return -1, nil
		}

		return 1, nil
	case Int:
		if left == right {
			return 0, nil
		}

		if left < right {
			return -1, nil
		}

		return 1, nil
	case RealInstance:
		leftVal := float64(left)
		if leftVal == right.Value {
			return 0, nil
		}

		if leftVal < right.Value {
			return -1, nil
		}

		return 1, nil
	default:
		return 0, utilities.OperatorNotSupportedError(op, left, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func intOperator(
	operator common.Operator,
	handler func(common.State, Int, common.Value) (common.Value, error),
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
			left, ok := (*args)[0].(Int)
			if !ok {
				return nil, utilities.InvalidUseOfOperator(operator, left, (*args)[1])
			}

			right := (*args)[1]
			result, err := handler(state, left, right)
			if err != nil {
				return nil, err
			}

			if result == nil {
				return nil, utilities.OperatorNotSupportedError(operator, left, right)
			}

			return result, nil
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

func intOperator_Pow(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case RealInstance:
		return NewRealInstance(math.Pow(float64(left), other.Value)), nil
	case Int:
		return Int(math.Pow(float64(left), float64(other))), nil
	case Bool:
		return Int(math.Pow(float64(left), boolToFloat64(other))), nil
	default:
		return nil, nil
	}
}

func intOperator_Mul(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return left * Int(boolToInt64(other)), nil
	case Int:
		return left * other, nil
	case RealInstance:
		return NewRealInstance(float64(left) * other.Value), nil
	case StringInstance:
		count := int(left)
		if count <= 0 {
			return NewStringInstance(""), nil
		}

		return NewStringInstance(strings.Repeat(other.Value, count)), nil
	case ListInstance:
		count := int(left)
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

func intOperator_Div(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		if other {
			return NewRealInstance(float64(left)), nil
		}
	case Int:
		if other != 0 {
			return NewRealInstance(float64(left) / float64(other)), nil
		}
	case RealInstance:
		if other.Value != 0.0 {
			return NewRealInstance(float64(left) / other.Value), nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення на нуль")
}

func intOperator_Modulo(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		if other {
			return left % Int(boolToInt64(other)), nil
		}
	case Int:
		if other != 0 {
			return left % other, nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення за модулем на нуль")
}

func intOperator_Add(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return left + Int(boolToInt64(other)), nil
	case Int:
		return left + other, nil
	case RealInstance:
		return NewRealInstance(float64(left) + other.Value), nil
	default:
		return nil, nil
	}
}

func intOperator_Sub(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return left - Int(boolToInt64(other)), nil
	case Int:
		return left - other, nil
	case RealInstance:
		return NewRealInstance(float64(left) - other.Value), nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseLeftShift(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return left << Int(boolToInt64(other)), nil
	case Int:
		return left << other, nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseRightShift(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return left >> Int(boolToInt64(other)), nil
	case Int:
		return left >> other, nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseAnd(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return left & Int(boolToInt64(other)), nil
	case Int:
		return left & other, nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseXor(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return left ^ Int(boolToInt64(other)), nil
	case Int:
		return left ^ other, nil
	default:
		return nil, nil
	}
}

func intOperator_BitwiseOr(_ common.State, left Int, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return left | Int(boolToInt64(other)), nil
	case Int:
		return left | other, nil
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
			return Int(0), nil
		},
	}
}
