package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
)

type VariableNode struct {
	Variable models.Token
}

func NewVariableNode(variable models.Token) VariableNode {
	return VariableNode{Variable: variable}
}

func (n VariableNode) String() string {
	return "VariableNode"
}
