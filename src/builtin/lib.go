package builtin

var FunctionsList = map[string] func (...string) (string, error) {
	// I/O
	"друк": Print,
	"друкр": PrintLn,

	// Math
	"сума": Sum,
	"лог10": Log10,

	// OS
	"середовище": GetEnv,
}
