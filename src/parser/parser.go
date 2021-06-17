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

func (p *Parser) current() *models.Token {
	if p.pos < len(p.tokens) {
		pos := p.pos - 1
		if pos < 0 {
			pos = 0
		}

		return &p.tokens[pos]
	}

	return nil
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

func (p *Parser) parseVariableOrConstant() (ast.ExpressionNode, *models.Token, error) {
	if number := p.match(models.TokenTypesList[models.RealNumber]); number != nil {
		return ast.NewRealTypeNode(*number), nil, nil
	}

	if number := p.match(models.TokenTypesList[models.IntegerNumber]); number != nil {
		return ast.NewIntegerTypeNode(*number), nil, nil
	}

	if stringToken := p.match(models.TokenTypesList[models.String]); stringToken != nil {
		return ast.NewStringTypeNode(*stringToken), nil, nil
	}

	if boolean := p.match(models.TokenTypesList[models.Bool]); boolean != nil {
		return ast.NewBoolTypeNode(*boolean), nil, nil
	}

	if listStart := p.match(models.TokenTypesList[models.LSquareBracket]); listStart != nil {
		var values []ast.ExpressionNode
		if p.match(models.TokenTypesList[models.RSquareBracket]) != nil {
			return ast.NewListTypeNode(*listStart, values), nil, nil
		}

		for {
			valueNode, err := p.parseFormula()
			if err != nil {
				return nil, nil, err
			}

			values = append(values, valueNode)
			if p.match(models.TokenTypesList[models.Comma]) == nil {
				_, err := p.require(models.TokenTypesList[models.RSquareBracket])
				if err != nil {
					return nil, nil, err
				}

				break
			}
		}

		return ast.NewListTypeNode(*listStart, values), nil, nil
	}

	if dictStart := p.match(models.TokenTypesList[models.LCurlyBracket]); dictStart != nil {
		if p.match(models.TokenTypesList[models.RCurlyBracket]) != nil {
			return ast.NewDictionaryTypeNode(*dictStart), nil, nil
		}

		dict := ast.NewDictionaryTypeNode(*dictStart)
		for {
			keyNode, err := p.parseFormula()
			if err != nil {
				return nil, nil, err
			}

			_, err = p.require(models.TokenTypesList[models.Colon])
			if err != nil {
				return nil, nil, err
			}

			valueNode, err := p.parseFormula()
			if err != nil {
				return nil, nil, err
			}

			dict.Map[keyNode] = valueNode
			if p.match(models.TokenTypesList[models.Comma]) == nil {
				_, err := p.require(models.TokenTypesList[models.RCurlyBracket])
				if err != nil {
					return nil, nil, err
				}

				break
			}
		}

		return dict, nil, nil
	}

	if name := p.match(models.TokenTypesList[models.Name]); name != nil {
		if p.match(models.TokenTypesList[models.LPar]) != nil {
			return nil, name, nil
		}

		err := p.checkForKeyword(name.Text)
		if err != nil {
			return nil, nil, err
		}

		var variable ast.ExpressionNode = ast.NewVariableNode(*name)
		randomAccessOp, err := p.parseRandomAccessOperation(name, variable)
		if err != nil {
			return nil, nil, err
		}

		if randomAccessOp != nil {
			variable = randomAccessOp
		}

		if dot := p.match(models.TokenTypesList[models.AttrAccessOp]); dot != nil {
			variable, err = p.parseAttrAccess(variable)
			if err != nil {
				return nil, nil, err
			}
		}

		return variable, nil, nil
	}

	return nil, nil, errors.New("очікується змінна або вираз")
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
			op := ast.NewRandomAccessOperationNode(name, expr, indexNode, lSquareBracket.Row)
			if lSquareBracket = p.match(models.TokenTypesList[models.LSquareBracket]); lSquareBracket != nil {
				p.pos--
				//ast.NewBinOperationNode()
				return p.parseRandomAccessOperation(nil, op)
			}

			op.Base = name
			return op, nil
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

	return expr, nil
}

// parseFunctionCall parses call of function.
// It assumes that name and left round bracket are parsed successfully.
func (p *Parser) parseFunctionCall(name *models.Token, parent ast.ExpressionNode) (ast.ExpressionNode, error) {
	lPar := p.current()
	if lPar != nil && lPar.Type.Name == models.LPar {
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

		return ast.NewFunctionCallNode(*name, parent, args), nil
	}

	return nil, errors.New("очікується відкриваюча дужка")
}

func (p *Parser) parseExpression() (ast.ExpressionNode, error) {
	variableNode, nameToken, err := p.parseVariableOrConstant()
	if err != nil {
		return nil, err
	}

	if variableNode != nil {
		randomAccessOp, err := p.parseRandomAccessOperation(nameToken, variableNode)
		if err != nil {
			return nil, err
		}

		if randomAccessOp != nil {
			variableNode = randomAccessOp
		}

		return variableNode, nil
	}

	if nameToken != nil {
		funcCallNode, err := p.parseFunctionCall(nameToken, nil)
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

		if dot := p.match(models.TokenTypesList[models.AttrAccessOp]); dot != nil {
			return p.parseAttrAccess(funcCallNode)
		}

		return funcCallNode, nil
	}

	return nil, errors.New("очікується змінна, або виклик функції")
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

		var err error
		var variableNode ast.ExpressionNode = ast.NewVariableNode(*name)
		variableNode, err = p.parseRandomAccessOperation(name, variableNode)
		if err != nil {
			return nil, err
		}

		if p.match(models.TokenTypesList[models.AttrAccessOp]) != nil {
			variableNode, err = p.parseAttrAccess(variableNode)
		} else {
			err := p.checkForKeyword(name.Text)
			if err != nil {
				return nil, err
			}
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
