package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strconv"
)

type BoolNode struct {
	Value models.Token

	rowNumber int
}

func NewBoolNode(token models.Token) BoolNode {
	return BoolNode{
		Value: token,
		rowNumber: token.Row,
	}
}

func (n BoolNode) String() string {
	return n.Value.Text
}

func (n BoolNode) RowNumber() int {
	return n.rowNumber
}

func (n BoolNode) AsBool() (bool, error) {
	return strconv.ParseBool(n.Value.Text)
}
