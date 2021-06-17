package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type RandomAccessOperationNode struct {
	Base    *models.Token
	Operand ExpressionNode
	Index   ExpressionNode
	IsSet   bool

	rowNumber int
}

func NewRandomAccessOperationNode(
	base *models.Token, operand, index ExpressionNode, rowNumber int,
) RandomAccessOperationNode {
	return RandomAccessOperationNode{
		Base:      base,
		Operand:   operand,
		Index:     index,
		rowNumber: rowNumber,
	}
}

func (n RandomAccessOperationNode) String() string {
	return n.Operand.String() + "[" + n.Index.String() + "]"
}

func (n RandomAccessOperationNode) RowNumber() int {
	return n.rowNumber
}
