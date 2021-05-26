package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type UnaryOperationNode struct {
	Operator models.Token
	Operand ExpressionNode
}

func NewUnaryOperationNode(operator models.Token, operand ExpressionNode) UnaryOperationNode {
	return UnaryOperationNode{
		Operator:  operator,
		Operand:  operand,
	}
}

func (n UnaryOperationNode) String() string {
	return "UnaryOperationNode"
}
