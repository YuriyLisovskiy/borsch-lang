package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type PackageType struct {
	IsBuiltin bool
	Name      string
	Parent    string
	Object    *ObjectType
}

func NewPackageType(isBuiltin bool, name, parent string, attrs map[string]ValueType) PackageType {
	attrs["__документ__"] = &NilType{} // TODO: set doc
	return PackageType{
		IsBuiltin: isBuiltin,
		Name:      name,
		Parent:    parent,
		Object:    newObjectType(PackageTypeHash, attrs),
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

func (t PackageType) TypeHash() int {
	return PackageTypeHash
}

func (t PackageType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t PackageType) AsBool() bool {
	return true
}

func (t PackageType) GetAttr(name string) (ValueType, error) {
	return t.Object.GetAttribute(name)
}

// SetAttr assumes that attribute already exists.
func (t PackageType) SetAttr(name string, value ValueType) (ValueType, error) {
	err := t.Object.SetAttribute(name, value)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t PackageType) Pow(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Plus() (ValueType, error) {
	return nil, nil
}

func (t PackageType) Minus() (ValueType, error) {
	return nil, nil
}

func (t PackageType) BitwiseNot() (ValueType, error) {
	return nil, nil
}

func (t PackageType) Mul(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Div(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Mod(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Add(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Sub(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) BitwiseLeftShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) BitwiseRightShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) BitwiseAnd(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) BitwiseXor(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) BitwiseOr(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
	case PackageType:
		return -2, util.RuntimeError(
			fmt.Sprintf(
				"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
				"%s", t.TypeName(), right.TypeName(),
			),
		)
	default:
		return -2, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				"%s", t.TypeName(), right.TypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t PackageType) Not() (ValueType, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t PackageType) And(other ValueType) (ValueType, error) {
	return BoolType{Value: other.AsBool()}, nil
}

func (t PackageType) Or(ValueType) (ValueType, error) {
	return BoolType{Value: true}, nil
}
