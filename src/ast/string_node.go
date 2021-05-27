package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strings"
)

type StringNode struct {
	Content models.Token
}

func NewStringNode(token models.Token) StringNode {
	token.Text = strings.TrimPrefix(strings.TrimSuffix(token.Text, "\""), "\"")
	return StringNode{Content: token}
}

func (n StringNode) String() string {
	return "\"" + n.Content.Text + "\""
}
