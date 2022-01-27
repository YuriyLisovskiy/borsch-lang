package ops

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

var opTypeNames = map[Operator]string{
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

var opTypeCaptions = map[Operator]string{
	PowOp:               "__оператор_степеня__",             // **
	ModuloOp:            "__оператор_ділення_за_модулем__",  // %
	AddOp:               "__оператор_суми__",                // +
	SubOp:               "__оператор_різниці__",             // -
	MulOp:               "__оператор_добутку__",             // *
	DivOp:               "__оператор_частки__",              // /
	UnaryMinus:          "__оператор_мінус__",               // -
	UnaryPlus:           "__оператор_плюс__",                // +
	AndOp:               "__оператор_і__",                   // &&
	OrOp:                "__оператор_або__",                 // ||
	NotOp:               "__оператор_не__",                  // !
	UnaryBitwiseNotOp:   "__оператор_побітового_не__",       // ~
	BitwiseLeftShiftOp:  "__оператор_зсуву_ліворуч__",       // <<
	BitwiseRightShiftOp: "__оператор_зсуву_праворуч__",      // >>
	BitwiseAndOp:        "__оператор_побітового_і__",        // &
	BitwiseXorOp:        "__оператор_побітового_XOR__",      // ^   TODO: підібрати відповідник до XOR
	BitwiseOrOp:         "__оператор_побітового_або__",      // |
	EqualsOp:            "__оператор_рівності__",            // ==
	NotEqualsOp:         "__оператор_нерівності__",          // !=
	GreaterOp:           "__оператор_більше__",              // >
	GreaterOrEqualsOp:   "__оператор_більше_або_дорівнює__", // >=
	LessOp:              "__оператор_менше__",               // <
	LessOrEqualsOp:      "__оператор_менше_або_дорівнює__",  // <=
}

func (op Operator) Sign() string {
	if op >= 0 && int(op) < len(opTypeNames) {
		return opTypeNames[op]
	}

	panic(
		fmt.Sprintf(
			"Unable to retrieve description for operator '%d', please add it to 'opTypeNames' map first",
			op,
		),
	)
}

func (op Operator) Name() string {
	if op >= 0 && int(op) < len(opTypeCaptions) {
		return opTypeCaptions[op]
	}

	panic(
		fmt.Sprintf(
			"Unable to retrieve caption for operator '%d', please add it to 'opTypeCaptions' map first",
			op,
		),
	)
}
