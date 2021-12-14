package ast

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/models"

type RealTypeNode struct {
	Value models.Token

	rowNumber int
}

func NewRealTypeNode(token models.Token) RealTypeNode {
	return RealTypeNode{
		Value:     token,
		rowNumber: token.Row,
	}
}

func (n RealTypeNode) String() string {
	return n.Value.Text
}

func (n RealTypeNode) RowNumber() int {
	return n.rowNumber
}
