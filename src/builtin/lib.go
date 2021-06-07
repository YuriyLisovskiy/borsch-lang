package builtin

import "github.com/YuriyLisovskiy/borsch/src/builtin/types"

var FunctionsList = map[string]func(...types.ValueType) (types.ValueType, error){
	// I/O
	"друк":  Print,
	"друкр": PrintLn,
	"ввід":  Input,

	// Conversion
	"цілий":    ToInteger,
	"дійсний":  ToReal,
	"рядок":    ToString,
	"логічний": ToBool,
	// TODO: список()

	// Common
	"паніка":     Panic,
	"середовище": GetEnv,
	"довжина":    Length,

	// Containers' manipulation
	// TODO: додати
}
