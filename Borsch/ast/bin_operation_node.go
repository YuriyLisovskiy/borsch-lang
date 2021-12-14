package ast

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/models"

type BinOperationNode struct {
	Operator  models.Token
	LeftNode  ExpressionNode
	RightNode ExpressionNode

	rowNumber int
}

func NewBinOperationNode(
	operator models.Token, leftNode ExpressionNode, rightNode ExpressionNode,
) BinOperationNode {
	return BinOperationNode{
		Operator:  operator,
		LeftNode:  leftNode,
		RightNode: rightNode,
		rowNumber: operator.Row,
	}
}

func (n BinOperationNode) String() string {
	return n.LeftNode.String() + " " + n.Operator.String() + " " + n.RightNode.String()
}

func (n BinOperationNode) RowNumber() int {
	return n.rowNumber
}
