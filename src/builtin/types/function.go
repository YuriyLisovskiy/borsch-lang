package types

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/models"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

type FunctionParameter struct {
	TypeHash   int
	Name       string
	IsVariadic bool
}

func (fp FunctionParameter) TypeName() string {
	return GetTypeName(fp.TypeHash)
}

type FunctionType struct {
	Name       string
	Parameters []FunctionParameter
	Code       []models.Token
	Callable   func([]ValueType, map[string]ValueType) (ValueType, error)
	ReturnType int
	IsBuiltin bool
}

func NewFunctionType(
	name string, parameters []FunctionParameter, code []models.Token, returnType int,
) FunctionType {
	return FunctionType{
		Name:       name,
		Parameters: parameters,
		Code:       code,
		Callable:   nil,
		ReturnType: returnType,
		IsBuiltin: false,
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
