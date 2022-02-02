package common

import "fmt"

type Operator int

const (
	// math
	PowOp Operator = iota
	ModuloOp
	AddOp
	SubOp
	MulOp
	DivOp
	UnaryMinus
	UnaryPlus

	// logical
	AndOp
	OrOp
	NotOp

	// bitwise
	UnaryBitwiseNotOp
	BitwiseLeftShiftOp
	BitwiseRightShiftOp
	BitwiseAndOp
	BitwiseXorOp
	BitwiseOrOp

	// conditional
	EqualsOp
	NotEqualsOp
	GreaterOp
	GreaterOrEqualsOp
	LessOp
	LessOrEqualsOp
)

var opTypesToSignatures = map[Operator]string{
	PowOp:               "**",
	ModuloOp:            "%",
	AddOp:               "+",
	SubOp:               "-",
	MulOp:               "*",
	DivOp:               "/",
	UnaryMinus:          "-",
	UnaryPlus:           "+",
	AndOp:               "&&",
	OrOp:                "||",
	NotOp:               "!",
	UnaryBitwiseNotOp:   "~",
	BitwiseLeftShiftOp:  "<<",
	BitwiseRightShiftOp: ">>",
	BitwiseAndOp:        "&",
	BitwiseXorOp:        "^",
	BitwiseOrOp:         "|",
	EqualsOp:            "==",
	NotEqualsOp:         "!=",
	GreaterOp:           ">",
	GreaterOrEqualsOp:   ">=",
	LessOp:              "<",
	LessOrEqualsOp:      "<=",
}

var opNames = []string{
	"__оператор_степеня__",             // **
	"__оператор_ділення_за_модулем__",  // %
	"__оператор_суми__",                // +
	"__оператор_різниці__",             // -
	"__оператор_добутку__",             // *
	"__оператор_частки__",              // /
	"__оператор_мінус__",               // -
	"__оператор_плюс__",                // +
	"__оператор_і__",                   // &&
	"__оператор_або__",                 // ||
	"__оператор_не__",                  // !
	"__оператор_побітового_не__",       // ~
	"__оператор_зсуву_ліворуч__",       // <<
	"__оператор_зсуву_праворуч__",      // >>
	"__оператор_побітового_і__",        // &
	"__оператор_побітового_XOR__",      // ^   TODO: підібрати відповідник до XOR
	"__оператор_побітового_або__",      // |
	"__оператор_рівності__",            // ==
	"__оператор_нерівності__",          // !=
	"__оператор_більше__",              // >
	"__оператор_більше_або_дорівнює__", // >=
	"__оператор_менше__",               // <
	"__оператор_менше_або_дорівнює__",  // <=
}

func (op Operator) Sign() string {
	if op >= 0 && int(op) < len(opTypesToSignatures) {
		return opTypesToSignatures[op]
	}

	panic(
		fmt.Sprintf(
			"Unable to retrieve description for operator '%d', please add it to 'opTypesToSignatures' map first",
			op,
		),
	)
}

func (op Operator) Name() string {
	if op >= 0 && int(op) < len(opNames) {
		return opNames[op]
	}

	panic(
		fmt.Sprintf(
			"Unable to retrieve caption for operator '%d', please add it to 'opNames' map first",
			op,
		),
	)
}

func IsOperator(name string) bool {
	for _, current := range opNames {
		if current == name {
			return true
		}
	}

	return false
}
