package builtin

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
)

var BuiltinScope map[string]common.Type

func init() {
	BuiltinScope = map[string]common.Type{

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
	}

	types.BuiltinPackage.Attributes = BuiltinScope
}
