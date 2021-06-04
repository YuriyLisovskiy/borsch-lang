package parser

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/models"
)

func (p *Parser) parseForLoop() (ast.ExpressionNode, error) {
	if forNode := p.match(models.TokenTypesList[models.For]); forNode != nil {
		_, err := p.require(models.TokenTypesList[models.LPar])
		if err != nil {
			return nil, err
		}

		indexVar, err := p.require(models.TokenTypesList[models.Name])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s змінної для порядкового номера", err.Error()))
		}

		if p.match(models.TokenTypesList[models.Comma]) != nil {
			itemVar, err := p.require(models.TokenTypesList[models.Name])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s змінної елемента послідовності", err.Error()))
			}

			_, err = p.require(models.TokenTypesList[models.Colon])
			if err != nil {
				return nil, err
			}

			container, err := p.parseFormula()
			if err != nil {
				return nil, err
			}

			_, err = p.require(models.TokenTypesList[models.RPar])
			if err != nil {
				return nil, err
			}

			body, err := p.readScope()
			if err != nil {
				return nil, err
			}

			return ast.NewForEachNode(*forNode, *indexVar, *itemVar, container, body), nil
		} else {
			// TODO: зробити цикл 'для' без конкретного ітерабельного контейнера.
			return nil, errors.New("некоректний синтаксис")
		}
	}

	return nil, nil
}
