package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/methods"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

var (
	BuiltinPackage = types.PackageNew(
		"builtin", nil, &ContextImpl{
			scopes:        []map[string]types.Object{},
			parentContext: nil,
		},
	)

	GlobalScope map[string]types.Object
)

func init() {
	types.ErrorClass.AddAttributes(types.MakeErrorClassAttributes(BuiltinPackage))
	// types.TypeErrorClass.AddAttributes(types.MakeErrorClassAttributes(BuiltinPackage))

	addMethod := methods.MakeAdd(BuiltinPackage)
	printlnMethod := methods.MakePrintln(BuiltinPackage)
	assertMethod := methods.MakeAssert(BuiltinPackage)

	GlobalScope = map[string]types.Object{
		types.BoolClass.Name:   types.BoolClass,
		types.IntClass.Name:    types.IntClass,
		types.ListClass.Name:   types.ListClass,
		types.RealClass.Name:   types.RealClass,
		types.StringClass.Name: types.StringClass,
		types.TupleClass.Name:  types.TupleClass,

		types.TypeClass.Name: types.TypeClass,

		types.ErrorClass.Name:                types.ErrorClass,
		types.RuntimeErrorClass.Name:         types.RuntimeErrorClass,
		types.TypeErrorClass.Name:            types.TypeErrorClass,
		types.AssertionErrorClass.Name:       types.AssertionErrorClass,
		types.ZeroDivisionErrorClass.Name:    types.ZeroDivisionErrorClass,
		types.IndexOutOfRangeErrorClass.Name: types.IndexOutOfRangeErrorClass,

		addMethod.Name:     addMethod,
		printlnMethod.Name: printlnMethod,
		assertMethod.Name:  assertMethod,

		types.ErrorClass.Name:     types.ErrorClass,
		types.TypeErrorClass.Name: types.TypeErrorClass,
	}

	types.Initialized = true
}
