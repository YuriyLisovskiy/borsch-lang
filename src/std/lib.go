package std

var FunctionsList = map[string] func (...string) (string, error) {
	// I/O
	"стд::друк": Print,
	"стд::друкр": Println,

	// Math
	"стд::сума": Sum,
	"стд::лог10": Log10,
}
