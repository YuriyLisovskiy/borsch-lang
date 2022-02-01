package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
)

type FunctionArgument struct {
	// nil means any type
	Type       *Class
	Name       string
	IsVariadic bool
	IsNullable bool
}

func (fa *FunctionArgument) String() string {
	res := fa.Name
	if fa.IsVariadic {
		res += "..."
	}

	return res + fa.GetTypeName()
}

func (fa FunctionArgument) GetTypeName() string {
	res := ""
	if fa.Type != Nil {
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
	if r.Type != Nil {
		return r.Type.GetTypeName()
	}

	return common.AnyTypeName
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
	handler func(common.Context, *[]common.Type, *map[string]common.Type) (common.Type, error),
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
			typeName:    common.FunctionTypeName,
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

func (t FunctionInstance) String(common.Context) string {
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

func (t FunctionInstance) Representation(ctx common.Context) string {
	return t.String(ctx)
}

func (t FunctionInstance) AsBool(common.Context) bool {
	return true
}

func (t FunctionInstance) GetTypeName() string {
	return t.GetPrototype().GetTypeName()
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

	return t.GetPrototype().GetAttribute(name)
}

func (t FunctionInstance) GetPrototype() *Class {
	return Function
}

func (t *FunctionInstance) GetParent() common.Type {
	return t.package_
}

func (t *FunctionInstance) GetContext() common.Context {
	if t.package_ != nil {
		return t.package_.GetContext()
	}

	return nil
}

func newFunctionClass() *Class {
	initAttributes := func() map[string]common.Type {
		return mergeAttributes(
			map[string]common.Type{
				ops.CallOperatorName: NewFunctionInstance(
					ops.CallOperatorName,
					[]FunctionArgument{
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
					func(ctx common.Context, args *[]common.Type, kwargs *map[string]common.Type) (
						common.Type,
						error,
					) {
						function := (*args)[0].(*FunctionInstance)
						slicedArgs := (*args)[1:]
						slicedKwargs := *kwargs
						delete(slicedKwargs, "я")
						if err := CheckFunctionArguments(ctx, function, &slicedArgs, &slicedKwargs); err != nil {
							return nil, err
						}

						return function.Call(ctx, &slicedArgs, &slicedKwargs)
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
			makeLogicalOperators(Function),
			makeCommonOperators(Function),
		)
	}

	return NewBuiltinClass(
		common.FunctionTypeName,
		BuiltinPackage,
		initAttributes,
		"",  // TODO: add doc
		nil, // CAUTION: segfault may be thrown when using without nil check!
	)
}
