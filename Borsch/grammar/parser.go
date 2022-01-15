package grammar

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
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

func (p *ParserImpl) NewContext(packageFilename string, parentPackage common.Type) common.Context {
	parentPackageName := ""
	if parentPackage != nil {
		parentPackageName = parentPackage.(*types.PackageInstance).Name
	}
	return &ContextImpl{
		package_: types.NewPackageInstance(false, packageFilename, parentPackageName, map[string]common.Type{}),
	}
}
