package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type NilType struct {
}

func (t NilType) String() string {
	return "нуль"
}

func (t NilType) Representation() string {
	return t.String()
}

func (t NilType) TypeHash() int {
	return NilTypeHash
}

func (t NilType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t NilType) AsBool() bool {
	return false
}

func (t NilType) GetAttr(name string) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t NilType) SetAttr(name string, _ ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t NilType) Pow(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) Plus() (ValueType, error) {
	return nil, nil
}

func (t NilType) Minus() (ValueType, error) {
	return nil, nil
}

func (t NilType) BitwiseNot() (ValueType, error) {
	return nil, nil
}

func (t NilType) Mul(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) Div(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) Mod(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) Add(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) Sub(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) BitwiseLeftShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) BitwiseRightShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) BitwiseAnd(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) BitwiseXor(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) BitwiseOr(ValueType) (ValueType, error) {
	return nil, nil
}

func (t NilType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
		return 0, nil
	default:
		return 0, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	}
}

func (t NilType) Not() (ValueType, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t NilType) And(ValueType) (ValueType, error) {
	return BoolType{Value: false}, nil
}

func (t NilType) Or(other ValueType) (ValueType, error) {
	return BoolType{Value: other.AsBool()}, nil
}
