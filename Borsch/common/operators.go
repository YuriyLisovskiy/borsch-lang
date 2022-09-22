package common

import "fmt"

type OperatorHash int

const (
	// math
	PowOp OperatorHash = iota
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

var opTypesToSignatures = map[OperatorHash]string{
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

var opSignaturesToHashes = map[string]OperatorHash{
	"**": PowOp,
	"%":  ModuloOp,
	"+":  AddOp,
	"-":  SubOp,
	"*":  MulOp,
	"/":  DivOp,
	"_-": UnaryMinus,
	"_+": UnaryPlus,
	"&&": AndOp,
	"||": OrOp,
	"!":  NotOp,
	"~":  UnaryBitwiseNotOp,
	"<<": BitwiseLeftShiftOp,
	">>": BitwiseRightShiftOp,
	"&":  BitwiseAndOp,
	"^":  BitwiseXorOp,
	"|":  BitwiseOrOp,
	"==": EqualsOp,
	"!=": NotEqualsOp,
	">":  GreaterOp,
	">=": GreaterOrEqualsOp,
	"<":  LessOp,
	"<=": LessOrEqualsOp,
}

var opNames = []string{
	"**",        // **
	"%",         // %
	"+",         // +
	"-",         // -
	"*",         // *
	"/",         // /
	"унарний -", // -
	"унарний +", // +
	"&&",        // &&
	"||",        // ||
	"!",         // !
	"~",         // ~
	"<<",        // <<
	">>",        // >>
	"&",         // &
	"^",         // ^   TODO: підібрати відповідник до XOR
	"|",         // |
	"==",        // ==
	"!=",        // !=
	">",         // >
	">=",        // >=
	"<",         // <
	"<=",        // <=
}

func OperatorHashFromString(signature string) OperatorHash {
	if opHash, ok := opSignaturesToHashes[signature]; ok {
		return opHash
	}

	panic(
		fmt.Sprintf(
			"Unable to create hash for operator '%s', please add it to 'opSignaturesToHashes' map first",
			signature,
		),
	)
}

func (op OperatorHash) Sign() string {
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

func (op OperatorHash) Name() string {
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
