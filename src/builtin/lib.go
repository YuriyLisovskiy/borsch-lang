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
	"середовище": GetEnv,
	"паніка": Panic,
	//"тип": Type,
	//"довжина": Length,
}
