package ast

import "github.com/YuriyLisovskiy/borsch/src/models"

type FunctionCallNode struct {
	FunctionName models.Token
	Args         []ExpressionNode
}

func NewFunctionCallNode(functionName models.Token, args []ExpressionNode) FunctionCallNode {
	return FunctionCallNode{
		FunctionName: functionName,
		Args:         args,
	}
}

func (n FunctionCallNode) String() string {
	return n.FunctionName.Text + "(...)"
}
