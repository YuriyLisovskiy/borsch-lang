package parser

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
)

func (p *Parser) parseClassDefinition() (ast.ExpressionNode, error) {
	if p.match(models.TokenTypesList[models.ClassDef]) != nil {
		name, err := p.require(models.TokenTypesList[models.Name])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s класу", err.Error()))
		}

		// TODO: parse inherited classes here

		_, err = p.require(models.TokenTypesList[models.LCurlyBracket])
		if err != nil {
			return nil, err
		}

		// TODO: parse class scope

		_, err = p.require(models.TokenTypesList[models.RCurlyBracket])
		if err != nil {
			return nil, err
		}

		// TODO: parse and set doc
		classNode := ast.NewClassDefNode(*name, models.Token{})
		return classNode, nil
	}

	return nil, nil
}
