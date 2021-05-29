package builtin

var FunctionsList = map[string] func (...ValueType) (ValueType, error) {
	// I/O
	"друк": Print,
	"друкр": PrintLn,

	// TODO: remove from built-ins, it will be implemented in std lib.
	// Math
	"лог10": Log10,

	// OS
	"середовище": GetEnv,

	// Cast
	"ціле": CastToInt,
	"дійсне": CastToReal,
	"рядок": CastToString,
}
