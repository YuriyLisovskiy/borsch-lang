package types

import (
	"errors"
	"fmt"
	"math"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type Bool bool

func NewBoolInstance(value bool) Bool {
	return Bool(value)
}

func (t Bool) GetClass() *Class {
	return BoolType
}

func (t Bool) GetTypeName() string {
	return t.GetClass().GetName()
}

func (t Bool) GetOperator(name string) (common.Value, error) {
	if attr, err := t.GetClass().getAttribute(name); err == nil {
		return attr, nil
	}

	return nil, utilities.OperatorNotFoundError(t.GetTypeName(), name)
}

func (t Bool) GetAttribute(name string) (common.Value, error) {
	if attr, err := t.GetClass().getAttribute(name); err == nil {
		return attr, nil
	}

	return nil, utilities.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t Bool) SetAttribute(name string, _ common.Value) error {
	return utilities.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t Bool) HasAttribute(name string) bool {
	return t.GetClass().HasAttribute(name)
}

func (t Bool) String(common.State) (string, error) {
	if t {
		return "істина", nil
	}

	return "хиба", nil
}

func (t Bool) Representation(state common.State) (string, error) {
	return t.String(state)
}

func (t Bool) AsBool(common.State) (bool, error) {
	return bool(t), nil
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
	left, ok := self.(Bool)
	if !ok {
		return 0, utilities.IncorrectUseOfFunctionError("compareBooleans")
	}

	switch right := other.(type) {
	case NilInstance:
	case Bool:
		if left == right {
			return 0, nil
		}
	case Int, RealInstance:
		rightBool, err := right.AsBool(state)
		if err != nil {
			return 0, err
		}

		if bool(left) == rightBool {
			return 0, nil
		}
	default:
		return 0, utilities.OperatorNotSupportedError(op, left, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func evalUnaryOperatorWithBooleans(_ common.State, operator common.Operator, value common.Value) (common.Value, error) {
	if self, ok := value.(Bool); ok {
		switch operator {
		case common.UnaryPlus:
			return Int(boolToInt64(self)), nil
		case common.UnaryMinus:
			return Int(-boolToInt64(self)), nil
		case common.UnaryBitwiseNotOp:
			return Int(^boolToInt64(self)), nil
		default:
			return nil, utilities.InternalOperatorError(operator)
		}
	}

	return nil, utilities.BadOperandForUnaryOperatorError(operator)
}

func boolOperator(
	operator common.Operator,
	handler func(common.State, Bool, common.Value) (common.Value, error),
) common.Value {
	return NewFunctionInstance(
		operator.Name(),
		[]FunctionParameter{
			{
				Type:       BoolType,
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
			left, ok := (*args)[0].(Bool)
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

func boolOperator_Pow(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case RealInstance:
		return NewRealInstance(math.Pow(boolToFloat64(left), other.Value)), nil
	case Int:
		return Int(math.Pow(boolToFloat64(left), float64(other))), nil
	case Bool:
		return Int(math.Pow(boolToFloat64(left), boolToFloat64(other))), nil
	default:
		return nil, nil
	}
}

func boolOperator_Mul(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return Int(boolToInt64(left) * boolToInt64(other)), nil
	case Int:
		return Int(boolToInt64(left)) * other, nil
	case RealInstance:
		return NewRealInstance(boolToFloat64(left) * other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_Div(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		if other {
			return NewRealInstance(boolToFloat64(left)), nil
		}
	case Int:
		if other != 0 {
			return NewRealInstance(boolToFloat64(left) / float64(other)), nil
		}
	case RealInstance:
		if other.Value != 0.0 {
			return NewRealInstance(boolToFloat64(left) / other.Value), nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення на нуль")
}

func boolOperator_Modulo(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		if other {
			return Int(boolToInt64(left) % boolToInt64(other)), nil
		}
	case Int:
		if other != 0 {
			return Int(boolToInt64(left)) % other, nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення за модулем на нуль")
}

func boolOperator_Add(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return Int(boolToInt64(left) + boolToInt64(other)), nil
	case Int:
		return Int(boolToInt64(left)) + other, nil
	case RealInstance:
		return NewRealInstance(boolToFloat64(left) + other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_Sub(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return Int(boolToInt64(left) - boolToInt64(other)), nil
	case Int:
		return Int(boolToInt64(left)) - other, nil
	case RealInstance:
		return NewRealInstance(boolToFloat64(left) - other.Value), nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseLeftShift(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return Int(boolToInt64(left) << boolToInt64(other)), nil
	case Int:
		return Int(boolToInt64(left)) << other, nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseRightShift(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return Int(boolToInt64(left) >> boolToInt64(other)), nil
	case Int:
		return Int(boolToInt64(left)) >> other, nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseAnd(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return Int(boolToInt64(left) & boolToInt64(other)), nil
	case Int:
		return Int(boolToInt64(left)) & other, nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseXor(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return Int(boolToInt64(left) ^ boolToInt64(other)), nil
	case Int:
		return Int(boolToInt64(left)) ^ other, nil
	default:
		return nil, nil
	}
}

func boolOperator_BitwiseOr(_ common.State, left Bool, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Bool:
		return Int(boolToInt64(left) | boolToInt64(other)), nil
	case Int:
		return Int(boolToInt64(left)) | other, nil
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
					common.ConstructorName: makeVariadicConstructor(BoolType, toBool, ""),

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
				MakeUnaryOperators(BoolType, Integer, evalUnaryOperatorWithBooleans),
				MakeLogicalOperators(BoolType),
				MakeComparisonOperators(BoolType, compareBooleans),
				MakeCommonOperators(BoolType),
			)
		},
		GetEmptyInstance: func() (common.Value, error) {
			return NewBoolInstance(false), nil
		},
	}
}
