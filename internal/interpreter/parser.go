package interpreter

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type ParserImpl struct {
	parser *participle.Parser
}

func NewParser() (*ParserImpl, error) {
	parser, err := participle.Build(
		&Package{},
		participle.UseLookahead(2),
		participle.Unquote("String", "RawString", "Char"),
		participle.Map(identMapper, "Ident"),
	)
	if err != nil {
		return nil, err
	}

	return &ParserImpl{parser: parser}, nil
}

func (p *ParserImpl) Parse(filename string, code string) (Evaluatable, error) {
	ast := &Package{}
	err := p.parser.ParseString(filename, code, ast)
	if err != nil {
		return nil, err
	}

	return ast, nil
}

var runes = map[rune]rune{
	'a': 'а',
	'c': 'с',
	'e': 'е',
	'i': 'і',
	'o': 'о',
	'p': 'р',
	'x': 'х',
	'y': 'у',
}

func identMapper(token lexer.Token) (lexer.Token, error) {
	value := []rune(token.Value)
	for i := range value {
		value[i] = getFixedRuneOrDefault(value[i])
	}

	token.Value = string(value)
	return token, nil
}

func getFixedRuneOrDefault(r rune) rune {
	if _, ok := runes[r]; ok {
		return runes[r]
	}

	return r
}
