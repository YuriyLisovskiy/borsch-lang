package types

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/Borsch/util"
)

type FunctionArgument struct {
	TypeHash   int
	Name       string
	IsVariadic bool
	IsNullable bool
}

func (fa *FunctionArgument) String() string {
	res := fa.Name
	if fa.IsVariadic {
		res += "..."
	}

	res += GetTypeName(fa.TypeHash)
	if fa.IsNullable {
		res += "?"
	}

	return res
}

func (fa FunctionArgument) TypeName() string {
	res := GetTypeName(fa.TypeHash)
	if fa.IsNullable {
		res += "?"
	}

	return res
}

type FunctionReturnType struct {
	TypeHash   int
	IsNullable bool
}

func (r *FunctionReturnType) String() string {
	res := GetTypeName(r.TypeHash)
	if r.IsNullable {
		res += "?"
	}

	return res
}

type FunctionType struct {
	Name       string
	Arguments  []FunctionArgument
	Callable   func([]ValueType, map[string]ValueType) (ValueType, error)
	ReturnType FunctionReturnType
	IsBuiltin  bool
	Attributes map[string]ValueType
}

func NewFunctionType(
	name string, arguments []FunctionArgument, returnType FunctionReturnType,
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

	return fmt.Sprintf(template, t.Name, t.ReturnType.String())
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

func (t FunctionType) AsBool() bool {
	return true
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

func (t FunctionType) Pow(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) Plus() (ValueType, error) {
	return nil, nil
}

func (t FunctionType) Minus() (ValueType, error) {
	return nil, nil
}

func (t FunctionType) BitwiseNot() (ValueType, error) {
	return nil, nil
}

func (t FunctionType) Mul(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) Div(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) Mod(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) Add(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) Sub(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) BitwiseLeftShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) BitwiseRightShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) BitwiseAnd(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) BitwiseXor(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) BitwiseOr(ValueType) (ValueType, error) {
	return nil, nil
}

func (t FunctionType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
	case FunctionType:
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

func (t FunctionType) Not() (ValueType, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t FunctionType) And(other ValueType) (ValueType, error) {
	return BoolType{Value: other.AsBool()}, nil
}

func (t FunctionType) Or(ValueType) (ValueType, error) {
	return BoolType{Value: true}, nil
}
