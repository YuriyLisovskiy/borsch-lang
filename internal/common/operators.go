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

	// logical
	AndOp
	OrOp

	// bitwise
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

	// other operators
	ConstructorOp
	CallOp

	// math: unary operators
	UnaryMinus
	UnaryPlus

	// logical: unary operators
	NotOp

	// bitwise: unary operators
	UnaryBitwiseNotOp

	// other: unary operators
	LengthOp
	BoolOp
	IntOp
	RealOp
	StringOp
	RepresentationOp
)

var opTypesToSignatures = map[OperatorHash]string{
	PowOp:               "**",
	ModuloOp:            "%",
	AddOp:               "+",
	SubOp:               "-",
	MulOp:               "*",
	DivOp:               "/",
	AndOp:               "&&",
	OrOp:                "||",
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

	ConstructorOp: "__конструктор__",
	CallOp:        "__виклик__",

	UnaryMinus:        "-",
	UnaryPlus:         "+",
	NotOp:             "!",
	UnaryBitwiseNotOp: "~",
	LengthOp:          "__довжина__",
	BoolOp:            "__логічне__",
	IntOp:             "__ціле__",
	RealOp:            "__дійсне__",
	StringOp:          "__рядок__",
	RepresentationOp:  "__представлення__",
}

var opSignaturesToHashes = map[string]OperatorHash{
	"**": PowOp,
	"%":  ModuloOp,
	"+":  AddOp,
	"-":  SubOp,
	"*":  MulOp,
	"/":  DivOp,
	"&&": AndOp,
	"||": OrOp,
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

	"__конструктор__": ConstructorOp,
	"__виклик__":      CallOp,

	"_-":                UnaryMinus,
	"_+":                UnaryPlus,
	"!":                 NotOp,
	"~":                 UnaryBitwiseNotOp,
	"__довжина__":       LengthOp,
	"__логічне__":       BoolOp,
	"__ціле__":          IntOp,
	"__дійсне__":        RealOp,
	"__рядок__":         StringOp,
	"__представлення__": RepresentationOp,
}

var opNames = []string{
	"**", // **
	"%",  // %
	"+",  // +
	"-",  // -
	"*",  // *
	"/",  // /
	"&&", // &&
	"||", // ||
	"<<", // <<
	">>", // >>
	"&",  // &
	"^",  // ^   TODO: підібрати відповідник до XOR
	"|",  // |
	"==", // ==
	"!=", // !=
	">",  // >
	">=", // >=
	"<",  // <
	"<=", // <=

	"__конструктор__",
	"__виклик__",

	"унарний -", // -
	"унарний +", // +
	"!",         // !
	"~",         // ~
	"__довжина__",
	"__логічне__",
	"__ціле__",
	"__дійсне__",
	"__рядок__",
	"__представлення__",
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

func (op OperatorHash) IsUnary() bool {
	return op >= UnaryMinus
}

func (op OperatorHash) IsBinary() bool {
	return op <= LessOrEqualsOp
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
