package grammar

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/participle/v2"
)

type Parser struct {
	parser *participle.Parser
}

func NewParser() (*Parser, error) {
	parser, err := participle.Build(&Package{}, participle.UseLookahead(2))
	if err != nil {
		return nil, err
	}

	return &Parser{parser: parser}, nil
}

func (p *Parser) Parse(filename string, code string) (*Package, error) {
	packageAST := &Package{}
	err := p.parser.ParseString(filename, code, packageAST)
	if err != nil {
		switch parseError := err.(type) {
		case participle.UnexpectedTokenError:
			stacktrace := fmt.Sprintf(
				"  Файл \"%s\", рядок %d, позиція %d,\n    %s\n    %s\nСинтаксична помилка: %s",
				filename,
				parseError.Position().Line,
				parseError.Position().Column,
				parseError.Unexpected.Value,
				strings.Repeat(" ", utf8.RuneCountInString(parseError.Unexpected.Value))+"^",
				parseError.Message(),
			)
			return nil, errors.New(fmt.Sprintf("Відстеження (стек викликів):\n%s", stacktrace))
		default:
			return nil, err
		}
	}

	return packageAST, nil
}
