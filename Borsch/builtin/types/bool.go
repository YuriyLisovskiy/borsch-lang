package types

import (
	"errors"
	"math"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
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

func compareBooleans(state common.State, op common.Operator, self common.Value, other common.Value) (int, error) {
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
		rightBool, err := right.AsBool(state)
		if err != nil {
			return 0, err
		}

		if left.Value == rightBool {
			return 0, nil
		}
	default:
		return 0, util.OperatorNotSupportedError(op, left.GetTypeName(), right.GetTypeName())
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newBoolBinaryOperator(
	name string,
	doc string,
	handler func(BoolInstance, common.Value) (common.Value, error),
) *FunctionInstance {
	return newBinaryMethod(
		name,
		Bool,
		Any,
		doc,
		func(_ common.State, left common.Value, right common.Value) (common.Value, error) {
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
	handler func(BoolInstance) (common.Value, error),
) *FunctionInstance {
	return newUnaryMethod(
		name, Bool, Any, doc, func(_ common.State, left common.Value) (common.Value, error) {
			if leftInstance, ok := left.(BoolInstance); ok {
				return handler(leftInstance)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newBoolClass() *Class {
	initAttributes := func(attrs *map[string]common.Value) {
		*attrs = MergeAttributes(
			map[string]common.Value{
				// TODO: add doc
				common.ConstructorName: newBuiltinConstructor(Bool, ToBool, ""),
				common.PowOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.PowOp.Name(), "", func(self BoolInstance, other common.Value) (common.Value, error) {
						switch o := other.(type) {
						case RealInstance:
							return NewRealInstance(math.Pow(boolToFloat64(self.Value), o.Value)), nil
						case IntegerInstance:
							return NewIntegerInstance(
								int64(
									math.Pow(
										boolToFloat64(self.Value),
										float64(o.Value),
									),
								),
							), nil
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
				common.UnaryPlus.Name(): newBoolUnaryOperator(
					// TODO: add doc
					common.UnaryPlus.Name(), "", func(self BoolInstance) (common.Value, error) {
						return NewIntegerInstance(boolToInt64(self.Value)), nil
					},
				),
				common.UnaryMinus.Name(): newBoolUnaryOperator(
					// TODO: add doc
					common.UnaryMinus.Name(), "", func(self BoolInstance) (common.Value, error) {
						return NewIntegerInstance(-boolToInt64(self.Value)), nil
					},
				),
				common.UnaryBitwiseNotOp.Name(): newBoolUnaryOperator(
					// TODO: add doc
					common.UnaryBitwiseNotOp.Name(), "", func(self BoolInstance) (common.Value, error) {
						return NewIntegerInstance(^boolToInt64(self.Value)), nil
					},
				),
				common.MulOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.MulOp.Name(), "", func(self BoolInstance, other common.Value) (common.Value, error) {
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
				common.DivOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.DivOp.Name(), "", func(self BoolInstance, other common.Value) (common.Value, error) {
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
				common.ModuloOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.ModuloOp.Name(), "", func(self BoolInstance, other common.Value) (common.Value, error) {
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
				common.AddOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.AddOp.Name(), "", func(self BoolInstance, other common.Value) (common.Value, error) {
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
				common.SubOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.SubOp.Name(), "", func(self BoolInstance, other common.Value) (common.Value, error) {
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
				common.BitwiseLeftShiftOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.BitwiseLeftShiftOp.Name(),
					"",
					func(self BoolInstance, other common.Value) (common.Value, error) {
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
				common.BitwiseRightShiftOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.BitwiseRightShiftOp.Name(),
					"",
					func(self BoolInstance, other common.Value) (common.Value, error) {
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
				common.BitwiseAndOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.BitwiseAndOp.Name(), "", func(self BoolInstance, other common.Value) (common.Value, error) {
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
				common.BitwiseXorOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.BitwiseXorOp.Name(), "", func(self BoolInstance, other common.Value) (common.Value, error) {
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
				common.BitwiseOrOp.Name(): newBoolBinaryOperator(
					// TODO: add doc
					common.BitwiseOrOp.Name(), "", func(self BoolInstance, other common.Value) (common.Value, error) {
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
			MakeLogicalOperators(Bool),
			MakeComparisonOperators(Bool, compareBooleans),
			MakeCommonOperators(Bool),
		)
	}

	return NewClass(
		common.BoolTypeName,
		nil,
		BuiltinPackage,
		initAttributes,
		func() (common.Value, error) {
			return NewBoolInstance(false), nil
		},
	)
}
