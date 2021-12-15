package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type PackageType struct {
	Object

	IsBuiltin bool
	Name      string
	Parent    string
}

func NewPackageType(isBuiltin bool, name, parent string, attrs map[string]Type) PackageType {
	attrs["__документ__"] = &NilType{} // TODO: set doc
	return PackageType{
		IsBuiltin: isBuiltin,
		Name:      name,
		Parent:    parent,
		Object:    *newBuiltinObject(PackageTypeHash, attrs),
	}
}

func (t PackageType) String() string {
	name := t.Name
	if t.IsBuiltin {
		name = "АТБ"
	}

	return fmt.Sprintf("<пакет '%s'>", name)
}

func (t PackageType) Representation() string {
	return t.String()
}

func (t PackageType) AsBool() bool {
	return true
}

func (t PackageType) SetAttribute(name string, value Type) (Type, error) {
	err := t.Object.SetAttribute(name, value)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t PackageType) Pow(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) Plus() (Type, error) {
	return nil, nil
}

func (t PackageType) Minus() (Type, error) {
	return nil, nil
}

func (t PackageType) BitwiseNot() (Type, error) {
	return nil, nil
}

func (t PackageType) Mul(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) Div(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) Mod(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) Add(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) Sub(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) BitwiseLeftShift(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) BitwiseRightShift(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) BitwiseAnd(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) BitwiseXor(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) BitwiseOr(Type) (Type, error) {
	return nil, nil
}

func (t PackageType) CompareTo(other Type) (int, error) {
	switch right := other.(type) {
	case NilType:
	case PackageType:
		return -2, util.RuntimeError(
			fmt.Sprintf(
				"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
				"%s", t.GetTypeName(), right.GetTypeName(),
			),
		)
	default:
		return -2, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				"%s", t.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t PackageType) Not() (Type, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t PackageType) And(other Type) (Type, error) {
	return BoolType{Value: other.AsBool()}, nil
}

func (t PackageType) Or(Type) (Type, error) {
	return BoolType{Value: true}, nil
}
