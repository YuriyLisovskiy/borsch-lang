package parser

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/lang/ast"
	"github.com/YuriyLisovskiy/borsch/lang/models"
)

func (p *Parser) parseForLoop() (ast.ExpressionNode, error) {
	if forNode := p.match(models.TokenTypesList[models.For]); forNode != nil {
		_, err := p.require(models.TokenTypesList[models.LPar])
		if err != nil {
			return nil, err
		}

		indexVar, err := p.require(models.TokenTypesList[models.Name])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s змінної для порядкового номера або ключа", err.Error()))
		}

		if p.match(models.TokenTypesList[models.Comma]) != nil {
			itemVar, err := p.require(models.TokenTypesList[models.Name])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s змінної елемента послідовності або значення", err.Error()))
			}

			if indexVar.Text == itemVar.Text && indexVar.Text != "_" {
				return nil, errors.New("неможливо створити пару змінних з однаковими назвами")
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
