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
	"список":   ToList,
	"словник":  ToDictionary,

	// Common
	"паніка":     Panic,
	"середовище": GetEnv,
	"підтвердити": Assert,

	// System calls
	"вихід": Exit,

	// Utilities
	"довжина": Length,
	"додати": AppendToList,
	"вилучити": RemoveFromDictionary,
}
