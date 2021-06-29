package ast

import (
	"github.com/YuriyLisovskiy/borsch/lang/models"
	"strconv"
)

type BoolTypeNode struct {
	Value models.Token

	rowNumber int
}

func NewBoolTypeNode(token models.Token) BoolTypeNode {
	return BoolTypeNode{
		Value:     token,
		rowNumber: token.Row,
	}
}

func (n BoolTypeNode) String() string {
	return n.Value.Text
}

func (n BoolTypeNode) RowNumber() int {
	return n.rowNumber
}

func (n BoolTypeNode) AsBool() (bool, error) {
	return strconv.ParseBool(n.Value.Text)
}
