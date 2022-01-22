package types

import (
	"errors"
	"fmt"
	"math"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type BoolInstance struct {
	Object
	Value bool
}

func NewBoolInstance(value bool) BoolInstance {
	return BoolInstance{
		Value: value,
		Object: Object{
			typeName:    common.BoolTypeName,
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t BoolInstance) String(common.Context) string {
	if t.Value {
		return "істина"
	}

	return "хиба"
}

func (t BoolInstance) Representation(ctx common.Context) string {
	return t.String(ctx)
}

func (t BoolInstance) AsBool(common.Context) bool {
	return t.Value
}

func (t BoolInstance) SetAttribute(name string, _ common.Type) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	if t.Object.HasAttribute(name) || t.GetPrototype().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t BoolInstance) GetAttribute(name string) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetPrototype().GetAttribute(name)
}

func (BoolInstance) GetPrototype() *Class {
	return Bool
}

func compareBooleans(ctx common.Context, self common.Type, other common.Type) (int, error) {
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
		if left.Value == right.AsBool(ctx) {
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
	handler func(BoolInstance, common.Type) (common.Type, error),
) *FunctionInstance {
	return newBinaryMethod(
		name,
		Bool,
		Any,
		doc,
		func(ctx common.Context, left common.Type, right common.Type) (common.Type, error) {
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
	handler func(BoolInstance) (common.Type, error),
) *FunctionInstance {
	return newUnaryMethod(
		name, Bool, Any, doc, func(ctx common.Context, left common.Type) (common.Type, error) {
			if leftInstance, ok := left.(BoolInstance); ok {
				return handler(leftInstance)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newBoolClass() *Class {
	attributes := mergeAttributes(
		map[string]common.Type{
			// TODO: add doc
			ops.ConstructorName: newBuiltinConstructor(Bool, ToBool, ""),
			ops.PowOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.PowOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.UnaryPlus.Caption(), "", func(self BoolInstance) (common.Type, error) {
					return NewIntegerInstance(boolToInt64(self.Value)), nil
				},
			),
			ops.UnaryMinus.Caption(): newBoolUnaryOperator(
				// TODO: add doc
				ops.UnaryMinus.Caption(), "", func(self BoolInstance) (common.Type, error) {
					return NewIntegerInstance(-boolToInt64(self.Value)), nil
				},
			),
			ops.UnaryBitwiseNotOp.Caption(): newBoolUnaryOperator(
				// TODO: add doc
				ops.UnaryBitwiseNotOp.Caption(), "", func(self BoolInstance) (common.Type, error) {
					return NewIntegerInstance(^boolToInt64(self.Value)), nil
				},
			),
			ops.MulOp.Caption(): newBoolBinaryOperator(
				// TODO: add doc
				ops.MulOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.DivOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.ModuloOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.AddOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.SubOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseLeftShiftOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseRightShiftOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseAndOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseXorOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseOrOp.Caption(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
		makeLogicalOperators(Bool),
		makeComparisonOperators(Bool, compareBooleans),
		makeCommonOperators(Bool),
	)
	return NewBuiltinClass(
		common.BoolTypeName,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewBoolInstance(false), nil
		},
	)
}
