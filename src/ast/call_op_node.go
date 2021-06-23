package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strings"
)

type CallOpNode struct {
	CallableName models.Token
	Parent       ExpressionNode
	Parameters   []ExpressionNode

	rowNumber int
}

func NewCallOpNode(callableName models.Token, parent ExpressionNode, parameters []ExpressionNode) CallOpNode {
	return CallOpNode{
		CallableName: callableName,
		Parent:       parent,
		Parameters:   parameters,
		rowNumber:    callableName.Row,
	}
}

func (n CallOpNode) String() string {
	var parameters []string
	for _, arg := range n.Parameters {
		parameters = append(parameters, arg.String())
	}

	parent := ""
	if n.Parent != nil {
		parent = n.Parent.String() + "."
	}

	return parent + n.CallableName.Text + "(" + strings.Join(parameters, ", ") + ")"
}

func (n CallOpNode) RowNumber() int {
	return n.rowNumber
}
