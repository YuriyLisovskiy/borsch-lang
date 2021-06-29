package ast

import (
	"github.com/YuriyLisovskiy/borsch/lang/models"
	"strings"
)

type CallOpNode struct {
	CallableName models.Token
	Parameters   []ExpressionNode

	rowNumber int
}

func NewCallOpNode(callableName models.Token, parameters []ExpressionNode) CallOpNode {
	return CallOpNode{
		CallableName: callableName,
		Parameters:   parameters,
		rowNumber:    callableName.Row,
	}
}

func (n CallOpNode) String() string {
	var parameters []string
	for _, arg := range n.Parameters {
		parameters = append(parameters, arg.String())
	}

	return n.CallableName.Text + "(" + strings.Join(parameters, ", ") + ")"
}

func (n CallOpNode) RowNumber() int {
	return n.rowNumber
}
