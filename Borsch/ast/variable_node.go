package ast

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
)

type VariableNode struct {
	Variable models.Token
	Doc      string

	rowNumber int
}

func NewVariableNode(variable models.Token, doc string) VariableNode {
	return VariableNode{
		Variable:  variable,
		Doc:       doc,
		rowNumber: variable.Row,
	}
}

func (n VariableNode) String() string {
	return n.Variable.String()
}

func (n VariableNode) RowNumber() int {
	return n.rowNumber
}
