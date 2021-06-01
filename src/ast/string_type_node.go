package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strings"
)

type StringTypeNode struct {
	Value models.Token

	rowNumber int
}

func NewStringTypeNode(token models.Token) StringTypeNode {
	token.Text = strings.TrimPrefix(strings.TrimSuffix(token.Text, "\""), "\"")
	return StringTypeNode{
		Value:     token,
		rowNumber: token.Row,
	}
}

func (n StringTypeNode) String() string {
	return "\"" + n.Value.Text + "\""
}

func (n StringTypeNode) RowNumber() int {
	return n.rowNumber
}
