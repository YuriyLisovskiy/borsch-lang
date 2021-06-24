package types

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

type FunctionArgument struct {
	TypeHash   int
	Name       string
	IsVariadic bool
}

func (fp FunctionArgument) TypeName() string {
	return GetTypeName(fp.TypeHash)
}

type FunctionType struct {
	Name       string
	Arguments  []FunctionArgument
	Callable   func([]ValueType, map[string]ValueType) (ValueType, error)
	ReturnType int
	IsBuiltin  bool
	Attributes map[string]ValueType
}

func NewFunctionType(
	name string, arguments []FunctionArgument, returnType int,
	fn func([]ValueType, map[string]ValueType) (ValueType, error),
) FunctionType {
	function := FunctionType{
		Name:       name,
		Arguments:  arguments,
		Callable:   fn,
		ReturnType: returnType,
		IsBuiltin:  false,
		Attributes: map[string]ValueType{},
	}

	return function
}

func (t FunctionType) String() string {
	template := "функція '%s' з типом результату '%s'"
	if t.IsBuiltin {
		template = "вбудована " + template
	}

	return fmt.Sprintf(template, t.Name, GetTypeName(t.ReturnType))
}

func (t FunctionType) Representation() string {
	return t.String()
}

func (t FunctionType) TypeHash() int {
	return FunctionTypeHash
}

func (t FunctionType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t FunctionType) GetAttr(name string) (ValueType, error) {
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

func (t FunctionType) SetAttr(name string, value ValueType) (ValueType, error) {
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
}
