package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strings"
)

type FunctionCallNode struct {
	FunctionName models.Token
	Args         []ExpressionNode

	rowNumber int
}

func NewFunctionCallNode(name models.Token, args []ExpressionNode) FunctionCallNode {
	return FunctionCallNode{
		FunctionName: name,
		Args:         args,
		rowNumber:    name.Row,
	}
}

func (n FunctionCallNode) String() string {
	var args []string
	for _, arg := range n.Args {
		args = append(args, arg.String())
	}

	return n.FunctionName.Text + "(" + strings.Join(args, ", ") + ")"
}

func (n FunctionCallNode) RowNumber() int {
	return n.rowNumber
}
