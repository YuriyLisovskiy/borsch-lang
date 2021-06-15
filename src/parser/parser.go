package parser

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
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

func (p *Parser) checkForKeyword(name string) error {
	if _, ok := builtin.RegisteredIdentifiers[name]; ok {
		return errors.New(fmt.Sprintf(
			"неможливо використати ідентифікатор '%s', осткільки він є вбудованим",
			name,
		))
	}

	return nil
}

func (p *Parser) parseVariableOrConstant() (ast.ExpressionNode, error) {
	if number := p.match(models.TokenTypesList[models.RealNumber]); number != nil {
		return ast.NewRealTypeNode(*number), nil
	}

	if number := p.match(models.TokenTypesList[models.IntegerNumber]); number != nil {
		return ast.NewIntegerTypeNode(*number), nil
	}

	if stringToken := p.match(models.TokenTypesList[models.String]); stringToken != nil {
		return ast.NewStringTypeNode(*stringToken), nil
	}

	if boolean := p.match(models.TokenTypesList[models.Bool]); boolean != nil {
		return ast.NewBoolTypeNode(*boolean), nil
	}

	if listStart := p.match(models.TokenTypesList[models.LSquareBracket]); listStart != nil {
		var values []ast.ExpressionNode
		if p.match(models.TokenTypesList[models.RSquareBracket]) != nil {
			return ast.NewListTypeNode(*listStart, values), nil
		}

		for {
			valueNode, err := p.parseFormula()
			if err != nil {
				return nil, err
			}

			values = append(values, valueNode)
			if p.match(models.TokenTypesList[models.Comma]) == nil {
				_, err := p.require(models.TokenTypesList[models.RSquareBracket])
				if err != nil {
					return nil, err
				}

				break
			}
		}

		return ast.NewListTypeNode(*listStart, values), nil
	}

	if dictStart := p.match(models.TokenTypesList[models.LCurlyBracket]); dictStart != nil {
		if p.match(models.TokenTypesList[models.RCurlyBracket]) != nil {
			return ast.NewDictionaryTypeNode(*dictStart), nil
		}

		dict := ast.NewDictionaryTypeNode(*dictStart)
		for {
			keyNode, err := p.parseFormula()
			if err != nil {
				return nil, err
			}

			_, err = p.require(models.TokenTypesList[models.Colon])
			if err != nil {
				return nil, err
			}

			valueNode, err := p.parseFormula()
			if err != nil {
				return nil, err
			}

			dict.Map[keyNode] = valueNode
			if p.match(models.TokenTypesList[models.Comma]) == nil {
				_, err := p.require(models.TokenTypesList[models.RCurlyBracket])
				if err != nil {
					return nil, err
				}

				break
			}
		}

		return dict, nil
	}

	if name := p.match(models.TokenTypesList[models.Name]); name != nil {
		if p.match(models.TokenTypesList[models.LPar]) != nil {
			p.pos--
			return nil, nil
		}

		err := p.checkForKeyword(name.Text)
		if err != nil {
			return nil, err
		}

		variable := ast.NewVariableNode(*name)

		// TODO: implement attr set access!
		if dot := p.match(models.TokenTypesList[models.AttrAccessOp]); dot != nil {
			return p.parseGetAttr(variable)
		}

		return variable, nil
	}

	return nil, errors.New("очікується змінна або вираз")
}

func (p *Parser) parseRandomAccessOperation(name *models.Token, expr ast.ExpressionNode) (ast.ExpressionNode, error) {
	if lSquareBracket := p.match(models.TokenTypesList[models.LSquareBracket]); lSquareBracket != nil {
		indexNode, err := p.parseFormula()
		if err != nil {
			return nil, err
		}

		token, err := p.require(
			models.TokenTypesList[models.RSquareBracket], models.TokenTypesList[models.Colon],
		)
		if err != nil {
			return nil, err
		}

		if token.Type.Name == models.RSquareBracket {
			if name != nil {
				return ast.NewRandomAccessSetOperationNode(*name, indexNode, lSquareBracket.Row), nil
			} else if expr != nil {
				return ast.NewRandomAccessGetOperationNode(expr, indexNode, lSquareBracket.Row), nil
			}

			panic(errors.New("unknown operation got"))
		}

		if token.Type.Name != models.Colon {
			panic(errors.New("got invalid token"))
		}

		rIndexNode, err := p.parseFormula()
		if err != nil {
			return nil, err
		}

		_, err = p.require(models.TokenTypesList[models.RSquareBracket])
		if err != nil {
			return nil, err
		}

		return ast.NewListSlicingNode(expr, indexNode, rIndexNode, lSquareBracket.Row), nil
	}

	return nil, nil
}

func (p *Parser) parseFunctionCall(parent ast.ExpressionNode) (ast.ExpressionNode, error) {
	name := p.match(models.TokenTypesList[models.Name])
	if name != nil {
		lPar := p.match(models.TokenTypesList[models.LPar])
		if lPar != nil {
			var args []ast.ExpressionNode
			if p.match(models.TokenTypesList[models.RPar]) != nil {
				return ast.NewFunctionCallNode(*name, parent, args), nil
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

			result := ast.NewFunctionCallNode(*name, parent, args)

			// TODO: implement attr set access!
			if dot := p.match(models.TokenTypesList[models.AttrAccessOp]); dot != nil {
				return p.parseGetAttr(result)
			}

			return result, nil
		}

		return nil, errors.New("очікується відкриваюча дужка")
	}

	return nil, errors.New("очікується виклик функції")
}

func (p *Parser) parseExpression() (ast.ExpressionNode, error) {
	variableNode, err := p.parseVariableOrConstant()
	if err != nil {
		return nil, err
	}

	if variableNode != nil {
		randomAccessOp, err := p.parseRandomAccessOperation(nil, variableNode)
		if err != nil {
			return nil, err
		}

		if randomAccessOp != nil {
			variableNode = randomAccessOp
		}

		return variableNode, nil
	}

	p.pos--
	funcCallNode, err := p.parseFunctionCall(nil)
	if err != nil {
		return nil, err
	}

	randomAccessOp, err := p.parseRandomAccessOperation(nil, funcCallNode)
	if err != nil {
		return nil, err
	}

	if randomAccessOp != nil {
		funcCallNode = randomAccessOp
	}

	return funcCallNode, nil
}

func (p *Parser) parseIncludeDirective() (ast.ExpressionNode, error) {
	isStd := false
	includeDirective := p.match(models.TokenTypesList[models.IncludeStdDirective])
	if includeDirective != nil {
		isStd = true
	} else {
		includeDirective = p.match(models.TokenTypesList[models.IncludeDirective])
	}

	if includeDirective != nil {
		arrow := p.match(models.TokenTypesList[models.Arrow])
		name := ""
		if arrow != nil {
			token, err := p.require(models.TokenTypesList[models.Name])
			if err != nil {
				return nil, err
			}

			name = token.Text
		}

		return ast.NewIncludeDirectiveNode(*includeDirective, name, isStd), nil
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

		err := p.checkForKeyword(name.Text)
		if err != nil {
			return nil, err
		}

		variableNode, err := p.parseRandomAccessOperation(name, nil)
		if err != nil {
			return nil, err
		}

		if variableNode == nil {
			variableNode = ast.NewVariableNode(*name)
		}

		assignOperator := p.match(models.TokenTypesList[models.Assign])
		if assignOperator != nil {
			rightExpressionNode, err := p.parseFormula()
			if err != nil {
				return nil, err
			}

			binaryNode := ast.NewBinOperationNode(*assignOperator, variableNode, rightExpressionNode)
			return binaryNode, nil
		}

		p.pos--
	}

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

	ifNode, err := p.parseIfSequence()
	if err != nil {
		return nil, err
	}

	if p.pos < 0 {
		p.pos = 0
	}

	if ifNode != nil {
		return ifNode, nil
	}

	forNode, err := p.parseForLoop()
	if err != nil {
		return nil, err
	}

	if p.pos < 0 {
		p.pos = 0
	}

	if forNode != nil {
		return forNode, nil
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

	//codeNode, err := p.parseExpression()
	codeNode, err := p.parseFormula()
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

func (p *Parser) skipSemicolons() {
	for p.match(models.TokenTypesList[models.Semicolon]) != nil {
	}
}

func (p *Parser) Parse() (*ast.AST, error) {
	asTree := ast.NewAST()
	for p.pos < len(p.tokens) {
		p.skipSemicolons()
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
		p.skipSemicolons()
	}

	return asTree, nil
}
