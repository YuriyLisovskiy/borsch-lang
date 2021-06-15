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
	IncludeStdDirective
	IncludeDirective
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
	LCurlyBracket
	RCurlyBracket
	LSquareBracket
	RSquareBracket
	Comma
	Name
	AttrAccessOp
)

var tokenTypeNames = map[int]string{
	SingleLineComment:   "однорядковий коментар",
	MultiLineComment:    "багаторядковий коментар",
	IncludeStdDirective: "директива підключення файлу стандартної бібліотеки",
	IncludeDirective:    "директива підключення файлу",
	Arrow:               "стрілка",
	RealNumber:          "дійсне число",
	IntegerNumber:       "ціле число",
	String:              "рядок",
	Bool:                "логічний тип",
	Semicolon:           "крапка з комою",
	Colon:               "двокрапка",
	Space:               "пропуск",
	ExponentOp:          "оператор піднесення до степеня",
	ModuloOp:            "оператор остачі від ділення",
	EqualsOp:            "умова рівності",
	NotEqualsOp:         "умова нерівності",
	GreaterOrEqualsOp:   "умова 'більше або дорівнює'",
	GreaterOp:           "умова 'більше'",
	LessOrEqualsOp:      "умова 'менше або дорівнює'",
	LessOp:              "умова 'менше'",
	Assign:              "оператор присвоєння",
	Add:                 "оператор суми",
	Sub:                 "оператор різниці",
	Mul:                 "оператор добутку",
	Div:                 "оператор частки",
	LPar:                "відкриваюча дужка",
	RPar:                "закриваюча дужка",
	If:                  "якщо",
	Else:                "інакше",
	For:                 "для",
	AndOp:               "оператор логічного 'і'",
	OrOp:                "оператор логічного 'або'",
	NotOp:               "оператор логічного заперечення",
	LCurlyBracket:       "відкриваюча фігурна дужка",
	RCurlyBracket:       "закриваюча фігурна дужка",
	LSquareBracket:      "відкриваюча квадратна дужка",
	RSquareBracket:      "закриваюча квадратна дужка",
	Comma:               "кома",
	Name:                "назва",
	AttrAccessOp:        "оператор доступу до атрибута",
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
	IncludeStdDirective: {
		Name: IncludeStdDirective,
		Regex: regexp.MustCompile(
			"^@\\s*<\\s*([^.\\\\/<\\r\\n].*[^>\\r\\n])\\s*>",
		),
	},
	IncludeDirective: {
		Name: IncludeDirective,
		Regex: regexp.MustCompile(
			"^@\\s*\"\\s*([^\"\\r\\n].*[^\"\\r\\n])\\s*\"",
		),
	},
	Arrow: {
		Name: Arrow,
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
	Name: {
		Name:  Name,
		Regex: regexp.MustCompile("^" + RawNameRegex),
	},
	AttrAccessOp: {
		Name:  AttrAccessOp,
		Regex: regexp.MustCompile("^\\."),
	},
}
