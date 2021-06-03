package parser

import (
	"errors"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/models"
)

func (p *Parser) readScope() ([]models.Token, error) {
	_, err := p.require(models.TokenTypesList[models.LCurlyBracket])
	if err != nil {
		return nil, err
	}

	fromPos := p.pos
	for p.pos < len(p.tokens) {
		if p.match(models.TokenTypesList[models.RCurlyBracket]) != nil {
			p.pos--
			break
		}

		p.pos++
	}

	_, err = p.require(models.TokenTypesList[models.RCurlyBracket])
	if err != nil {
		return nil, err
	}

	var result []models.Token
	if fromPos < p.pos {
		result = p.tokens[fromPos:p.pos-1]
	}

	return result, err
}

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

		ifNode := ast.NewIfSequenceNode(conditionNode, blockOfCode)
		for p.match(models.TokenTypesList[models.Else]) != nil {
			if p.match(models.TokenTypesList[models.If]) != nil {
				_, err := p.require(models.TokenTypesList[models.LPar])
				if err != nil {
					return nil, err
				}

				conditionNode, err = p.parseExpression()
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
