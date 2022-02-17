package builtin

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

var GlobalScope map[string]common.Object

func init() {
	types.Init()
	initFunctions()
	initClasses()

	GlobalScope = map[string]common.Object{

		// I/O
		PrintFunction.Name:     PrintFunction,
		PrintLineFunction.Name: PrintLineFunction,
		InputFunction.Name:     InputFunction,

		// Common
		// PanicFunction.Name:     PanicFunction,
		EnvFunction.Name:       EnvFunction,
		AssertFunction.Name:    AssertFunction,
		CopyrightFunction.Name: CopyrightFunction,
		LicenceFunction.Name:   LicenceFunction,
		HelpFunction.Name:      HelpFunction,

		// System calls
		ExitFunction.Name: ExitFunction,

		// Conversion
		"дійсний":           types.RealClass,
		"логічний":          types.BoolClass,
		ImportFunction.Name: ImportFunction,
		"рядок":             types.StringClass,
		"словник":           types.DictClass,
		"список":            types.ListClass,
		// "функція": types.FunctionClass,
		"цілий":     types.IntClass,
		"довільний": types.AnyClass,
		"нульовий":  types.NilClass,

		// Utilities
		"довжина":   LengthFunction,
		"додати":    AddToListFunction,
		"копіювати": DeepCopyFunction,
		"тип":       TypeFunction,

		// Classes
		ErrorClass.GetName(): ErrorClass,
	}

	types.BuiltinPackage.SetAttributes(GlobalScope)
}
