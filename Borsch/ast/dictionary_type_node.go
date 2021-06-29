package ast

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/Borsch/models"
	"strings"
)

type DictionaryTypeNode struct {
	Map map[ExpressionNode]ExpressionNode

	rowNumber int
}

func NewDictionaryTypeNode(token models.Token) DictionaryTypeNode {
	return DictionaryTypeNode{
		Map: map[ExpressionNode]ExpressionNode{},
		rowNumber: token.Row,
	}
}

func (n DictionaryTypeNode) String() string {
	var values []string
	for key, value := range n.Map {
		values = append(values, fmt.Sprintf("%s: %s", key.String(), value.String()))
	}

	return "{" + strings.Join(values, ", ") + "}"
}

func (n DictionaryTypeNode) RowNumber() int {
	return n.rowNumber
}
