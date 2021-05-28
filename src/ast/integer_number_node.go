package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type IntegerNumberNode struct {
	Number models.Token

	rowNumber int
}

func NewIntegerNumberNode(number models.Token) IntegerNumberNode {
	return IntegerNumberNode{
		Number:    number,
		rowNumber: number.Row,
	}
}

func (n IntegerNumberNode) String() string {
	return n.Number.Text
}

func (n IntegerNumberNode) RowNumber() int {
	return n.rowNumber
}
