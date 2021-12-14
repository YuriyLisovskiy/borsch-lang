package ast

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
)

type ForEachNode struct {
	For       models.Token
	IndexVar  models.Token
	ItemVar   models.Token
	Container ExpressionNode
	Body      []models.Token
}

func NewForEachNode(
	forToken, indexVar, itemVar models.Token,
	container ExpressionNode,
	bodyTokens []models.Token,
) ForEachNode {
	return ForEachNode{
		For:       forToken,
		IndexVar:  indexVar,
		ItemVar:   itemVar,
		Container: container,
		Body:      bodyTokens,
	}
}

func (n ForEachNode) String() string {
	return n.For.String()
}

func (n ForEachNode) RowNumber() int {
	return n.For.Row
}
