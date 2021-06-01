package builtin

var FunctionsList = map[string] func (...ValueType) (ValueType, error) {
	// I/O
	"друк": Print,
	"друкр": PrintLn,
	"ввід": Input,

	// Cast
	"ціле": CastToInt,
	"дійсне": CastToReal,
	"рядок": CastToString,
	"логічне": CastToBool,

	// Common
	"паніка": Panic,
	"середовище": GetEnv,
	"довжина": Length,
	//"тип": Type,
}
