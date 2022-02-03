package std

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"

var ErrorClass *types.Class = nil

func Init() {
	ErrorClass = newErrorClass()

	ErrorClass.InitAttributes()
}
