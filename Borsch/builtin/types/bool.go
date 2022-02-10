package types

import (
	"errors"
	"fmt"
	"math"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type BoolInstance struct {
	BuiltinInstance
	Value bool
}

func NewBoolInstance(value bool) BoolInstance {
	return BoolInstance{
		BuiltinInstance: BuiltinInstance{
			ClassInstance: *NewClassInstance(Bool, nil),
		},
		Value: value,
	}
}

func (t BoolInstance) String(common.State) (string, error) {
	if t.Value {
		return "істина", nil
	}

	return "хиба", nil
}

func (t BoolInstance) Representation(state common.State) (string, error) {
	return t.String(state)
}

func (t BoolInstance) AsBool(common.State) (bool, error) {
	return t.Value, nil
}

func toBool(state common.State, args ...common.Value) (common.Value, error) {
	if len(args) == 0 {
		return NewBoolInstance(false), nil
	}

	if len(args) != 1 {
		return nil, errors.New(
			fmt.Sprintf(
				"функція 'логічний()' приймає лише один аргумент (отримано %d)", len(args),
			),
		)
	}

	boolValue, err := args[0].AsBool(state)
	if err != nil {
		return nil, err
	}

	return NewBoolInstance(boolValue), err
}

func compareBooleans(state common.State, op common.Operator, self common.Value, other common.Value) (int, error) {
	left, ok := self.(BoolInstance)
	if !ok {
		return 0, utilities.IncorrectUseOfFunctionError("compareBooleans")
	}

	switch right := other.(type) {
	case NilInstance:
	case BoolInstance:
		if left.Value == right.Value {
			return 0, nil
		}
	case IntegerInstance, RealInstance:
		rightBool, err := right.AsBool(state)
		if err != nil {
			return 0, err
		}

		if left.Value == rightBool {
			return 0, nil
		}
	default:
		return 0, utilities.OperatorNotSupportedError(op, left, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func evalUnaryOperatorWithBooleans(_ common.State, operator common.Operator, value common.Value) (common.Value, error) {
	if self, ok := value.(BoolInstance); ok {
		switch operator {
		case common.UnaryPlus:
			return NewIntegerInstance(boolToInt64(self.Value)), nil
		case common.UnaryMinus:
			return NewIntegerInstance(-boolToInt64(self.Value)), nil
		case common.UnaryBitwiseNotOp:
			return NewIntegerInstance(^boolToInt64(self.Value)), nil
		default:
			return nil, utilities.InternalOperatorError(operator)
		}
	}

	return nil, utilities.BadOperandForUnaryOperatorError(operator)
}

func boolOperator(
	operator common.Operator,
	handler func(common.State, BoolInstance, common.Value) (common.Value, error),
) common.Value {
	return NewFunctionInstance(
		operator.Name(),
		[]FunctionParameter{
			{
				Type:       Bool,
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
			left, ok := (*args)[0].(BoolInstance)
			if !ok {
				return nil, utilities.InvalidUseOfOperator(operator, left, (*args)[1])
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

func boolOperator_Pow(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case RealInstance:
		return NewRealInstance(math.Pow(boolToFloat64(left.Value), other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(int64(math.Pow(boolToFloat64(left.Value), float64(other.Value)))), nil
	case BoolInstance:
		return NewIntegerInstance(int64(math.Pow(boolToFloat64(left.Value), boolToFloat64(other.Value)))), nil
	default:
		return nil, nil
	}
}

func boolOperator_Mul(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(boolToInt64(left.Value) * boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(boolToInt64(left.Value) * other.Value), nil
	case RealInstance:
		return NewRealInstance(boolToFloat64(left.Value) * other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_Div(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		if other.Value {
			return NewRealInstance(boolToFloat64(left.Value)), nil
		}
	case IntegerInstance:
		if other.Value != 0 {
			return NewRealInstance(boolToFloat64(left.Value) / float64(other.Value)), nil
		}
	case RealInstance:
		if other.Value != 0.0 {
			return NewRealInstance(boolToFloat64(left.Value) / other.Value), nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення на нуль")
}

func boolOperator_Modulo(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		if other.Value {
			return NewIntegerInstance(boolToInt64(left.Value) % boolToInt64(other.Value)), nil
		}
	case IntegerInstance:
		if other.Value != 0 {
			return NewIntegerInstance(boolToInt64(left.Value) % other.Value), nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення за модулем на нуль")
}

func boolOperator_Add(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(boolToInt64(left.Value) + boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(boolToInt64(left.Value) + other.Value), nil
	case RealInstance:
		return NewRealInstance(boolToFloat64(left.Value) + other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_Sub(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(boolToInt64(left.Value) - boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(boolToInt64(left.Value) - other.Value), nil
	case RealInstance:
		return NewRealInstance(boolToFloat64(left.Value) - other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseLeftShift(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(boolToInt64(left.Value) << boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(boolToInt64(left.Value) << other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseRightShift(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(boolToInt64(left.Value) >> boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(boolToInt64(left.Value) >> other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseAnd(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(boolToInt64(left.Value) & boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(boolToInt64(left.Value) & other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseXor(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(boolToInt64(left.Value) ^ boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(boolToInt64(left.Value) ^ other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseOr(_ common.State, left BoolInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case BoolInstance:
		return NewIntegerInstance(boolToInt64(left.Value) | boolToInt64(other.Value)), nil
	case IntegerInstance:
		return NewIntegerInstance(boolToInt64(left.Value) | other.Value), nil
	default:
		return nil, nil
	}
}

func newBoolClass() *Class {
	return &Class{
		Name:    common.BoolTypeName,
		IsFinal: true,
		Bases:   []*Class{},
		Parent:  BuiltinPackage,
		AttrInitializer: func(attrs *map[string]common.Value) {
			*attrs = MergeAttributes(
				map[string]common.Value{
					// TODO: add doc
					common.ConstructorName: makeVariadicConstructor(Bool, toBool, ""),

					common.PowOp.Name():    boolOperator(common.PowOp, boolOperator_Pow),
					common.MulOp.Name():    boolOperator(common.MulOp, boolOperator_Mul),
					common.DivOp.Name():    boolOperator(common.DivOp, boolOperator_Div),
					common.ModuloOp.Name(): boolOperator(common.ModuloOp, boolOperator_Modulo),
					common.AddOp.Name():    boolOperator(common.AddOp, boolOperator_Add),
					common.SubOp.Name():    boolOperator(common.SubOp, boolOperator_Sub),
					common.BitwiseLeftShiftOp.Name(): boolOperator(
						common.BitwiseLeftShiftOp,
						boolOperator_BitwiseLeftShift,
					),
					common.BitwiseRightShiftOp.Name(): boolOperator(
						common.BitwiseRightShiftOp,
						boolOperator_BitwiseRightShift,
					),
					common.BitwiseAndOp.Name(): boolOperator(common.BitwiseAndOp, boolOperator_BitwiseAnd),
					common.BitwiseXorOp.Name(): boolOperator(common.BitwiseXorOp, boolOperator_BitwiseXor),
					common.BitwiseOrOp.Name():  boolOperator(common.BitwiseOrOp, boolOperator_BitwiseOr),
				},
				MakeUnaryOperators(Bool, Integer, evalUnaryOperatorWithBooleans),
				MakeLogicalOperators(Bool),
				MakeComparisonOperators(Bool, compareBooleans),
				MakeCommonOperators(Bool),
			)
		},
		GetEmptyInstance: func() (common.Value, error) {
			return NewBoolInstance(false), nil
		},
	}
}
