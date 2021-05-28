package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strings"
)

type StringNode struct {
	Content models.Token

	rowNumber int
}

func NewStringNode(token models.Token) StringNode {
	token.Text = strings.TrimPrefix(strings.TrimSuffix(token.Text, "\""), "\"")
	return StringNode{
		Content: token,
		rowNumber: token.Row,
	}
}

func (n StringNode) String() string {
	return "\"" + n.Content.Text + "\""
}

func (n StringNode) RowNumber() int {
	return n.rowNumber
}
