package builtin

var FunctionsList = map[string] func (...string) (string, error) {
	// I/O
	"друк": Print,
	"друкр": Println,

	// Math
	"сума": Sum,
	"лог10": Log10,
}
