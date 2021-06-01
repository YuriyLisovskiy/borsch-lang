package src

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strings"
	"unicode/utf8"
)

type Parser struct {
	tokens   []models.Token
	pos      int
	fileName string
}

func NewParser(fileName string, tokens []models.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		pos:      0,
		fileName: fileName,
	}
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
			fmt.Sprintf("очікується %s", expected[0].Description()),
		)
	}

	return token, nil
}

func (p *Parser) matchBinaryOperator() *models.Token {
	return p.match(
		models.TokenTypesList[models.Add], models.TokenTypesList[models.Sub],
		models.TokenTypesList[models.Mul], models.TokenTypesList[models.Div],
		models.TokenTypesList[models.AndOp], models.TokenTypesList[models.OrOp],
		models.TokenTypesList[models.EqualsOp], models.TokenTypesList[models.NotEqualsOp],
		models.TokenTypesList[models.GreaterOp], models.TokenTypesList[models.GreaterOrEqualsOp],
		models.TokenTypesList[models.LessOp], models.TokenTypesList[models.LessOrEqualsOp],
	)
}

func (p *Parser) parseVariableOrConstant() (ast.ExpressionNode, error) {
	number := p.match(models.TokenTypesList[models.RealNumber])
	if number != nil {
		return ast.NewRealTypeNode(*number), nil
	}

	number = p.match(models.TokenTypesList[models.IntegerNumber])
	if number != nil {
		return ast.NewIntegerTypeNode(*number), nil
	}

	stringToken := p.match(models.TokenTypesList[models.String])
	if stringToken != nil {
		return ast.NewStringTypeNode(*stringToken), nil
	}

	boolean := p.match(models.TokenTypesList[models.Bool])
	if boolean != nil {
		return ast.NewBoolTypeNode(*boolean), nil
	}

	name := p.match(models.TokenTypesList[models.Name])
	if name != nil {
		if p.match(models.TokenTypesList[models.LPar]) != nil {
			p.pos--
			return nil, nil
		}

		return ast.NewVariableNode(*name), nil
	}

	return nil, errors.New("очікується змінна або число")
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

	operator := p.matchBinaryOperator()
	for operator != nil {
		rightNode, err := p.parseParentheses()
		if err != nil {
			return nil, err
		}

		leftNode = ast.NewBinOperationNode(*operator, leftNode, rightNode)
		operator = p.matchBinaryOperator()
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

		return nil, errors.New("очікується відкриваюча дужка")
	}

	return nil, errors.New("очікується виклик функції")
}

func (p *Parser) parseExpression() (ast.ExpressionNode, error) {
	notOperator := p.match(models.TokenTypesList[models.NotOp])
	if notOperator != nil {
		operandNode, err := p.parseParentheses()
		if err != nil {
			return nil, err
		}

		return ast.NewUnaryOperationNode(*notOperator, operandNode), nil
	}

	variableNode, err := p.parseVariableOrConstant()
	if err != nil {
		return nil, err
	}

	if variableNode != nil {
		operator := p.matchBinaryOperator()
		if operator != nil {
			rightNode, err := p.parseParentheses()
			if err != nil {
				return nil, err
			}

			return ast.NewBinOperationNode(*operator, variableNode, rightNode), nil
		}

		//p.pos--
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
	includeDirective := p.match(models.TokenTypesList[models.IncludeStdDirective])
	if includeDirective != nil {
		return ast.NewIncludeDirectiveNode(*includeDirective, true), nil
	}

	includeDirective = p.match(models.TokenTypesList[models.IncludeDirective])
	if includeDirective != nil {
		return ast.NewIncludeDirectiveNode(*includeDirective, false), nil
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

	if p.pos < 0 {
		p.pos = 0
	}

	if includeDirectiveNode != nil {
		return includeDirectiveNode, nil
	}

	assignmentNode, err := p.parseVariableAssignment()
	if err != nil {
		return nil, err
	}

	if p.pos < 0 {
		p.pos = 0
	}

	if assignmentNode != nil {
		_, err = p.require(models.TokenTypesList[models.Semicolon])
		if err != nil {
			return nil, err
		}

		return assignmentNode, nil
	}

	if p.pos < 0 {
		p.pos = 0
	}

	codeNode, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.pos < 0 {
		p.pos = 0
	}

	_, err = p.require(models.TokenTypesList[models.Semicolon])
	if err != nil {
		return nil, err
	}

	return codeNode, nil
}

func (p *Parser) Parse() (*ast.AST, error) {
	asTree := &ast.AST{}
	for p.pos < len(p.tokens) {
		codeNode, err := p.parseRow()
		if err != nil {
			if p.pos == 0 {
				p.pos = 1
			}

			tokenString := p.tokens[p.pos-1].String()
			return nil, errors.New(fmt.Sprintf(
				"  Файл \"%s\", рядок %d\n    %s\n    %s\nСинтаксична помилка: %s",
				p.fileName, p.tokens[p.pos-1].Row,
				tokenString, strings.Repeat(" ", utf8.RuneCountInString(tokenString)) + "^",
				err.Error(),
			))
		}

		asTree.AddNode(codeNode)
	}

	return asTree, nil
}
