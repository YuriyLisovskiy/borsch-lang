package types

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type RealInstance struct {
	BuiltinObject
	Value float64
}

func NewRealInstance(value float64) RealInstance {
	return RealInstance{
		BuiltinObject: BuiltinObject{
			CommonObject{
				Object: Object{
					typeName:    common.RealTypeName,
					Attributes:  nil,
					callHandler: nil,
				},
				prototype: Real,
			},
		},
		Value: value,
	}
}

func (t RealInstance) String(common.State) string {
	return strconv.FormatFloat(t.Value, 'f', -1, 64)
}

func (t RealInstance) Representation(state common.State) string {
	return t.String(state)
}

func (t RealInstance) AsBool(common.State) bool {
	return t.Value != 0.0
}

func compareReals(_ common.State, self, other common.Type) (int, error) {
	left, ok := self.(RealInstance)
	if !ok {
		return 0, util.IncorrectUseOfFunctionError("compareReals")
	}

	switch right := other.(type) {
	case NilInstance:
	case BoolInstance:
		rightVal := boolToFloat64(right.Value)
		if left.Value == rightVal {
			return 0, nil
		}

		if left.Value < rightVal {
			return -1, nil
		}

		return 1, nil
	case IntegerInstance:
		rightVal := float64(right.Value)
		if left.Value == rightVal {
			return 0, nil
		}

		if left.Value < rightVal {
			return -1, nil
		}

		return 1, nil
	case RealInstance:
		if left.Value == right.Value {
			return 0, nil
		}

		if left.Value < right.Value {
			return -1, nil
		}

		return 1, nil
	default:
		return 0, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор '%s' до значень типів '%s' та '%s'",
				"%s", left.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newRealBinaryOperator(
	name string,
	doc string,
	handler func(RealInstance, common.Type) (common.Type, error),
) *FunctionInstance {
	return newBinaryMethod(
		name,
		Real,
		Any,
		doc,
		func(_ common.State, left common.Type, right common.Type) (common.Type, error) {
			if leftInstance, ok := left.(RealInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newRealUnaryOperator(
	name string,
	doc string,
	handler func(RealInstance) (common.Type, error),
) *FunctionInstance {
	return newUnaryMethod(
		name, Real, Any, doc, func(_ common.State, left common.Type) (common.Type, error) {
			if leftInstance, ok := left.(RealInstance); ok {
				return handler(leftInstance)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newRealClass() *Class {
	initAttributes := func() map[string]common.Type {
		return mergeAttributes(
			map[string]common.Type{
				// TODO: add doc
				ops.ConstructorName: newBuiltinConstructor(Real, ToReal, ""),
				ops.PowOp.Name(): newRealBinaryOperator(
					// TODO: add doc
					ops.PowOp.Name(), "", func(self RealInstance, other common.Type) (common.Type, error) {
						switch o := other.(type) {
						case RealInstance:
							return NewRealInstance(math.Pow(self.Value, o.Value)), nil
						case IntegerInstance:
							return NewRealInstance(math.Pow(self.Value, float64(o.Value))), nil
						case BoolInstance:
							return NewRealInstance(math.Pow(self.Value, boolToFloat64(o.Value))), nil
						default:
							return nil, nil
						}
					},
				),
				ops.UnaryPlus.Name(): newRealUnaryOperator(
					// TODO: add doc
					ops.UnaryPlus.Name(), "", func(self RealInstance) (common.Type, error) {
						return self, nil
					},
				),
				ops.UnaryMinus.Name(): newRealUnaryOperator(
					// TODO: add doc
					ops.UnaryMinus.Name(), "", func(self RealInstance) (common.Type, error) {
						return NewRealInstance(-self.Value), nil
					},
				),
				ops.MulOp.Name(): newRealBinaryOperator(
					// TODO: add doc
					ops.MulOp.Name(), "", func(self RealInstance, other common.Type) (common.Type, error) {
						switch o := other.(type) {
						case BoolInstance:
							return NewRealInstance(self.Value * boolToFloat64(o.Value)), nil
						case IntegerInstance:
							return NewRealInstance(self.Value * float64(o.Value)), nil
						case RealInstance:
							return NewRealInstance(self.Value * o.Value), nil
						default:
							return nil, nil
						}
					},
				),
				ops.DivOp.Name(): newRealBinaryOperator(
					// TODO: add doc
					ops.DivOp.Name(), "", func(self RealInstance, other common.Type) (common.Type, error) {
						switch o := other.(type) {
						case BoolInstance:
							if o.Value {
								return NewRealInstance(self.Value), nil
							}
						case IntegerInstance:
							if o.Value != 0 {
								return NewRealInstance(self.Value / float64(o.Value)), nil
							}
						case RealInstance:
							if o.Value != 0.0 {
								return NewRealInstance(self.Value / o.Value), nil
							}
						default:
							return nil, nil
						}

						return nil, errors.New("ділення на нуль")
					},
				),
				ops.AddOp.Name(): newRealBinaryOperator(
					// TODO: add doc
					ops.AddOp.Name(), "", func(self RealInstance, other common.Type) (common.Type, error) {
						switch o := other.(type) {
						case BoolInstance:
							return NewRealInstance(self.Value + boolToFloat64(o.Value)), nil
						case IntegerInstance:
							return NewRealInstance(self.Value + float64(o.Value)), nil
						case RealInstance:
							return NewRealInstance(self.Value + o.Value), nil
						default:
							return nil, nil
						}
					},
				),
				ops.SubOp.Name(): newRealBinaryOperator(
					// TODO: add doc
					ops.SubOp.Name(), "", func(self RealInstance, other common.Type) (common.Type, error) {
						switch o := other.(type) {
						case BoolInstance:
							return NewRealInstance(self.Value - boolToFloat64(o.Value)), nil
						case IntegerInstance:
							return NewRealInstance(self.Value - float64(o.Value)), nil
						case RealInstance:
							return NewRealInstance(self.Value - o.Value), nil
						default:
							return nil, nil
						}
					},
				),
			},
			makeLogicalOperators(Real),
			makeComparisonOperators(Real, compareReals),
			makeCommonOperators(Real),
		)
	}

	return NewBuiltinClass(
		common.RealTypeName,
		BuiltinPackage,
		initAttributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewRealInstance(0), nil
		},
	)
}
