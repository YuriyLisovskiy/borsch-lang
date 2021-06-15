package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strings"
)

type FunctionCallNode struct {
	FunctionName models.Token
	Parent       ExpressionNode
	Args         []ExpressionNode

	rowNumber int
}

func NewFunctionCallNode(name models.Token, parent ExpressionNode, args []ExpressionNode) FunctionCallNode {
	return FunctionCallNode{
		FunctionName: name,
		Parent:       parent,
		Args:         args,
		rowNumber:    name.Row,
	}
}

func (n FunctionCallNode) String() string {
	var args []string
	for _, arg := range n.Args {
		args = append(args, arg.String())
	}

	parent := ""
	if n.Parent != nil {
		parent = n.Parent.String() + "."
	}

	return parent + n.FunctionName.Text + "(" + strings.Join(args, ", ") + ")"
}

func (n FunctionCallNode) RowNumber() int {
	return n.rowNumber
}
