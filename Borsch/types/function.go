package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
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
		return r.Type.GetTypeName()
	}

	return common.AnyTypeName
}

type FunctionInstance struct {
	CommonInstance
	package_    *PackageInstance
	address     string
	Name        string
	Parameters  []FunctionParameter
	ReturnTypes []FunctionReturnType
	IsMethod    bool
}

func NewFunctionInstance(
	name string,
	arguments []FunctionParameter,
	handler func(common.State, *[]common.Type, *map[string]common.Type) (common.Type, error),
	returnTypes []FunctionReturnType,
	isMethod bool,
	package_ *PackageInstance,
	doc string,
) *FunctionInstance {
	attributes := map[string]common.Type{}
	if package_ != nil {
		attributes[common.PackageAttributeName] = package_
	}

	if len(doc) > 0 {
		attributes[common.DocAttributeName] = NewStringInstance(doc)
	}

	function := &FunctionInstance{
		CommonInstance: CommonInstance{
			Object: Object{
				typeName:    common.FunctionTypeName,
				Attributes:  attributes,
				callHandler: handler,
			},
			prototype: Function,
		},
		package_:    package_,
		Name:        name,
		Parameters:  arguments,
		ReturnTypes: returnTypes,
		IsMethod:    isMethod,
	}

	function.address = fmt.Sprintf("%p", function)
	return function
}

func (t FunctionInstance) String(common.State) (string, error) {
	template := ""
	if t.Name == common.LambdaSignature {
		template = "функція " + common.LambdaSignature
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
	return fmt.Sprintf(fmt.Sprintf("<%s>", template), t.address), nil
}

func (t FunctionInstance) Representation(state common.State) (string, error) {
	return t.String(state)
}

func (t FunctionInstance) AsBool(common.State) (bool, error) {
	return true, nil
}

func (t *FunctionInstance) GetContext() common.Context {
	if t.package_ != nil {
		return t.package_.GetContext()
	}

	return nil
}

func (t *FunctionInstance) IsLambda() bool {
	return t.Name == common.LambdaSignature
}

func newFunctionClass() *Class {
	initAttributes := func() map[string]common.Type {
		return MergeAttributes(
			map[string]common.Type{
				common.CallOperatorName: NewFunctionInstance(
					common.CallOperatorName,
					[]FunctionParameter{
						{
							Type:       Function,
							Name:       "я",
							IsVariadic: false,
							IsNullable: false,
						},
						{
							Type:       nil,
							Name:       "значення",
							IsVariadic: true,
							IsNullable: true,
						},
					},
					func(state common.State, args *[]common.Type, kwargs *map[string]common.Type) (
						common.Type,
						error,
					) {
						function := (*args)[0].(*FunctionInstance)
						slicedArgs := (*args)[1:]
						slicedKwargs := *kwargs
						delete(slicedKwargs, "я")
						if err := CheckFunctionArguments(function, &slicedArgs, &slicedKwargs); err != nil {
							return nil, err
						}

						return function.Call(state, &slicedArgs, &slicedKwargs)
					},
					[]FunctionReturnType{
						{
							Type:       Any,
							IsNullable: true,
						},
					},
					true,
					nil,
					"", // TODO: add doc
				),
			},
			MakeLogicalOperators(Function),
			MakeCommonOperators(Function),
		)
	}

	return NewBuiltinClass(
		common.FunctionTypeName,
		nil,
		BuiltinPackage,
		initAttributes,
		"",  // TODO: add doc
		nil, // CAUTION: segfault may be thrown when using without nil check!
	)
}
