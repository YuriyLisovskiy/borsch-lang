package builtin

import (
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
)

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
	// TODO: словник()

	// Common
	"паніка":     Panic,
	"середовище": GetEnv,

	// System calls
	"вихід": Exit,

	// Containers utilities
	"довжина": Length,

	// List utilities
	"додати": AppendToList,

	// Dictionary utilities
	"вилучити": RemoveFromDictionary,
}
