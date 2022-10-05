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
	initializeObjectMethod := types.MakeObjectClassInitializer(BuiltinPackage)
	types.ObjectClass.AddAttributes(
		map[string]types.Object{
			initializeObjectMethod.Name: initializeObjectMethod,
		},
	)

	types.ErrorClass.AddAttributes(types.MakeErrorClassMethods(BuiltinPackage))
	types.ErrorClass.Operators = types.MakeErrorClassOperators(BuiltinPackage)

	addMethod := methods.MakeAdd(BuiltinPackage)
	assertMethod := methods.MakeAssert(BuiltinPackage)
	lenMethod := methods.MakeLen(BuiltinPackage)
	printlnMethod := methods.MakePrintln(BuiltinPackage)

	GlobalScope = map[string]types.Object{
		types.ObjectClass.Name: types.ObjectClass,
		types.TypeClass.Name:   types.TypeClass,

		types.BoolClass.Name:   types.BoolClass,
		types.IntClass.Name:    types.IntClass,
		types.ListClass.Name:   types.ListClass,
		types.RealClass.Name:   types.RealClass,
		types.StringClass.Name: types.StringClass,
		types.TupleClass.Name:  types.TupleClass,

		types.ErrorClass.Name:                types.ErrorClass,
		types.RuntimeErrorClass.Name:         types.RuntimeErrorClass,
		types.TypeErrorClass.Name:            types.TypeErrorClass,
		types.AssertionErrorClass.Name:       types.AssertionErrorClass,
		types.ZeroDivisionErrorClass.Name:    types.ZeroDivisionErrorClass,
		types.IndexOutOfRangeErrorClass.Name: types.IndexOutOfRangeErrorClass,

		addMethod.Name:     addMethod,
		assertMethod.Name:  assertMethod,
		lenMethod.Name:     lenMethod,
		printlnMethod.Name: printlnMethod,

		types.ErrorClass.Name:     types.ErrorClass,
		types.TypeErrorClass.Name: types.TypeErrorClass,
	}

	types.Initialized = true
}
