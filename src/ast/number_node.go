package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type NumberNode struct {
	Number models.Token
}

func NewNumberNode(number models.Token) NumberNode {
	return NumberNode{Number: number}
}

func (n NumberNode) String() string {
	return n.Number.Text
}
