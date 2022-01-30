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
	initAttributes := func() map[string]common.Type {
		return mergeAttributes(
			map[string]common.Type{
				// TODO: add doc
				ops.ConstructorName: newBuiltinConstructor(Bool, ToBool, ""),
				ops.PowOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.PowOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.UnaryPlus.Name(): newBoolUnaryOperator(
					// TODO: add doc
					ops.UnaryPlus.Name(), "", func(self BoolInstance) (common.Type, error) {
						return NewIntegerInstance(boolToInt64(self.Value)), nil
					},
				),
				ops.UnaryMinus.Name(): newBoolUnaryOperator(
					// TODO: add doc
					ops.UnaryMinus.Name(), "", func(self BoolInstance) (common.Type, error) {
						return NewIntegerInstance(-boolToInt64(self.Value)), nil
					},
				),
				ops.UnaryBitwiseNotOp.Name(): newBoolUnaryOperator(
					// TODO: add doc
					ops.UnaryBitwiseNotOp.Name(), "", func(self BoolInstance) (common.Type, error) {
						return NewIntegerInstance(^boolToInt64(self.Value)), nil
					},
				),
				ops.MulOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.MulOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.DivOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.DivOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.ModuloOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.ModuloOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.AddOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.AddOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.SubOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.SubOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseLeftShiftOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.BitwiseLeftShiftOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseRightShiftOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.BitwiseRightShiftOp.Name(),
					"",
					func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseAndOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.BitwiseAndOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseXorOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.BitwiseXorOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
				ops.BitwiseOrOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					ops.BitwiseOrOp.Name(), "", func(self BoolInstance, other common.Type) (common.Type, error) {
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
	}

	return NewBuiltinClass(
		common.BoolTypeName,
		BuiltinPackage,
		initAttributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewBoolInstance(false), nil
		},
	)
}
