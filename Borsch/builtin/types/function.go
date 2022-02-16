package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type FunctionParameter struct {
	// nil means any type
	Type       *Class
	Name       string
	IsVariadic bool
	IsNullable bool
}

func (fa *FunctionParameter) String() string {
	res := fa.Name
	if fa.IsVariadic {
		res += "..."
	}

	return res + ": " + fa.GetTypeName()
}

func (fa FunctionParameter) GetTypeName() string {
	res := ""
	if fa.Type != Any {
		res = fa.Type.GetTypeName()
	} else {
		res = common.AnyTypeName
	}

	if fa.IsNullable {
		res += "?"
	}

	return res
}

type FunctionReturnType struct {
	Type       *Class
	IsNullable bool
}

func (r *FunctionReturnType) String() string {
	res := r.GetTypeName()
	if r.IsNullable {
		res += "?"
	}

	return res
}

func (r *FunctionReturnType) GetTypeName() string {
	if r.Type != Any {
		return r.Type.GetName()
	}

	return common.AnyTypeName
}

type FunctionInstance struct {
	ClassInstance
	package_    *PackageInstance
	address     string
	Name        string
	Parameters  []FunctionParameter
	ReturnTypes []FunctionReturnType
	IsMethod    bool
	callFunc    func(common.State, *[]common.Value, *map[string]common.Value) (common.Value, error)
}

func NewFunctionInstance(
	name string,
	arguments []FunctionParameter,
	handler func(common.State, *[]common.Value, *map[string]common.Value) (common.Value, error),
	returnTypes []FunctionReturnType,
	isMethod bool,
	package_ *PackageInstance,
	doc string,
) *FunctionInstance {
	attributes := map[string]common.Value{}
	if package_ != nil {
		attributes[common.PackageAttributeName] = package_
	}

	if len(doc) > 0 {
		attributes[common.DocAttributeName] = NewStringInstance(doc)
	}

	function := &FunctionInstance{
		ClassInstance: ClassInstance{
			class:      Function,
			attributes: map[string]common.Value{},
			address:    "",
		},
		package_:    package_,
		Name:        name,
		Parameters:  arguments,
		ReturnTypes: returnTypes,
		IsMethod:    isMethod,
		callFunc:    handler,
	}

	function.address = fmt.Sprintf("%p", function)
	return function
}

func (i FunctionInstance) String(common.State) (string, error) {
	template := ""
	if i.Name == common.LambdaSignature {
		template = "функція " + common.LambdaSignature
	} else {
		if i.package_ != nil {
			template = "функція '%s'"
		} else {
			template = "метод '%s'"
		}

		template = fmt.Sprintf(template, i.Name)
	}

	template += " з адресою %s"
	return fmt.Sprintf(fmt.Sprintf("<%s>", template), i.address), nil
}

func (i FunctionInstance) Representation(state common.State) (string, error) {
	return i.String(state)
}

func (i FunctionInstance) AsBool(common.State) (bool, error) {
	return true, nil
}

func (i FunctionInstance) Call(state common.State, args *[]common.Value, kwargs *map[string]common.Value) (
	common.Value,
	error,
) {
	if i.callFunc != nil {
		return i.callFunc(state, args, kwargs)
	}

	return nil, utilities.ObjectIsNotCallable("", i.GetTypeName())
}

func (i *FunctionInstance) GetContext() common.Context {
	if i.package_ != nil {
		return i.package_.GetContext()
	}

	return nil
}

func (i *FunctionInstance) IsLambda() bool {
	return i.Name == common.LambdaSignature
}

func functionOperator_Call(name string) common.Value {
	return NewFunctionInstance(
		name,
		[]FunctionParameter{
			{
				Type:       Function,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       Any,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(state common.State, args *[]common.Value, kwargs *map[string]common.Value) (
			common.Value,
			error,
		) {
			function := (*args)[0].(*FunctionInstance)
			functionArgs := (*args)[1:]
			functionKwargs := *kwargs
			delete(functionKwargs, "я")
			if err := CheckFunctionArguments(function, &functionArgs, &functionKwargs); err != nil {
				return nil, err
			}

			return function.Call(state, &functionArgs, &functionKwargs)
		},
		[]FunctionReturnType{
			{
				Type:       Any,
				IsNullable: false,
			},
		},
		true,
		nil,
		"", // TODO: add doc
	)
}

func newFunctionClass() *Class {
	return &Class{
		Name:    common.FunctionTypeName,
		IsFinal: true,
		Bases:   []*Class{},
		Parent:  BuiltinPackage,
		AttrInitializer: func(attrs *map[string]common.Value) {
			*attrs = MergeAttributes(
				map[string]common.Value{
					common.CallOperatorName: functionOperator_Call(common.CallOperatorName),
				},
				MakeLogicalOperators(Function),
				MakeCommonOperators(Function),
			)
		},
		GetEmptyInstance: func() (common.Value, error) {
			panic("unreachable")
		},
	}
}
