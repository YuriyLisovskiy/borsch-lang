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

type IntegerType struct {
	Object

	Value    int64
	package_ *PackageType
}

func NewIntegerType(value string) (IntegerType, error) {
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return IntegerType{}, util.RuntimeError(err.Error())
	}

	return IntegerType{
		Object:   *newIntegerObject(),
		Value:    number,
		package_: BuiltinPackage,
	}, nil
}

func (t IntegerType) String() string {
	return fmt.Sprintf("%d", t.Value)
}

func (t IntegerType) Representation() string {
	return t.String()
}

func (t IntegerType) AsBool() bool {
	return t.Value != 0
}

func (t IntegerType) SetAttribute(name string, _ Type) (Type, error) {
	if t.Object.HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func compareTo(self IntegerType, other Type) (int, error) {
	switch right := other.(type) {
	case NilType:
	case BoolType:
		rightVal := boolToInt64(right.Value)
		if self.Value == rightVal {
			return 0, nil
		}

		if self.Value < rightVal {
			return -1, nil
		}

		return 1, nil
	case IntegerType:
		if self.Value == right.Value {
			return 0, nil
		}

		if self.Value < right.Value {
			return -1, nil
		}

		return 1, nil
	case RealType:
		leftVal := float64(self.Value)
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

func newIntegerBinaryMethod(
	name string,
	doc string,
	handler func(IntegerType, Type) (Type, error),
) FunctionType {
	return newBinaryMethod(
		name, IntegerTypeHash, doc, func(left Type, right Type) (Type, error) {
			return handler(left.(IntegerType), right)
		},
	)
}

func newIntegerUnaryMethod(
	name string,
	doc string,
	handler func(IntegerType) (Type, error),
) FunctionType {
	return newUnaryMethod(
		name, IntegerTypeHash, doc, func(left Type) (Type, error) {
			return handler(left.(IntegerType))
		},
	)
}

func newComparisonMethod(operator ops.Operator, doc string, checker func(res int) bool) FunctionType {
	return newIntegerBinaryMethod(
		operator.Caption(), doc, func(self IntegerType, other Type) (Type, error) {
			res, err := compareTo(self, other)
			if err != nil {
				return nil, err
			}

			return BoolType{Value: checker(res)}, nil
		},
	)
}

func newIntegerObject() *Object {
	return newBuiltinObject(
		IntegerTypeHash,
		map[string]Type{
			"__документ__": &NilType{}, // TODO: set doc
			"__пакет__":    BuiltinPackage,
			ops.PowOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.PowOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case RealType:
						return RealType{
							Value: math.Pow(float64(self.Value), o.Value),
						}, nil
					case IntegerType:
						return IntegerType{
							Value: int64(math.Pow(float64(self.Value), float64(o.Value))),
						}, nil
					case BoolType:
						return IntegerType{
							Value: int64(math.Pow(float64(self.Value), boolToFloat64(o.Value))),
						}, nil
					default:
						return nil, nil
					}
				},
			),
			ops.UnaryPlus.Caption(): newIntegerUnaryMethod(
				// TODO: add doc
				ops.UnaryPlus.Caption(), "", func(self IntegerType) (Type, error) {
					return self, nil
				},
			),
			ops.UnaryMinus.Caption(): newIntegerUnaryMethod(
				// TODO: add doc
				ops.UnaryMinus.Caption(), "", func(self IntegerType) (Type, error) {
					return IntegerType{Value: -self.Value}, nil
				},
			),
			ops.UnaryBitwiseNotOp.Caption(): newIntegerUnaryMethod(
				// TODO: add doc
				ops.UnaryBitwiseNotOp.Caption(), "", func(self IntegerType) (Type, error) {
					return IntegerType{Value: ^self.Value}, nil
				},
			),
			ops.MulOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.MulOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						return IntegerType{
							Value: self.Value * boolToInt64(o.Value),
						}, nil
					case IntegerType:
						return IntegerType{
							Value: self.Value * o.Value,
						}, nil
					case RealType:
						return RealType{
							Value: float64(self.Value) * o.Value,
						}, nil
					case StringType:
						count := int(self.Value)
						if count <= 0 {
							return StringType{Value: ""}, nil
						}

						return StringType{
							Value: strings.Repeat(o.Value, count),
						}, nil
					case ListType:
						count := int(self.Value)
						list := NewListType()
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
			ops.DivOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.DivOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						if o.Value {
							return RealType{
								Value: float64(self.Value),
							}, nil
						}
					case IntegerType:
						if o.Value != 0 {
							return RealType{
								Value: float64(self.Value) / float64(o.Value),
							}, nil
						}
					case RealType:
						if o.Value != 0.0 {
							return RealType{
								Value: float64(self.Value) / o.Value,
							}, nil
						}
					default:
						return nil, nil
					}

					return nil, errors.New("ділення на нуль")
				},
			),
			ops.ModuloOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.ModuloOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						if o.Value {
							return IntegerType{
								Value: self.Value % boolToInt64(o.Value),
							}, nil
						}
					case IntegerType:
						if o.Value != 0 {
							return IntegerType{
								Value: self.Value % o.Value,
							}, nil
						}
					default:
						return nil, nil
					}

					return nil, errors.New("ділення за модулем на нуль")
				},
			),
			ops.AddOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.AddOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						return IntegerType{
							Value: self.Value + boolToInt64(o.Value),
						}, nil
					case IntegerType:
						return IntegerType{
							Value: self.Value + o.Value,
						}, nil
					case RealType:
						return RealType{
							Value: float64(self.Value) + o.Value,
						}, nil
					default:
						return nil, nil
					}
				},
			),
			ops.SubOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.SubOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						return IntegerType{
							Value: self.Value - boolToInt64(o.Value),
						}, nil
					case IntegerType:
						return IntegerType{
							Value: self.Value - o.Value,
						}, nil
					case RealType:
						return RealType{
							Value: float64(self.Value) - o.Value,
						}, nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseLeftShiftOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.BitwiseLeftShiftOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						return IntegerType{Value: self.Value << boolToInt64(o.Value)}, nil
					case IntegerType:
						return IntegerType{Value: self.Value << o.Value}, nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseRightShiftOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.BitwiseRightShiftOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						return IntegerType{Value: self.Value >> boolToInt64(o.Value)}, nil
					case IntegerType:
						return IntegerType{Value: self.Value >> o.Value}, nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseAndOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.BitwiseAndOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						return IntegerType{Value: self.Value & boolToInt64(o.Value)}, nil
					case IntegerType:
						return IntegerType{Value: self.Value & o.Value}, nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseXorOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.BitwiseXorOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						return IntegerType{Value: self.Value ^ boolToInt64(o.Value)}, nil
					case IntegerType:
						return IntegerType{Value: self.Value ^ o.Value}, nil
					default:
						return nil, nil
					}
				},
			),
			ops.BitwiseOrOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.BitwiseOrOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					switch o := other.(type) {
					case BoolType:
						return IntegerType{Value: self.Value | boolToInt64(o.Value)}, nil
					case IntegerType:
						return IntegerType{Value: self.Value | o.Value}, nil
					default:
						return nil, nil
					}
				},
			),
			ops.NotOp.Caption(): newIntegerUnaryMethod(
				// TODO: add doc
				ops.NotOp.Caption(), "", func(self IntegerType) (Type, error) {
					return BoolType{Value: !self.AsBool()}, nil
				},
			),
			ops.AndOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.AndOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					return logicalAnd(self, other)
				},
			),
			ops.OrOp.Caption(): newIntegerBinaryMethod(
				// TODO: add doc
				ops.OrOp.Caption(), "", func(self IntegerType, other Type) (Type, error) {
					return logicalOr(self, other)
				},
			),
			ops.EqualsOp.Caption(): newComparisonMethod(
				// TODO: add doc
				ops.EqualsOp, "", func(res int) bool {
					return res == 0
				},
			),
			ops.NotEqualsOp.Caption(): newComparisonMethod(
				// TODO: add doc
				ops.NotEqualsOp, "", func(res int) bool {
					return res != 0
				},
			),
			ops.GreaterOp.Caption(): newComparisonMethod(
				// TODO: add doc
				ops.GreaterOp, "", func(res int) bool {
					return res == 1
				},
			),
			ops.GreaterOrEqualsOp.Caption(): newComparisonMethod(
				// TODO: add doc
				ops.GreaterOrEqualsOp, "", func(res int) bool {
					return res == 0 || res == 1
				},
			),
			ops.LessOp.Caption(): newComparisonMethod(
				// TODO: add doc
				ops.LessOp, "", func(res int) bool {
					return res == -1
				},
			),
			ops.LessOrEqualsOp.Caption(): newComparisonMethod(
				// TODO: add doc
				ops.LessOrEqualsOp, "", func(res int) bool {
					return res == 0 || res == -1
				},
			),
		},
	)
}
