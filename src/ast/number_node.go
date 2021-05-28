package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type NumberNode struct {
	Number models.Token

	rowNumber int
}

func NewNumberNode(number models.Token) NumberNode {
	return NumberNode{
		Number:    number,
		rowNumber: number.Row,
	}
}

func (n NumberNode) String() string {
	return n.Number.Text
}

func (n NumberNode) RowNumber() int {
	return n.rowNumber
}
