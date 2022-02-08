package builtin

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"

var (
	PrintFunction     *types.FunctionInstance
	PrintLineFunction *types.FunctionInstance
	InputFunction     *types.FunctionInstance
	PanicFunction     *types.FunctionInstance
	EnvFunction       *types.FunctionInstance
	AssertFunction    *types.FunctionInstance
	CopyrightFunction *types.FunctionInstance
	LicenceFunction   *types.FunctionInstance
	HelpFunction      *types.FunctionInstance
	ExitFunction      *types.FunctionInstance
	ImportFunction    *types.FunctionInstance
	LengthFunction    *types.FunctionInstance
	AddToListFunction *types.FunctionInstance
	DeepCopyFunction  *types.FunctionInstance
	TypeFunction      *types.FunctionInstance
)

func initFunctions() {
	PrintFunction = makePrintFunction()
	PrintLineFunction = makePrintLnFunction()
	InputFunction = makeInputFunction()
	PanicFunction = makePanicFunction()
	EnvFunction = makeEnvFunction()
	AssertFunction = makeAssertFunction()
	CopyrightFunction = makeCopyrightFunction()
	LicenceFunction = makeLicenseFunction()
	HelpFunction = makeHelpFunction()
	ExitFunction = makeExitFunction()
	ImportFunction = makeImportFunction()
	LengthFunction = makeLengthFunction()
	AddToListFunction = makeAddToListFunction()
	DeepCopyFunction = makeDeepCopyFunction()
	TypeFunction = makeTypeFunction()
}
