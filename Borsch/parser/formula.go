package parser

import (
	"errors"
	"github.com/YuriyLisovskiy/borsch/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch/Borsch/models"
)

type stack []interface{}

func (s stack) empty() bool {
	return len(s) == 0
}

func (s *stack) push(elem interface{}) {
	*s = append(*s, elem)
}

func (s *stack) top() interface{} {
	if s.empty() {
		panic("unable to get top value from empty stack")
	}

	return (*s)[len(*s)-1]
}

func (s *stack) pop() {
	if s.empty() {
		panic("unable to pop from empty stack")
	}

	*s = (*s)[:len(*s)-1]
}

func precedence(op models.Token) int {
	switch op.Type.Name {
	case models.OrOp:
		return 1
	case models.AndOp:
		return 2
	case models.NotOp:
		return 3
	case models.LessOp, models.LessOrEqualsOp, models.EqualsOp, models.NotEqualsOp, models.GreaterOp, models.GreaterOrEqualsOp:
		return 4
	case models.BitwiseOrOp:
		return 5
	case models.BitwiseXorOp:
		return 6
	case models.BitwiseAndOp:
		return 7
	case models.BitwiseLeftShiftOp, models.BitwiseRightShiftOp:
		return 8
	case models.Add, models.Sub:
		if op.IsUnaryOperator {
			return 11
		}

		return 9
	case models.Mul, models.Div, models.ModuloOp:
		return 10
	case models.ExponentOp:
		return 12
	case models.BitwiseNotOp:
		return 11
	default:
		return 0
	}
}

func buildOperationNode(nodes, operators *stack) {
	op := operators.top().(models.Token)
	operators.pop()
	var resultNode ast.ExpressionNode
	if op.IsUnaryOperator {
		resultNode = ast.NewUnaryOperationNode(op, nodes.top().(ast.ExpressionNode))
		nodes.pop()
	} else {
		rightNode := nodes.top().(ast.ExpressionNode)
		nodes.pop()
		leftNode := nodes.top().(ast.ExpressionNode)
		nodes.pop()
		resultNode = ast.NewBinOperationNode(op, leftNode, rightNode)
	}

	nodes.push(resultNode)
}

func (p *Parser) matchBinaryOperator() *models.Token {
	return p.match(
		models.TokenTypesList[models.ExponentOp], models.TokenTypesList[models.ModuloOp],
		models.TokenTypesList[models.Add], models.TokenTypesList[models.Sub],
		models.TokenTypesList[models.Mul], models.TokenTypesList[models.Div],
		models.TokenTypesList[models.AndOp], models.TokenTypesList[models.OrOp],
		models.TokenTypesList[models.EqualsOp], models.TokenTypesList[models.NotEqualsOp],
		models.TokenTypesList[models.GreaterOp], models.TokenTypesList[models.GreaterOrEqualsOp],
		models.TokenTypesList[models.LessOp], models.TokenTypesList[models.LessOrEqualsOp],
		models.TokenTypesList[models.BitwiseLeftShiftOp], models.TokenTypesList[models.BitwiseRightShiftOp],
		models.TokenTypesList[models.BitwiseAndOp], models.TokenTypesList[models.BitwiseXorOp],
		models.TokenTypesList[models.BitwiseOrOp],
	)
}

func (p *Parser) parseUnaryOperator() *models.Token {
	op := p.match(
		models.TokenTypesList[models.NotOp],
		models.TokenTypesList[models.Sub], models.TokenTypesList[models.Add],
		models.TokenTypesList[models.BitwiseNotOp],
	)
	if op != nil {
		op.IsUnaryOperator = true
		return op
	}

	return nil
}

func (p *Parser) parseExpression() (ast.ExpressionNode, error) {
	variableNode, nameToken, err := p.parseVariableOrConstant()
	if err != nil {
		return nil, err
	}

	if variableNode != nil {
		randomAccessOp, err := p.parseRandomAccessOperation(variableNode)
		if err != nil {
			return nil, err
		}

		if randomAccessOp != nil {
			variableNode = randomAccessOp
		}

		return variableNode, nil
	}

	if nameToken != nil {
		funcCallNode, err := p.parseFunctionCall(nameToken)
		if err != nil {
			return nil, err
		}

		randomAccessOp, err := p.parseRandomAccessOperation(funcCallNode)
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

		randomAccessOpNode, err := p.parseRandomAccessOperation(node)
		if err != nil {
			return nil, err
		}

		if randomAccessOpNode != nil {
			node = randomAccessOpNode
		}

		return node, nil
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseFormula() (ast.ExpressionNode, error) {
	nodes := stack{}
	operators := stack{}

	unaryOp := p.parseUnaryOperator()
	if unaryOp != nil {
		operators = append(operators, *unaryOp)
	}

	node, err := p.parseParentheses()
	if err != nil {
		return nil, err
	}

	nodes.push(node)
	operator := p.matchBinaryOperator()
	for operator != nil {
		for !operators.empty() && precedence(operators.top().(models.Token)) >= precedence(*operator) {
			buildOperationNode(&nodes, &operators)
		}

		unaryOp = p.parseUnaryOperator()
		node, err = p.parseParentheses()
		if err != nil {
			return nil, err
		}

		nodes.push(node)
		operators = append(operators, *operator)
		if unaryOp != nil {
			operators = append(operators, *unaryOp)
		}

		operator = p.matchBinaryOperator()
	}

	for !operators.empty() {
		buildOperationNode(&nodes, &operators)
	}

	return nodes.top().(ast.ExpressionNode), nil
}
