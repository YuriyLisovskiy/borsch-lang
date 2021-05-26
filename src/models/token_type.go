package models

import (
	"fmt"
	"regexp"
)

const (
	Number = iota
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
	FilePath
	SingleLineComment
)

var TokenTypeNames = []string{
	"Number", "Name", "Semicolon", "Space", "Assign",
	"Add", "Sub", "Mul", "Div",
	"LPar", "RPar", "LAngleBracket", "RAngleBracket",
	"Comma", "IncludeDirective", "FilePath", "SingleLineComment",
}

type TokenType struct {
	Name  int // iota
	Regex *regexp.Regexp
}

func (tt *TokenType) String() string {
	return fmt.Sprintf("[%d | %s]", tt.Name, tt.Regex.String())
}

var TokenTypesList = map[int]TokenType{
	Number: {
		Name:  Number,
		Regex: regexp.MustCompile("^[0-9]+(\\.[0-9]+)?"),
	},
	Name: {
		Name:  Name,
		Regex: regexp.MustCompile("^(стд::)?[А-ЩЬЮЯҐЄІЇа-щьюяґєії_][А-ЩЬЮЯҐЄІЇа-щьюяґєії_0-9]*"),
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
		Regex: regexp.MustCompile("^#підключити\\s*<\\s*([^<\\s\\r\\n].*[^>\\s\\r\\n])\\s*>"),
	},
	SingleLineComment: {
		Name:  SingleLineComment,
		Regex: regexp.MustCompile("^~[^\\n\\r]*.*[^\\n\\r]*"),
	},
}
