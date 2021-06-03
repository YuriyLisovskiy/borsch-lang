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

type IfSequenceNode struct {
	Blocks    []ConditionBlock
	ElseBlock []models.Token
}

func NewIfSequenceNode(mainCondition ExpressionNode, mainConditionTokens []models.Token) IfSequenceNode {
	return IfSequenceNode{
		Blocks: []ConditionBlock{
			NewConditionBlock(mainCondition, mainConditionTokens),
		},
	}
}

func (n IfSequenceNode) String() string {
	return n.Blocks[0].String()
}

func (n IfSequenceNode) RowNumber() int {
	return n.Blocks[0].RowNumber()
}
