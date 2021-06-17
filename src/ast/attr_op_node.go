package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
)

type AttrOpNode struct {
	Base       *models.Token
	Expression ExpressionNode
	Attr       models.Token

	rowNumber int
}

func NewGetAttrOpNode(base *models.Token, expression ExpressionNode, attr models.Token) AttrOpNode {
	return AttrOpNode{
		Base:       base,
		Expression: expression,
		Attr:       attr,
		rowNumber:  attr.Row,
	}
}

func (n AttrOpNode) String() string {
	return n.Expression.String() + "." + n.Attr.String()
}

func (n AttrOpNode) RowNumber() int {
	return n.rowNumber
}
