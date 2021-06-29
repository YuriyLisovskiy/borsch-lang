package ast

import "github.com/YuriyLisovskiy/borsch/Borsch/models"

type UnaryOperationNode struct {
	Operator models.Token
	Operand  ExpressionNode

	rowNumber int
}

func NewUnaryOperationNode(
	operator models.Token, operand ExpressionNode,
) UnaryOperationNode {
	return UnaryOperationNode{
		Operator:  operator,
		Operand:   operand,
		rowNumber: operator.Row,
	}
}

func (n UnaryOperationNode) String() string {
	return n.Operator.String() + " (" + n.Operand.String() + ")"
}

func (n UnaryOperationNode) RowNumber() int {
	return n.rowNumber
}
