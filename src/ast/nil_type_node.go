package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type NilTypeNode struct {
	Value string

	rowNumber int
}

func NewNilTypeNode(token *models.Token) NilTypeNode {
	return NilTypeNode{
		Value:     token.Text,
		rowNumber: token.Row,
	}
}

func (n NilTypeNode) String() string {
	return n.Value
}

func (n NilTypeNode) RowNumber() int {
	return n.rowNumber
}
