package types

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type BoolInstance struct {
	Object
	Value bool
}

func NewBoolInstanceFromString(value string) (BoolInstance, error) {
	switch value {
	case "істина":
		value = "t"
	case "хиба":
		value = "f"
	}

	boolean, err := strconv.ParseBool(value)
	if err != nil {
		return BoolInstance{}, util.RuntimeError(err.Error())
	}

	return NewBoolInstance(boolean), nil
}

func NewBoolInstance(value bool) BoolInstance {
	return BoolInstance{
		Value: value,
		Object: Object{
			typeName:    GetTypeName(BoolTypeHash),
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t BoolInstance) String() string {
	if t.Value {
		return "істина"
	}

	return "хиба"
}

func (t BoolInstance) Representation() string {
	return t.String()
}

func (t BoolInstance) GetTypeHash() uint64 {
	return t.GetClass().GetTypeHash()
}

func (t BoolInstance) AsBool() bool {
	return t.Value
}

func (t BoolInstance) SetAttribute(name string, _ Type) (Type, error) {
	if t.Object.HasAttribute(name) || t.GetClass().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t BoolInstance) GetAttribute(name string) (Type, error) {
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetClass().GetAttribute(name)
}

func (BoolInstance) GetClass() *Class {
	return Bool
}

func compareBooleans(self Type, other Type) (int, error) {
	left, ok := self.(BoolInstance)
	if !ok {
		return 0, util.IncorrectUseOfFunctionError("compareBooleans")
	}

	switch right := other.(type) {
	case NilInstance:
	case BoolInstance:
		if left.Value == right.Value {
			return 0, nil
		}
	case IntegerInstance, RealInstance:
		if left.Value == right.AsBool() {
			return 0, nil
		}
	default:
		return 0, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор '%s' до значень типів '%s' та '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newBoolBinaryOperator(
	name string,
	doc string,
	handler func(BoolInstance, Type) (Type, error),
) *FunctionInstance {
	return newBinaryMethod(
		name, BoolTypeHash, AnyTypeHash, doc, func(left Type, right Type) (Type, error) {
			if leftInstance, ok := left.(BoolInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newBoolUnaryOperator(
	name string,
	doc string,
	handler func(BoolInstance) (Type, error),
) *FunctionInstance {
	return newUnaryMethod(
		name, BoolTypeHash, AnyTypeHash, doc, func(left Type) (Type, error) {
			if leftInstance, ok := left.(BoolInstance); ok {
				return handler(leftInstance)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newBoolClass() *Class {
	attributes := mergeAttributes(
		map[string]Type{
			// TODO: add doc
			ops.ConstructorName: newBuiltinConstructor(BoolTypeHash, ToBool, ""),
			ops.PowOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.PowOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case RealInstance:
						return NewRealInstance(math.Pow(boolToFloat64(self.Value), o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(int64(math.Pow(boolToFloat64(self.Value), float64(o.Value)))), nil
					case BoolInstance:
						return NewIntegerInstance(
							int64(
								math.Pow(
									boolToFloat64(self.Value),
									boolToFloat64(o.Value),
								),
							),
						), nil
					default:
						return nil, nil
					}
				},
			),
			ops.UnaryPlus.Caption(): newBoolUnaryOperator(
				// TODO: add doc
				ops.UnaryPlus.Caption(), "", func(self BoolInstance) (Type, error) {
					return NewIntegerInstance(boolToInt64(self.Value)), nil
				},
			),
			ops.UnaryMinus.Caption(): newBoolUnaryOperator(
				// TODO: add doc
				ops.UnaryMinus.Caption(), "", func(self BoolInstance) (Type, error) {
					return NewIntegerInstance(-boolToInt64(self.Value)), nil
				},
			),
			ops.UnaryBitwiseNotOp.Caption(): newBoolUnaryOperator(
				// TODO: add doc
				ops.UnaryBitwiseNotOp.Caption(), "", func(self BoolInstance) (Type, error) {
					return NewIntegerInstance(^boolToInt64(self.Value)), nil
				},
			),
			ops.MulOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.MulOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(boolToInt64(self.Value) * boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(boolToInt64(self.Value) * o.Value), nil
					case RealInstance:
						return NewRealInstance(boolToFloat64(self.Value) * o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.DivOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.DivOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						if o.Value {
							return NewRealInstance(boolToFloat64(self.Value)), nil
						}
					case IntegerInstance:
						if o.Value != 0 {
							return NewRealInstance(boolToFloat64(self.Value) / float64(o.Value)), nil
						}
					case RealInstance:
						if o.Value != 0.0 {
							return NewRealInstance(boolToFloat64(self.Value) / o.Value), nil
						}
					default:
						return nil, nil
					}

					return nil, errors.New("ділення на нуль")
				},
			),
			ops.ModuloOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.ModuloOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						if o.Value {
							return NewIntegerInstance(boolToInt64(self.Value) % boolToInt64(o.Value)), nil
						}
					case IntegerInstance:
						if o.Value != 0 {
							return NewIntegerInstance(boolToInt64(self.Value) % o.Value), nil
						}
					default:
						return nil, nil
					}

					return nil, errors.New("ділення за модулем на нуль")
				},
			),
			ops.AddOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.AddOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(boolToInt64(self.Value) + boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(boolToInt64(self.Value) + o.Value), nil
					case RealInstance:
						return NewRealInstance(boolToFloat64(self.Value) + o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.SubOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.SubOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(boolToInt64(self.Value) - boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(boolToInt64(self.Value) - o.Value), nil
					case RealInstance:
						return NewRealInstance(boolToFloat64(self.Value) - o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseLeftShiftOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.BitwiseLeftShiftOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(boolToInt64(self.Value) << boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(boolToInt64(self.Value) << o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseRightShiftOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.BitwiseRightShiftOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(boolToInt64(self.Value) >> boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(boolToInt64(self.Value) >> o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseAndOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.BitwiseAndOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(boolToInt64(self.Value) & boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(boolToInt64(self.Value) & o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseXorOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.BitwiseXorOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(boolToInt64(self.Value) ^ boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(boolToInt64(self.Value) ^ o.Value), nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseOrOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.BitwiseOrOp.Caption(), "", func(self BoolInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolInstance:
						return NewIntegerInstance(boolToInt64(self.Value) | boolToInt64(o.Value)), nil
					case IntegerInstance:
						return NewIntegerInstance(boolToInt64(self.Value) | o.Value), nil
					default:
						return nil, nil
					}
				},
			),
		},
		makeLogicalOperators(BoolTypeHash),
		makeComparisonOperators(BoolTypeHash, compareBooleans),
		makeCommonOperators(BoolTypeHash),
	)
	return NewBuiltinClass(
		BoolTypeHash,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (Type, error) {
			return NewBoolInstance(false), nil
		},
	)
}
