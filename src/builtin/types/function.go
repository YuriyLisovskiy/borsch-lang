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
}

func NewFunctionType(
	name string, arguments []FunctionArgument, returnType int,
	fn func([]ValueType, map[string]ValueType) (ValueType, error),
) FunctionType {
	return FunctionType{
		Name:       name,
		Arguments:  arguments,
		Callable:   fn,
		ReturnType: returnType,
		IsBuiltin:  false,
	}
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
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t FunctionType) SetAttr(name string, _ ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}
