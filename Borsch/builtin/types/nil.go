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

func (t NilType) GetTypeHash() uint64 {
	return NilTypeHash
}

func (t NilType) GetTypeName() string {
	return GetTypeName(t.GetTypeHash())
}

func (t NilType) AsBool() bool {
	return false
}

func (t NilType) GetAttribute(name string) (Type, error) {
	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t NilType) SetAttribute(name string, _ Type) (Type, error) {
	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t NilType) Pow(Type) (Type, error) {
	return nil, nil
}

func (t NilType) Plus() (Type, error) {
	return nil, nil
}

func (t NilType) Minus() (Type, error) {
	return nil, nil
}

func (t NilType) BitwiseNot() (Type, error) {
	return nil, nil
}

func (t NilType) Mul(Type) (Type, error) {
	return nil, nil
}

func (t NilType) Div(Type) (Type, error) {
	return nil, nil
}

func (t NilType) Mod(Type) (Type, error) {
	return nil, nil
}

func (t NilType) Add(Type) (Type, error) {
	return nil, nil
}

func (t NilType) Sub(Type) (Type, error) {
	return nil, nil
}

func (t NilType) BitwiseLeftShift(Type) (Type, error) {
	return nil, nil
}

func (t NilType) BitwiseRightShift(Type) (Type, error) {
	return nil, nil
}

func (t NilType) BitwiseAnd(Type) (Type, error) {
	return nil, nil
}

func (t NilType) BitwiseXor(Type) (Type, error) {
	return nil, nil
}

func (t NilType) BitwiseOr(Type) (Type, error) {
	return nil, nil
}

func (t NilType) CompareTo(other Type) (int, error) {
	switch right := other.(type) {
	case NilType:
		return 0, nil
	default:
		return 0, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.GetTypeName(), right.GetTypeName(),
		))
	}
}

func (t NilType) Not() (Type, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t NilType) And(Type) (Type, error) {
	return BoolType{Value: false}, nil
}

func (t NilType) Or(other Type) (Type, error) {
	return BoolType{Value: other.AsBool()}, nil
}
