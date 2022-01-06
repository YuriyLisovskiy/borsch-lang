package ast

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
)

type ClassDefNode struct {
	Name       models.Token
	Doc        string
	Attributes []ExpressionNode

	rowNumber int
}

func NewClassDefNode(name models.Token, doc string, attributes []ExpressionNode) ClassDefNode {
	return ClassDefNode{
		Name:       name,
		Doc:        doc,
		Attributes: attributes,
		rowNumber:  name.Row,
	}
}

func (n ClassDefNode) String() string {
	return "клас " + n.Name.Text
}

func (n ClassDefNode) RowNumber() int {
	return n.rowNumber
}
