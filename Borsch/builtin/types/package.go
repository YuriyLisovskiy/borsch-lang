package types

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/Borsch/util"
)

type PackageType struct {
	IsBuiltin  bool
	Name       string
	Parent     string
	Attributes map[string]ValueType
}

func NewPackageType(isBuiltin bool, name, parent string, attrs map[string]ValueType) PackageType {
	return PackageType{
		IsBuiltin:  isBuiltin,
		Name:       name,
		Parent:     parent,
		Attributes: attrs,
	}
}

func (t PackageType) String() string {
	builtinStr := ""
	if t.IsBuiltin {
		builtinStr = " (вбудований)"
	}

	return fmt.Sprintf("<пакет '%s'%s>", t.Name, builtinStr)
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

func (t PackageType) GetAttr(name string) (ValueType, error) {
	if name == "__атрибути__" {
		dict := NewDictionaryType()
		for key, val := range t.Attributes {
			err := dict.SetElement(StringType{key}, val)
			if err != nil {
				return nil, err
			}
		}

		return dict, nil
	}

	if val, ok := t.Attributes[name]; ok {
		return val, nil
	}

	return nil, util.AttributeError(t.TypeName(), name)
}

// SetAttr assumes that attribute already exists.
func (t PackageType) SetAttr(name string, value ValueType) (ValueType, error) {
	if val, ok := t.Attributes[name]; ok {
		if val.TypeHash() == value.TypeHash() {
			t.Attributes[name] = value
			return t, nil
		}

		return nil, util.RuntimeError(fmt.Sprintf(
			"неможливо записати значення типу '%s' у атрибут '%s' з типом '%s'",
			value.TypeName(), name, val.TypeName(),
		))
	}

	t.Attributes[name] = value
	return t, nil

	//return nil, util.AttributeError(t.TypeName(), name)
}

func (t PackageType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
	case PackageType:
		return -2, util.RuntimeError(fmt.Sprintf(
			"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	default:
		return -2, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t PackageType) Add(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Sub(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Mul(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Div(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Pow(ValueType) (ValueType, error) {
	return nil, nil
}

func (t PackageType) Mod(ValueType) (ValueType, error) {
	return nil, nil
}
