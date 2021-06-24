package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
)

type AttrOpNode struct {
	Base       *models.Token
	Expression ExpressionNode
	Attr       ExpressionNode

	rowNumber int
}

func NewGetAttrOpNode(base *models.Token, expression ExpressionNode, attr ExpressionNode, rowNumber int) AttrOpNode {
	return AttrOpNode{
		Base:       base,
		Expression: expression,
		Attr:       attr,
		rowNumber:  rowNumber,
	}
}

func (n AttrOpNode) String() string {
	res := n.Expression.String() + "."
	switch attr := n.Attr.(type) {
	case VariableNode:
		res += attr.Variable.Text
	case CallOpNode:
		res = attr.String()
	default:
		panic("fatal error")
	}

	return res
}

func (n AttrOpNode) RowNumber() int {
	return n.rowNumber
}
