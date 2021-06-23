package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strings"
)

type CallOpNode struct {
	CallableName models.Token
	Parent       ExpressionNode
	Args         []ExpressionNode

	rowNumber int
}

func NewCallOpNode(callableName models.Token, parent ExpressionNode, args []ExpressionNode) CallOpNode {
	return CallOpNode{
		CallableName: callableName,
		Parent:       parent,
		Args:         args,
		rowNumber:    callableName.Row,
	}
}

func (n CallOpNode) String() string {
	var args []string
	for _, arg := range n.Args {
		args = append(args, arg.String())
	}

	parent := ""
	if n.Parent != nil {
		parent = n.Parent.String() + "."
	}

	return parent + n.CallableName.Text + "(" + strings.Join(args, ", ") + ")"
}

func (n CallOpNode) RowNumber() int {
	return n.rowNumber
}
