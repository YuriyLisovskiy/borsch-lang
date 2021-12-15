package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type FunctionArgument struct {
	TypeHash   uint64
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
	TypeHash   uint64
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
	Object

	package_   *PackageType
	Name       string
	Arguments  []FunctionArgument
	ReturnType FunctionReturnType
}

func NewFunctionType(
	name string,
	arguments []FunctionArgument,
	handler func([]Type, map[string]Type) (Type, error),
	returnType FunctionReturnType,
	package_ *PackageType,
	doc string,
) FunctionType {
	function := FunctionType{
		Name:       name,
		Arguments:  arguments,
		ReturnType: returnType,
		package_:   package_,
		Object: *newBuiltinObject(
			FunctionTypeHash, map[string]Type{
				"__документ__": &StringType{Value: doc},
				"__пакет__":    package_,
			},
		),
	}

	function.CallHandler = handler
	return function
}

func (t FunctionType) String() string {
	template := ""
	if t.package_ != nil {
		template = "функція '%s'"
		if t.package_.IsBuiltin {
			template = "вбудована " + template
		}
	} else {
		template = "метод '%s'"
	}

	template += " з типом результату '%s'"
	return fmt.Sprintf(fmt.Sprintf("<%s>", template), t.Name, t.ReturnType.String())
}

func (t FunctionType) Representation() string {
	return t.String()
}

func (t FunctionType) AsBool() bool {
	return true
}

func (t FunctionType) SetAttribute(name string, value Type) (Type, error) {
	err := t.Object.SetAttribute(name, value)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t FunctionType) Pow(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) Plus() (Type, error) {
	return nil, nil
}

func (t FunctionType) Minus() (Type, error) {
	return nil, nil
}

func (t FunctionType) BitwiseNot() (Type, error) {
	return nil, nil
}

func (t FunctionType) Mul(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) Div(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) Mod(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) Add(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) Sub(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) BitwiseLeftShift(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) BitwiseRightShift(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) BitwiseAnd(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) BitwiseXor(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) BitwiseOr(Type) (Type, error) {
	return nil, nil
}

func (t FunctionType) CompareTo(other Type) (int, error) {
	switch right := other.(type) {
	case NilType:
	case FunctionType:
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

func (t FunctionType) Not() (Type, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t FunctionType) And(other Type) (Type, error) {
	return BoolType{Value: other.AsBool()}, nil
}

func (t FunctionType) Or(Type) (Type, error) {
	return BoolType{Value: true}, nil
}
