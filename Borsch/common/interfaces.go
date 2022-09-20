package common

import "github.com/alecthomas/participle/v2/lexer"

type Statement interface {
	String() string
	Position() lexer.Position
}
