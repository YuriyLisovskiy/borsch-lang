package ast

import (
	"github.com/YuriyLisovskiy/borsch/lang/models"
	"strings"
)

type ListTypeNode struct {
	Values []ExpressionNode

	rowNumber int
}

func NewListTypeNode(token models.Token, values []ExpressionNode) ListTypeNode {
	return ListTypeNode{
		Values:    values,
		rowNumber: token.Row,
	}
}

func (n ListTypeNode) String() string {
	var values []string
	for _, value := range n.Values {
		values = append(values, value.String())
	}

	return "[" + strings.Join(values, ", ") + "]"
}

func (n ListTypeNode) RowNumber() int {
	return n.rowNumber
}
