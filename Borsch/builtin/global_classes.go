package builtin

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"

var ErrorClass *types.Class = nil

func initClasses() {
	// def
	ErrorClass = newErrorClass()

	// init
	types.InitClass(ErrorClass)
}
