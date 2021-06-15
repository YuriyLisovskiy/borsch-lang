package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
)

type GetAttrOpNode struct {
	Parent ExpressionNode
	Attr   models.Token

	rowNumber int
}

func NewGetAttrOpNode(parent ExpressionNode, attr models.Token) GetAttrOpNode {
	return GetAttrOpNode{
		Parent:  parent,
		Attr:     attr,
		rowNumber: attr.Row,
	}
}

func (n GetAttrOpNode) String() string {
	return n.Parent.String() + "." + n.Attr.String()
}

func (n GetAttrOpNode) RowNumber() int {
	return n.rowNumber
}
