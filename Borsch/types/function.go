package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
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

type FunctionInstance struct {
	Object
	package_    *PackageInstance
	address     string
	Name        string
	Arguments   []FunctionArgument
	ReturnTypes []FunctionReturnType
	IsMethod    bool
}

func NewFunctionInstance(
	name string,
	arguments []FunctionArgument,
	handler func(interface{}, *[]common.Type, *map[string]common.Type) (common.Type, error),
	returnTypes []FunctionReturnType,
	isMethod bool,
	package_ *PackageInstance,
	doc string,
) *FunctionInstance {
	attributes := map[string]common.Type{}
	if package_ != nil {
		attributes[ops.PackageAttributeName] = package_
	}

	if len(doc) > 0 {
		attributes[ops.DocAttributeName] = NewStringInstance(doc)
	}

	function := &FunctionInstance{
		Object: Object{
			typeName:    GetTypeName(FunctionTypeHash),
			Attributes:  attributes,
			callHandler: handler,
		},
		package_:    package_,
		Name:        name,
		Arguments:   arguments,
		ReturnTypes: returnTypes,
		IsMethod:    isMethod,
	}

	function.address = fmt.Sprintf("%p", function)
	return function
}

func (t FunctionInstance) String() string {
	template := ""
	if t.Name == "" {
		template = "функція <лямбда>"
	} else {
		if t.package_ != nil {
			template = "функція '%s'"
			if t.package_.IsBuiltin {
				template = "вбудована " + template
			}
		} else {
			template = "метод '%s'"
		}

		template = fmt.Sprintf(template, t.Name)
	}

	template += " з адресою %s"
	return fmt.Sprintf(fmt.Sprintf("<%s>", template), t.address)
}

func (t FunctionInstance) Representation() string {
	return t.String()
}

func (t FunctionInstance) GetTypeHash() uint64 {
	return t.GetClass().GetTypeHash()
}

func (t FunctionInstance) AsBool() bool {
	return true
}

func (t FunctionInstance) SetAttribute(name string, value common.Type) (common.Type, error) {
	err := t.Object.SetAttribute(name, value)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t FunctionInstance) GetAttribute(name string) (common.Type, error) {
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetClass().GetAttribute(name)
}

func (t FunctionInstance) GetClass() *Class {
	return Function
}

func newFunctionClass() *Class {
	attributes := mergeAttributes(
		map[string]common.Type{
			ops.CallOperatorName: NewFunctionInstance(
				ops.CallOperatorName,
				[]FunctionArgument{
					{
						TypeHash:   FunctionTypeHash,
						Name:       "я",
						IsVariadic: false,
						IsNullable: false,
					},
					{
						TypeHash:   AnyTypeHash,
						Name:       "значення",
						IsVariadic: true,
						IsNullable: true,
					},
				},
				func(_ interface{}, args *[]common.Type, kwargs *map[string]common.Type) (common.Type, error) {
					function := (*args)[0].(*FunctionInstance)
					slicedArgs := (*args)[1:]
					slicedKwargs := *kwargs
					delete(slicedKwargs, "я")
					if err := CheckFunctionArguments(function, &slicedArgs, &slicedKwargs); err != nil {
						return nil, err
					}

					return function.Call(nil, &slicedArgs, &slicedKwargs)
				},
				[]FunctionReturnType{
					{
						TypeHash:   AnyTypeHash,
						IsNullable: true,
					},
				},
				true,
				nil,
				"", // TODO: add doc
			),
		},
		makeLogicalOperators(FunctionTypeHash),
		makeCommonOperators(FunctionTypeHash),
	)
	return NewBuiltinClass(
		FunctionTypeHash,
		BuiltinPackage,
		attributes,
		"",  // TODO: add doc
		nil, // CAUTION: segfault may be thrown when using without nil check!
	)
}
