package models

import (
	"fmt"
	"regexp"
)

type Token struct {
	Type            TokenType
	Text            string
	Pos             int
	Row             int
	IsUnaryOperator bool
}

func (t *Token) String() string {
	return t.Text
}

const (
	SingleLineComment = iota
	MultiLineComment
	ImportStdDirective
	ImportDirective
	Arrow
	RealNumber
	IntegerNumber
	String
	Bool
	Semicolon
	Colon
	Space
	ExponentOp
	ModuloOp
	BitwiseLeftShiftOp
	BitwiseRightShiftOp
	BitwiseAndOp
	BitwiseXorOp
	BitwiseOrOp
	EqualsOp
	NotEqualsOp
	GreaterOrEqualsOp
	GreaterOp
	LessOrEqualsOp
	LessOp
	Assign
	Add
	Sub
	Mul
	Div
	LPar
	RPar
	If
	Else
	For
	AndOp
	OrOp
	NotOp
	BitwiseNotOp
	LCurlyBracket
	RCurlyBracket
	LSquareBracket
	RSquareBracket
	Comma
	FunctionDef
	Return
	Nil
	Name
	TripleDot
	QuestionMark
	AttrAccessOp
)

var tokenTypeNames = map[int]string{
	SingleLineComment:  "однорядковий коментар",
	MultiLineComment:   "багаторядковий коментар",
	ImportStdDirective: "директива підключення файлу стандартної бібліотеки",
	ImportDirective:    "директива підключення файлу",
	Arrow:              "стрілка",
	RealNumber:         "дійсне число",
	IntegerNumber:      "ціле число",
	String:             "рядок",
	Bool:               "логічний тип",
	Semicolon:          "крапка з комою",
	Colon:              "двокрапка",
	Space:              "пропуск",
	ExponentOp:         "оператор піднесення до степеня",
	ModuloOp:           "оператор остачі від ділення",
	EqualsOp:           "умова рівності",
	NotEqualsOp:        "умова нерівності",
	GreaterOrEqualsOp:  "умова 'більше або дорівнює'",
	GreaterOp:          "умова 'більше'",
	LessOrEqualsOp:     "умова 'менше або дорівнює'",
	LessOp:             "умова 'менше'",
	Assign:             "оператор присвоєння",
	Add:                "оператор суми",
	Sub:                "оператор різниці",
	Mul:                "оператор добутку",
	Div:                "оператор частки",
	LPar:               "відкриваюча дужка",
	RPar:               "закриваюча дужка",
	If:                 "якщо",
	Else:               "інакше",
	For:                "для",
	AndOp:              "оператор логічного 'і'",
	OrOp:               "оператор логічного 'або'",
	NotOp:              "оператор логічного заперечення",
	BitwiseNotOp:       "унірний оператор заперечення",
	LCurlyBracket:      "відкриваюча фігурна дужка",
	RCurlyBracket:      "закриваюча фігурна дужка",
	LSquareBracket:     "відкриваюча квадратна дужка",
	RSquareBracket:     "закриваюча квадратна дужка",
	Comma:              "кома",
	FunctionDef:        "визначення функції",
	Return:             "повернення значення",
	Nil:                "нуль",
	Name:               "назва",
	TripleDot:          "три крапки",
	QuestionMark:       "знак запитання",
	AttrAccessOp:       "оператор доступу до атрибута",
}

type TokenType struct {
	Name  int // iota
	Regex *regexp.Regexp
}

func (t *TokenType) String() string {
	return fmt.Sprintf("[%d | %s]", t.Name, t.Regex.String())
}

func (t TokenType) Description() string {
	if description, ok := tokenTypeNames[t.Name]; ok {
		return description
	}

	panic(fmt.Sprintf(
		"Unable to retrieve description for '%d' token, please add it to 'tokenTypeNames' map first",
		t.Name,
	))
}

const RawNameRegex = "[А-ЩЬЮЯҐЄІЇа-щьюяґєії_][А-ЩЬЮЯҐЄІЇа-щьюяґєії_0-9]*"

var NameRegex = regexp.MustCompile("(" + RawNameRegex + ")")

var TokenTypesList = map[int]TokenType{
	SingleLineComment: {
		Name:  SingleLineComment,
		Regex: regexp.MustCompile("^//[^\\n\\r]+?(?:\\*\\)|[\\n\\r])"),
	},
	MultiLineComment: {
		Name:  MultiLineComment,
		Regex: regexp.MustCompile("^(/\\*)(.|\\n)*?(\\*/)"),
	},
	ImportStdDirective: {
		Name: ImportStdDirective,
		Regex: regexp.MustCompile(
			"^@\\s*'\\s*([^.\\\\/`\\r\\n].*[^`\\r\\n])\\s*'",
		),
	},
	ImportDirective: {
		Name: ImportDirective,
		Regex: regexp.MustCompile(
			"^@\\s*\"\\s*([^\"\\r\\n].*[^\"\\r\\n])\\s*\"",
		),
	},
	Arrow: {
		Name:  Arrow,
		Regex: regexp.MustCompile("^->"),
	},
	RealNumber: {
		Name:  RealNumber,
		Regex: regexp.MustCompile("^[0-9]+(\\.[0-9]+)"),
	},
	IntegerNumber: {
		Name:  IntegerNumber,
		Regex: regexp.MustCompile("^\\d+"),
	},
	String: {
		Name:  String,
		Regex: regexp.MustCompile("^\"(?:[^\"\\\\]|\\\\.)*\""),
	},
	Bool: {
		Name:  Bool,
		Regex: regexp.MustCompile("^(істина|хиба)"),
	},
	Semicolon: {
		Name:  Semicolon,
		Regex: regexp.MustCompile("^;"),
	},
	Colon: {
		Name:  Colon,
		Regex: regexp.MustCompile("^:"),
	},
	Space: {
		Name:  Space,
		Regex: regexp.MustCompile("^[\\s\\n\\t\\r]"),
	},
	ExponentOp: {
		Name:  ExponentOp,
		Regex: regexp.MustCompile("^\\*\\*"),
	},
	ModuloOp: {
		Name:  ModuloOp,
		Regex: regexp.MustCompile("^%"),
	},
	BitwiseLeftShiftOp: {
		Name: BitwiseLeftShiftOp,
		Regex: regexp.MustCompile("^<<"),
	},
	BitwiseRightShiftOp: {
		Name: BitwiseRightShiftOp,
		Regex: regexp.MustCompile("^>>"),
	},
	BitwiseAndOp: {
		Name: BitwiseAndOp,
		Regex: regexp.MustCompile("^&"),
	},
	BitwiseXorOp: {
		Name: BitwiseXorOp,
		Regex: regexp.MustCompile("^\\^"),
	},
	BitwiseOrOp: {
		Name: BitwiseOrOp,
		Regex: regexp.MustCompile("^|"),
	},
	EqualsOp: {
		Name:  EqualsOp,
		Regex: regexp.MustCompile("^=="),
	},
	NotEqualsOp: {
		Name:  NotEqualsOp,
		Regex: regexp.MustCompile("^!="),
	},
	GreaterOrEqualsOp: {
		Name:  GreaterOrEqualsOp,
		Regex: regexp.MustCompile("^>="),
	},
	GreaterOp: {
		Name:  GreaterOp,
		Regex: regexp.MustCompile("^>"),
	},
	LessOrEqualsOp: {
		Name:  LessOrEqualsOp,
		Regex: regexp.MustCompile("^<="),
	},
	LessOp: {
		Name:  LessOp,
		Regex: regexp.MustCompile("^<"),
	},
	Assign: {
		Name:  Assign,
		Regex: regexp.MustCompile("^="),
	},
	Add: {
		Name:  Add,
		Regex: regexp.MustCompile("^\\+"),
	},
	Sub: {
		Name:  Sub,
		Regex: regexp.MustCompile("^-"),
	},
	Mul: {
		Name:  Mul,
		Regex: regexp.MustCompile("^\\*"),
	},
	Div: {
		Name:  Div,
		Regex: regexp.MustCompile("^/"),
	},
	LPar: {
		Name:  LPar,
		Regex: regexp.MustCompile("^\\("),
	},
	RPar: {
		Name:  RPar,
		Regex: regexp.MustCompile("^\\)"),
	},
	If: {
		Name:  If,
		Regex: regexp.MustCompile("^якщо"),
	},
	Else: {
		Name:  Else,
		Regex: regexp.MustCompile("^інакше"),
	},
	For: {
		Name:  For,
		Regex: regexp.MustCompile("^для"),
	},
	AndOp: {
		Name:  AndOp,
		Regex: regexp.MustCompile("^&&"),
	},
	OrOp: {
		Name:  OrOp,
		Regex: regexp.MustCompile("^\\|\\|"),
	},
	NotOp: {
		Name:  NotOp,
		Regex: regexp.MustCompile("^!"),
	},
	BitwiseNotOp: {
		Name: BitwiseNotOp,
		Regex: regexp.MustCompile("^~"),
	},
	LCurlyBracket: {
		Name:  LCurlyBracket,
		Regex: regexp.MustCompile("^{"),
	},
	RCurlyBracket: {
		Name:  RCurlyBracket,
		Regex: regexp.MustCompile("^}"),
	},
	LSquareBracket: {
		Name:  LSquareBracket,
		Regex: regexp.MustCompile("^\\["),
	},
	RSquareBracket: {
		Name:  RSquareBracket,
		Regex: regexp.MustCompile("^]"),
	},
	Comma: {
		Name:  Comma,
		Regex: regexp.MustCompile("^,"),
	},
	FunctionDef: {
		Name:  FunctionDef,
		Regex: regexp.MustCompile("^функція"),
	},
	Return: {
		Name:  Return,
		Regex: regexp.MustCompile("^повернути"),
	},
	Nil: {
		Name:  Nil,
		Regex: regexp.MustCompile("^нуль"),
	},
	Name: {
		Name:  Name,
		Regex: regexp.MustCompile("^" + RawNameRegex),
	},
	TripleDot: {
		Name:  TripleDot,
		Regex: regexp.MustCompile("^\\.{3}"),
	},
	QuestionMark: {
		Name:  QuestionMark,
		Regex: regexp.MustCompile("^\\?"),
	},
	AttrAccessOp: {
		Name:  AttrAccessOp,
		Regex: regexp.MustCompile("^\\."),
	},
}
