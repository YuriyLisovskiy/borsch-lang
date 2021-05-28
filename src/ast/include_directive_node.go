package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
)

type IncludeDirectiveNode struct {
	Directive models.Token
	FilePath  string

	rowNumber int
}

func NewIncludeDirectiveNode(directive models.Token) IncludeDirectiveNode {
	matches := directive.Type.Regex.FindAllStringSubmatch(directive.Text, -1)
	return IncludeDirectiveNode{
		Directive: directive,
		FilePath:  matches[0][1],
		rowNumber: directive.Row,
	}
}

func (n IncludeDirectiveNode) String() string {
	return n.Directive.String()
}

func (n IncludeDirectiveNode) RowNumber() int {
	return n.rowNumber
}
