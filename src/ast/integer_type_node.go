package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type IntegerTypeNode struct {
	Value models.Token

	rowNumber int
}

func NewIntegerTypeNode(token models.Token) IntegerTypeNode {
	return IntegerTypeNode{
		Value:     token,
		rowNumber: token.Row,
	}
}

func (n IntegerTypeNode) String() string {
	return n.Value.Text
}

func (n IntegerTypeNode) RowNumber() int {
	return n.rowNumber
}
