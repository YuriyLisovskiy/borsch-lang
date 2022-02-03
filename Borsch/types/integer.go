package types

import (
	"errors"
	"fmt"
	"math"
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
			CommonInstance{
				Object: Object{
					typeName:    common.IntegerTypeName,
					Attributes:  nil,
					callHandler: nil,
				},
				prototype: Integer,
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

func compareIntegers(_ common.State, self common.Type, other common.Type) (int, error) {
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
				"неможливо застосувати оператор '%s' до значень типів '%s' та '%s'",
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
	handler func(IntegerInstance, common.Type) (common.Type, error),
) *FunctionInstance {
	return newBinaryMethod(
		name,
		Integer,
		Any,
		doc,
		func(_ common.State, left common.Type, right common.Type) (common.Type, error) {
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
	handler func(IntegerInstance) (common.Type, error),
) *FunctionInstance {
	return newUnaryMethod(
		name, Integer, Any, doc, func(_ common.State, left common.Type) (common.Type, error) {
			if leftInstance, ok := left.(IntegerInstance); ok {
				return handler(leftInstance)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newIntegerClass() *Class {
	initAttributes := func() map[string]common.Type {
		return MergeAttributes(
			map[string]common.Type{
				// TODO: add doc
				common.ConstructorName: newBuiltinConstructor(Integer, ToInteger, ""),
				common.PowOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.PowOp.Name(), "", func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.UnaryPlus.Name(): newIntegerUnaryOperator(
					// TODO: add doc
					common.UnaryPlus.Name(), "", func(self IntegerInstance) (common.Type, error) {
						return self, nil
					},
				),
				common.UnaryMinus.Name(): newIntegerUnaryOperator(
					// TODO: add doc
					common.UnaryMinus.Name(), "", func(self IntegerInstance) (common.Type, error) {
						return NewIntegerInstance(-self.Value), nil
					},
				),
				common.UnaryBitwiseNotOp.Name(): newIntegerUnaryOperator(
					// TODO: add doc
					common.UnaryBitwiseNotOp.Name(), "", func(self IntegerInstance) (common.Type, error) {
						return NewIntegerInstance(^self.Value), nil
					},
				),
				common.MulOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.MulOp.Name(), "", func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.DivOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.DivOp.Name(), "", func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.ModuloOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.ModuloOp.Name(), "", func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.AddOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.AddOp.Name(), "", func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.SubOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.SubOp.Name(), "", func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.BitwiseLeftShiftOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.BitwiseLeftShiftOp.Name(),
					"",
					func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.BitwiseRightShiftOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.BitwiseRightShiftOp.Name(),
					"",
					func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.BitwiseAndOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.BitwiseAndOp.Name(), "", func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.BitwiseXorOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.BitwiseXorOp.Name(), "", func(self IntegerInstance, other common.Type) (common.Type, error) {
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
				common.BitwiseOrOp.Name(): newIntegerBinaryOperator(
					// TODO: add doc
					common.BitwiseOrOp.Name(), "", func(self IntegerInstance, other common.Type) (common.Type, error) {
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
			MakeLogicalOperators(Integer),
			MakeComparisonOperators(Integer, compareIntegers),
			MakeCommonOperators(Integer),
		)
	}

	return NewBuiltinClass(
		common.IntegerTypeName,
		BuiltinPackage,
		initAttributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewIntegerInstance(0), nil
		},
	)
}
