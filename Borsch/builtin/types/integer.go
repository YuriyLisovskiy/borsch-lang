package types

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type IntegerInstance struct {
	Object
	Value int64
}

func NewIntegerInstanceFromString(value string) (IntegerInstance, error) {
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return IntegerInstance{}, util.RuntimeError(err.Error())
	}

	return NewIntegerInstance(number), nil
}

func NewIntegerInstance(value int64) IntegerInstance {
	return IntegerInstance{
		Object: Object{
			typeName:    GetTypeName(IntegerTypeHash),
			Attributes:  nil,
			callHandler: nil,
		},
		Value: value,
	}
}

func (t IntegerInstance) String() string {
	return fmt.Sprintf("%d", t.Value)
}

func (t IntegerInstance) Representation() string {
	return t.String()
}

func (t IntegerInstance) GetTypeHash() uint64 {
	return t.GetClass().GetTypeHash()
}

func (t IntegerInstance) AsBool() bool {
	return t.Value != 0
}

func (t IntegerInstance) SetAttribute(name string, _ Type) (Type, error) {
	if t.Object.HasAttribute(name) || t.GetClass().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t IntegerInstance) GetAttribute(name string) (Type, error) {
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetClass().GetAttribute(name)
}

func (IntegerInstance) GetClass() *Class {
	return Integer
}

func compareIntegers(self Type, other Type) (int, error) {
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
		return 0, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newIntegerBinaryOperator(
	name string,
	doc string,
	handler func(IntegerInstance, Type) (Type, error),
) *FunctionInstance {
	return newBinaryMethod(
		name, IntegerTypeHash, AnyTypeHash, doc, func(left Type, right Type) (Type, error) {
			if leftInstance, ok := left.(IntegerInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newIntegerUnaryOperator(
	name string,
	doc string,
	handler func(IntegerInstance) (Type, error),
) *FunctionInstance {
	return newUnaryMethod(
		name, IntegerTypeHash, AnyTypeHash, doc, func(left Type) (Type, error) {
			if leftInstance, ok := left.(IntegerInstance); ok {
				return handler(leftInstance)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newIntegerClass() *Class {
	attributes := mergeAttributes(
		map[string]Type{
			// TODO: add doc
			ops.ConstructorName: newBuiltinConstructor(IntegerTypeHash, ToInteger, ""),
			ops.PowOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.PowOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case RealInstance:
						return NewRealInstance(math.Pow(float64(self.Value), o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(int64(math.Pow(float64(self.Value), float64(o.Value)))), nil
					case BoolInstance:
						return NewIntegerInstance(int64(math.Pow(float64(self.Value), boolToFloat64(o.Value)))), nil
					default:
						return nil, nil
					}
				},
			),
			ops.UnaryPlus.Caption(): newIntegerUnaryOperator(
				// TODO: add doc
				ops.UnaryPlus.Caption(), "", func(self IntegerInstance) (Type, error) {
					return self, nil
				},
			),
			ops.UnaryMinus.Caption(): newIntegerUnaryOperator(
				// TODO: add doc
				ops.UnaryMinus.Caption(), "", func(self IntegerInstance) (Type, error) {
					return NewIntegerInstance(-self.Value), nil
				},
			),
			ops.UnaryBitwiseNotOp.Caption(): newIntegerUnaryOperator(
				// TODO: add doc
				ops.UnaryBitwiseNotOp.Caption(), "", func(self IntegerInstance) (Type, error) {
					return NewIntegerInstance(^self.Value), nil
				},
			),
			ops.MulOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.MulOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(self.Value * boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(self.Value * o.Value), nil
					case RealInstance:
						return NewRealInstance(float64(self.Value) * o.Value), nil
					case StringInstance:
						count := int(self.Value)
						if count <= 0 {
							return NewStringInstance(""), nil
						}

						return NewStringInstance(strings.Repeat(o.Value, count)), nil
					case ListInstance:
						count := int(self.Value)
						list := NewListInstance()
						if count > 0 {
							for c := 0; c < count; c++ {
								list.Values = append(list.Values, o.Values...)
							}
						}

						return list, nil
					default:
						return nil, nil
					}
				},
			),
			ops.DivOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.DivOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						if o.Value {
							return NewRealInstance(float64(self.Value)), nil
						}
					case IntegerInstance:
						if o.Value != 0 {
							return NewRealInstance(float64(self.Value) / float64(o.Value)), nil
						}
					case RealInstance:
						if o.Value != 0.0 {
							return NewRealInstance(float64(self.Value) / o.Value), nil
						}
					default:
						return nil, nil
					}

					return nil, errors.New("ділення на нуль")
				},
			),
			ops.ModuloOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.ModuloOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						if o.Value {
							return NewIntegerInstance(self.Value % boolToInt64(o.Value)), nil
						}
					case IntegerInstance:
						if o.Value != 0 {
							return NewIntegerInstance(self.Value % o.Value), nil
						}
					default:
						return nil, nil
					}

					return nil, errors.New("ділення за модулем на нуль")
				},
			),
			ops.AddOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.AddOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(self.Value + boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(self.Value + o.Value), nil
					case RealInstance:
						return NewRealInstance(float64(self.Value) + o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.SubOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.SubOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(self.Value - boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(self.Value - o.Value), nil
					case RealInstance:
						return NewRealInstance(float64(self.Value) - o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseLeftShiftOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.BitwiseLeftShiftOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(self.Value << boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(self.Value << o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseRightShiftOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.BitwiseRightShiftOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(self.Value >> boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(self.Value >> o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseAndOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.BitwiseAndOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(self.Value & boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(self.Value & o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseXorOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.BitwiseXorOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(self.Value ^ boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(self.Value ^ o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseOrOp.Caption(): newIntegerBinaryOperator(
				// TODO: add doc
				ops.BitwiseOrOp.Caption(), "", func(self IntegerInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(self.Value | boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(self.Value | o.Value), nil
					default:
						return nil, nil
					}
				},
			),
		},
		makeLogicalOperators(IntegerTypeHash),
		makeComparisonOperators(IntegerTypeHash, compareIntegers),
	)
	return NewBuiltinClass(
		IntegerTypeHash,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (Type, error) {
			return NewIntegerInstance(0), nil
		},
	)
}
