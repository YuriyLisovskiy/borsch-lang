package models

import (
	"fmt"
	"regexp"
)

const (
	SingleLineComment = iota
	MultiLineComment
	RealNumber
	IntegerNumber
	String
	Name
	Semicolon
	Space
	Assign
	Add
	Sub
	Mul
	Div
	LPar
	RPar
	LAngleBracket
	RAngleBracket
	Comma
	IncludeDirective
)

var TokenTypeNames = []string{
	"однорядковий коментар",
	"багаторядковий коментар",
	"дійсне число",
	"ціле число",
	"рядок",
	"назва",
	"крапка з комою",
	"пропуск",
	"оператор присвоєння",
	"оператор суми",
	"оператор різниці",
	"оператор добутку",
	"оператор частки",
	"відкриваюча дужка",
	"закриваюча дужка",
	"відкриваюча кутова дужка",
	"закриваюча кутова дужка",
	"кома",
	"директива підключення файлу",
}

type TokenType struct {
	Name  int // iota
	Regex *regexp.Regexp
}

func (tt *TokenType) String() string {
	return fmt.Sprintf("[%d | %s]", tt.Name, tt.Regex.String())
}

const nameRegex = "[А-ЩЬЮЯҐЄІЇа-щьюяґєії_][А-ЩЬЮЯҐЄІЇа-щьюяґєії_0-9]*"

var TokenTypesList = map[int]TokenType{
	SingleLineComment: {
		Name:  SingleLineComment,
		//Regex: regexp.MustCompile("^//[^\\n\\r]*.*[^\\n\\r]*"),
		Regex: regexp.MustCompile("^//[^\\n\\r]+?(?:\\*\\)|[\\n\\r])"),
	},
	MultiLineComment: {
		Name:  MultiLineComment,
		//Regex: regexp.MustCompile("^//[^\\n\\r]*.*[^\\n\\r]*"),
		Regex: regexp.MustCompile("^(/\\*)(.|\\n)*?(\\*/)"),
	},
	RealNumber: {
		Name:  RealNumber,
		Regex: regexp.MustCompile("^[0-9]+(\\.[0-9]+)"),
	},
	IntegerNumber: {
		Name:  IntegerNumber,
		//Regex: regexp.MustCompile("^[0-9]+([^.][0-9]+)?"),
		//Regex: regexp.MustCompile("^\\d+[^\\Df]?"),
		Regex: regexp.MustCompile("^\\d+"),
	},
	String: {
		Name:  String,
		Regex: regexp.MustCompile("^\"(?:[^\"\\\\]|\\\\.)*\""),
	},
	Name: {
		Name:  Name,
		Regex: regexp.MustCompile("^" + nameRegex),
	},
	Semicolon: {
		Name:  Semicolon,
		Regex: regexp.MustCompile("^;"),
	},
	Space: {
		Name:  Space,
		Regex: regexp.MustCompile("^[\\s\\n\\t\\r]"),
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
	LAngleBracket: {
		Name:  LAngleBracket,
		Regex: regexp.MustCompile("^<"),
	},
	RAngleBracket: {
		Name:  RAngleBracket,
		Regex: regexp.MustCompile("^>"),
	},
	Comma: {
		Name:  Comma,
		Regex: regexp.MustCompile("^,"),
	},
	IncludeDirective: {
		Name:  IncludeDirective,
		Regex: regexp.MustCompile(
			//"^@\\s*<\\s*([^<\\s\\r\\n].*[^>\\s\\r\\n])\\s*>\\sяк\\s(" + nameRegex + ")",
			"^@\\s*<\\s*([^<\\s\\r\\n].*[^>\\s\\r\\n])\\s*>",
		),
	},
}
