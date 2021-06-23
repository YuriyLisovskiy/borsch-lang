package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
)

type ConditionBlock struct {
	Condition ExpressionNode
	Tokens    []models.Token
}

func NewConditionBlock(condition ExpressionNode, tokens []models.Token) ConditionBlock {
	return ConditionBlock{
		Condition: condition,
		Tokens:    tokens,
	}
}

func (n ConditionBlock) String() string {
	return n.Condition.String()
}

func (n ConditionBlock) RowNumber() int {
	return n.Condition.RowNumber()
}

type IfNode struct {
	Blocks    []ConditionBlock
	ElseBlock []models.Token
}

func NewIfNode(mainCondition ExpressionNode, mainConditionTokens []models.Token) IfNode {
	return IfNode{
		Blocks: []ConditionBlock{
			NewConditionBlock(mainCondition, mainConditionTokens),
		},
	}
}

func (n IfNode) String() string {
	return n.Blocks[0].String()
}

func (n IfNode) RowNumber() int {
	return n.Blocks[0].RowNumber()
}
