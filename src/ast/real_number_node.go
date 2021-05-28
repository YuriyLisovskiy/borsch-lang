package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type RealNumberNode struct {
	Number models.Token

	rowNumber int
}

func NewRealNumberNode(number models.Token) RealNumberNode {
	return RealNumberNode{
		Number:    number,
		rowNumber: number.Row,
	}
}

func (n RealNumberNode) String() string {
	return n.Number.Text
}

func (n RealNumberNode) RowNumber() int {
	return n.rowNumber
}
