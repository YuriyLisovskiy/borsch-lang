package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type BinOperationNode struct {
	Operator models.Token
	LeftNode ExpressionNode
	RightNode ExpressionNode
}

func NewBinOperationNode(
	operator models.Token, leftNode ExpressionNode, rightNode ExpressionNode,
) BinOperationNode {
	return BinOperationNode{
		Operator:  operator,
		LeftNode:  leftNode,
		RightNode: rightNode,
	}
}

func (bo BinOperationNode) String() string {
	return "BinOperationNode"
}
