package types

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

type PackageType struct {
	IsBuiltin  bool
	Name       string
	Attributes map[string]ValueType
}

func NewPackageType(isBuiltin bool, name string, attrs map[string]ValueType) PackageType {
	return PackageType{
		IsBuiltin:  isBuiltin,
		Name:       name,
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

	return nil, util.AttributeError(t.TypeName(), name)
}
