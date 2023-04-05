package interpreter

import (
	methods2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/methods"
	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

var (
	BuiltinPackage = types2.PackageNew(
		"builtin", nil, &ContextImpl{
			scopes:        []map[string]types2.Object{},
			parentContext: nil,
		},
	)

	GlobalScope map[string]types2.Object
)

func init() {
	initializeObjectMethod := types2.MakeObjectClassInitializer(BuiltinPackage)
	types2.ObjectClass.AddAttributes(
		map[string]types2.Object{
			initializeObjectMethod.Name: initializeObjectMethod,
		},
	)

	types2.ErrorClass.AddAttributes(types2.MakeErrorClassMethods(BuiltinPackage))
	types2.ErrorClass.Operators = types2.MakeErrorClassOperators(BuiltinPackage)

	addMethod := methods2.MakeAdd(BuiltinPackage)
	assertMethod := methods2.MakeAssert(BuiltinPackage)
	lenMethod := methods2.MakeLen(BuiltinPackage)
	printlnMethod := methods2.MakePrintln(BuiltinPackage)

	GlobalScope = map[string]types2.Object{
		types2.ObjectClass.Name: types2.ObjectClass,
		types2.TypeClass.Name:   types2.TypeClass,

		types2.BoolClass.Name:   types2.BoolClass,
		types2.IntClass.Name:    types2.IntClass,
		types2.ListClass.Name:   types2.ListClass,
		types2.RealClass.Name:   types2.RealClass,
		types2.StringClass.Name: types2.StringClass,
		types2.TupleClass.Name:  types2.TupleClass,

		types2.ErrorClass.Name:                types2.ErrorClass,
		types2.RuntimeErrorClass.Name:         types2.RuntimeErrorClass,
		types2.TypeErrorClass.Name:            types2.TypeErrorClass,
		types2.AssertionErrorClass.Name:       types2.AssertionErrorClass,
		types2.ZeroDivisionErrorClass.Name:    types2.ZeroDivisionErrorClass,
		types2.IndexOutOfRangeErrorClass.Name: types2.IndexOutOfRangeErrorClass,

		addMethod.Name:     addMethod,
		assertMethod.Name:  assertMethod,
		lenMethod.Name:     lenMethod,
		printlnMethod.Name: printlnMethod,

		types2.ErrorClass.Name:     types2.ErrorClass,
		types2.TypeErrorClass.Name: types2.TypeErrorClass,
	}

	types2.Initialized = true
}
