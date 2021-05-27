package src

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/models"
)

type Parser struct {
	tokens   []models.Token
	pos      int
	fileName string
}

func NewParser(fileName string, tokens []models.Token) (*Parser, error) {
	return &Parser{
		tokens:   tokens,
		pos:      0,
		fileName: fileName,
	}, nil
}

func (p *Parser) match(expected ...models.TokenType) *models.Token {
	if p.pos < len(p.tokens) {
		currentToken := p.tokens[p.pos]
		for _, typ := range expected {
			if typ.Name == currentToken.Type.Name {
				p.pos++
				return &currentToken
			}
		}
	}

	return nil
}

func (p *Parser) require(expected ...models.TokenType) (*models.Token, error) {
	token := p.match(expected...)
	if token == nil {
		return nil, errors.New(
			fmt.Sprintf("На позиції %d очікується %s", p.pos, models.TokenTypeNames[expected[0].Name]),
		)
	}

	return token, nil
}

func (p *Parser) parseVariableOrConstant() (ast.ExpressionNode, error) {
	number := p.match(models.TokenTypesList[models.Number])
	if number != nil {
		return ast.NewNumberNode(*number), nil
	}

	stringToken := p.match(models.TokenTypesList[models.String])
	if stringToken != nil {
		return ast.NewStringNode(*stringToken), nil
	}

	name := p.match(models.TokenTypesList[models.Name])
	if name != nil {
		if p.match(models.TokenTypesList[models.LPar]) != nil {
			p.pos--
			return nil, nil
		}

		return ast.NewVariableNode(*name), nil
	}

	return nil, errors.New(
		fmt.Sprintf("Очікується змінна або число на позиції %d", p.pos),
	)
}

func (p *Parser) parseParentheses() (ast.ExpressionNode, error) {
	if p.match(models.TokenTypesList[models.LPar]) != nil {
		node, err := p.parseFormula()
		if err != nil {
			return nil, err
		}

		_, err = p.require(models.TokenTypesList[models.RPar])
		if err != nil {
			return nil, err
		}

		return node, err
	}

	return p.parseExpression()
}

func (p *Parser) parseFormula() (ast.ExpressionNode, error) {
	leftNode, err := p.parseParentheses()
	if err != nil {
		return nil, err
	}

	operator := p.match(
		models.TokenTypesList[models.Add], models.TokenTypesList[models.Sub],
		models.TokenTypesList[models.Mul], models.TokenTypesList[models.Div],
	)
	for operator != nil {
		rightNode, err := p.parseParentheses()
		if err != nil {
			return nil, err
		}

		leftNode = ast.NewBinOperationNode(*operator, leftNode, rightNode)
		operator = p.match(
			models.TokenTypesList[models.Add], models.TokenTypesList[models.Sub],
			models.TokenTypesList[models.Mul], models.TokenTypesList[models.Div],
		)
	}

	return leftNode, nil
}

func (p *Parser) parseFunctionCall() (ast.ExpressionNode, error) {
	name := p.match(models.TokenTypesList[models.Name])
	if name != nil {
		lPar := p.match(models.TokenTypesList[models.LPar])
		if lPar != nil {
			var args []ast.ExpressionNode
			if p.match(models.TokenTypesList[models.RPar]) != nil {
				return ast.NewFunctionCallNode(*name, args), nil
			}

			for {
				argNode, err := p.parseFormula()
				if err != nil {
					return nil, err
				}

				args = append(args, argNode)
				if p.match(models.TokenTypesList[models.Comma]) == nil {
					_, err := p.require(models.TokenTypesList[models.RPar])
					if err != nil {
						return nil, err
					}

					break
				}
			}

			return ast.NewFunctionCallNode(*name, args), nil
		}

		return nil, errors.New(
			fmt.Sprintf("Очікується відкриваюча дужка на позиції %d", p.pos),
		)
	}

	return nil, errors.New(
		fmt.Sprintf("Очікується виклик функції на позиції %d", p.pos),
	)
}

func (p *Parser) parseExpression() (ast.ExpressionNode, error) {
	variableNode, err := p.parseVariableOrConstant()
	if err != nil {
		return nil, err
	}

	if variableNode != nil {
		return variableNode, nil
	}

	p.pos--
	funcCallNode, err := p.parseFunctionCall()
	if err != nil {
		return nil, err
	}

	return funcCallNode, nil
}

func (p *Parser) parseIncludeDirective() (ast.ExpressionNode, error) {
	includeDirective := p.match(models.TokenTypesList[models.IncludeDirective])
	if includeDirective != nil {
		return ast.NewIncludeDirectiveNode(*includeDirective), nil
	}

	return nil, nil
}

func (p *Parser) parseVariableAssignment() (ast.ExpressionNode, error) {
	name := p.match(models.TokenTypesList[models.Name])
	if name != nil {
		if p.match(models.TokenTypesList[models.LPar]) != nil {
			p.pos -= 2
			return nil, nil
		}

		variableNode := ast.NewVariableNode(*name)
		assignOperator := p.match(models.TokenTypesList[models.Assign])
		if assignOperator != nil {
			rightExpressionNode, err := p.parseFormula()
			if err != nil {
				return nil, err
			}

			binaryNode := ast.NewBinOperationNode(*assignOperator, variableNode, rightExpressionNode)
			return binaryNode, nil
		}
	}

	p.pos -= 1
	return nil, nil
}

func (p *Parser) parseRow() (ast.ExpressionNode, error) {
	includeDirectiveNode, err := p.parseIncludeDirective()
	if err != nil {
		return nil, err
	}

	if includeDirectiveNode != nil {
		return includeDirectiveNode, nil
	}

	assignmentNode, err := p.parseVariableAssignment()
	if err != nil {
		return nil, err
	}

	if assignmentNode != nil {
		_, err = p.require(models.TokenTypesList[models.Semicolon])
		if err != nil {
			return nil, err
		}

		return assignmentNode, nil
	}

	codeNode, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.require(models.TokenTypesList[models.Semicolon])
	if err != nil {
		return nil, err
	}

	return codeNode, nil
}

func (p *Parser) Parse() (*ast.AST, error) {
	asTree := &ast.AST{}
	rowCounter := 1
	for p.pos < len(p.tokens) {
		codeNode, err := p.parseRow()
		if err != nil {
			return nil, errors.New(fmt.Sprintf(
				"  Файл \"%s\", рядок %d\nСинтаксична помилка: %s",
				p.fileName, rowCounter, err.Error(),
			))
		}

		asTree.AddNode(codeNode)
		rowCounter++
	}

	return asTree, nil
}
