package parser

import (
	"errors"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/models"
)

func (p *Parser) parseIfSequence() (ast.ExpressionNode, error) {
	if p.match(models.TokenTypesList[models.If]) != nil {
		_, err := p.require(models.TokenTypesList[models.LPar])
		if err != nil {
			return nil, err
		}

		conditionNode, err := p.parseFormula()
		if err != nil {
			return nil, err
		}

		_, err = p.require(models.TokenTypesList[models.RPar])
		if err != nil {
			return nil, err
		}

		blockOfCode, err := p.readScope()
		if err != nil {
			return nil, err
		}

		ifNode := ast.NewIfNode(conditionNode, blockOfCode)
		for p.match(models.TokenTypesList[models.Else]) != nil {
			if p.match(models.TokenTypesList[models.If]) != nil {
				_, err := p.require(models.TokenTypesList[models.LPar])
				if err != nil {
					return nil, err
				}

				conditionNode, err = p.parseFormula()
				if err != nil {
					return nil, err
				}

				_, err = p.require(models.TokenTypesList[models.RPar])
				if err != nil {
					return nil, err
				}

				blockOfCode, err = p.readScope()
				if err != nil {
					return nil, err
				}

				ifNode.Blocks = append(ifNode.Blocks, ast.NewConditionBlock(conditionNode, blockOfCode))
			} else {
				ifNode.ElseBlock, err = p.readScope()
				if err != nil {
					return nil, err
				}

				break
			}
		}

		return ifNode, nil
	}

	if p.match(models.TokenTypesList[models.Else]) != nil {
		return nil, errors.New("некоректний синтаксис")
	}

	return nil, nil
}
