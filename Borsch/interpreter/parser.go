package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
	"github.com/alecthomas/participle/v2"
)

var ParserInstance common.Parser

func init() {
	var err error
	ParserInstance, err = NewParser()
	if err != nil {
		panic(err)
	}
}

type ParserImpl struct {
	parser *participle.Parser
}

func NewParser() (*ParserImpl, error) {
	parser, err := participle.Build(
		&Package{},
		participle.UseLookahead(2),
		participle.Unquote("String", "Char"),
	)

	if err != nil {
		return nil, err
	}

	return &ParserImpl{parser: parser}, nil
}

func (p *ParserImpl) Parse(filename string, code string) (common.Evaluatable, error) {
	packageAST := &Package{}
	err := p.parser.ParseString(filename, code, packageAST)
	if err != nil {
		switch parseError := err.(type) {
		case participle.UnexpectedTokenError:
			err := utilities.ParseError(parseError.Position(), parseError.Unexpected.Value, parseError.Message())
			return nil, errors.New(fmt.Sprintf("Відстеження (стек викликів):\n%s", err))
		default:
			return nil, err
		}
	}

	return packageAST, nil
}
