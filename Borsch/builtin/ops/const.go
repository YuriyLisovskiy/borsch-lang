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

func (op Operator) Description() string {
	if op >= 0 && int(op) < len(opTypeNames) {
		return opTypeNames[op]
	}

	panic(fmt.Sprintf(
		"Unable to retrieve description for operator '%d', please add it to 'opTypeNames' map first",
		op,
	))
}
