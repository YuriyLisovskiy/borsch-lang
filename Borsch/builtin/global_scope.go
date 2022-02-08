package builtin

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

var GlobalScope map[string]common.Value

func init() {
	types.Init()
	initFunctions()
	initClasses()

	GlobalScope = map[string]common.Value{

		// I/O
		PrintFunction.Name:     PrintFunction,
		PrintLineFunction.Name: PrintLineFunction,
		InputFunction.Name:     InputFunction,

		// Common
		PanicFunction.Name:     PanicFunction,
		EnvFunction.Name:       EnvFunction,
		AssertFunction.Name:    AssertFunction,
		CopyrightFunction.Name: CopyrightFunction,
		LicenceFunction.Name:   LicenceFunction,
		HelpFunction.Name:      HelpFunction,

		// System calls
		ExitFunction.Name: ExitFunction,

		// Conversion
		"дійсний":           types.Real,
		"логічний":          types.Bool,
		ImportFunction.Name: ImportFunction,
		"рядок":             types.String,
		"словник":           types.Dictionary,
		"список":            types.List,
		// "функція": types.Function,
		"цілий":     types.Integer,
		"довільний": types.Any,

		// Utilities
		"довжина":   LengthFunction,
		"додати":    AddToListFunction,
		"копіювати": DeepCopyFunction,
		"тип":       TypeFunction,

		// Classes
		ErrorClass.GetTypeName(): ErrorClass,
	}

	types.BuiltinPackage.SetAttributes(GlobalScope)
}
