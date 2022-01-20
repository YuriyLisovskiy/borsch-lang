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
	Object
	Value float64
}

func NewRealInstance(value float64) RealInstance {
	return RealInstance{
		Value: value,
		Object: Object{
			typeName:    GetTypeName(RealTypeHash),
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t RealInstance) String(common.Context) string {
	return strconv.FormatFloat(t.Value, 'f', -1, 64)
}

func (t RealInstance) Representation(ctx common.Context) string {
	return t.String(ctx)
}

func (t RealInstance) GetTypeHash() uint64 {
	return t.GetPrototype().GetTypeHash()
}

func (t RealInstance) AsBool(common.Context) bool {
	return t.Value != 0.0
}

func (t RealInstance) SetAttribute(name string, _ common.Type) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	if t.Object.HasAttribute(name) || t.GetPrototype().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t RealInstance) GetAttribute(name string) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetPrototype().GetAttribute(name)
}

func (RealInstance) GetPrototype() *Class {
	return Real
}

func (t RealInstance) Div(other common.Type) (common.Type, error) {
	switch o := other.(type) {
	case BoolInstance:
		if o.Value {
			return NewRealInstance(t.Value), nil
		}
	case IntegerInstance:
		if o.Value != 0 {
			return NewRealInstance(t.Value / float64(o.Value)), nil
		}
	case RealInstance:
		if o.Value != 0.0 {
			return NewRealInstance(t.Value / o.Value), nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення на нуль")
}

func compareReals(_ common.Context, self, other common.Type) (int, error) {
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
		RealTypeHash,
		AnyTypeHash,
		doc,
		func(ctx common.Context, left common.Type, right common.Type) (common.Type, error) {
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
		name, RealTypeHash, AnyTypeHash, doc, func(ctx common.Context, left common.Type) (common.Type, error) {
			if leftInstance, ok := left.(RealInstance); ok {
				return handler(leftInstance)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newRealClass() *Class {
	attributes := mergeAttributes(
		map[string]common.Type{
			// TODO: add doc
			ops.ConstructorName: newBuiltinConstructor(RealTypeHash, ToReal, ""),
			ops.PowOp.Caption(): newRealBinaryOperator(
				// TODO: add doc
				ops.PowOp.Caption(), "", func(self RealInstance, other common.Type) (common.Type, error) {
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
			ops.UnaryPlus.Caption(): newRealUnaryOperator(
				// TODO: add doc
				ops.UnaryPlus.Caption(), "", func(self RealInstance) (common.Type, error) {
					return self, nil
				},
			),
			ops.UnaryMinus.Caption(): newRealUnaryOperator(
				// TODO: add doc
				ops.UnaryMinus.Caption(), "", func(self RealInstance) (common.Type, error) {
					return NewRealInstance(-self.Value), nil
				},
			),
			ops.MulOp.Caption(): newRealBinaryOperator(
				// TODO: add doc
				ops.MulOp.Caption(), "", func(self RealInstance, other common.Type) (common.Type, error) {
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
			ops.DivOp.Caption(): newRealBinaryOperator(
				// TODO: add doc
				ops.DivOp.Caption(), "", func(self RealInstance, other common.Type) (common.Type, error) {
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
			ops.AddOp.Caption(): newRealBinaryOperator(
				// TODO: add doc
				ops.AddOp.Caption(), "", func(self RealInstance, other common.Type) (common.Type, error) {
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
			ops.SubOp.Caption(): newRealBinaryOperator(
				// TODO: add doc
				ops.SubOp.Caption(), "", func(self RealInstance, other common.Type) (common.Type, error) {
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
		makeLogicalOperators(RealTypeHash),
		makeComparisonOperators(RealTypeHash, compareReals),
		makeCommonOperators(RealTypeHash),
	)
	return NewBuiltinClass(
		RealTypeHash,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewRealInstance(0), nil
		},
	)
}
