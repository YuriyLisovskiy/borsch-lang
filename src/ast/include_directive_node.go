package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
)

type IncludeDirectiveNode struct {
	Directive models.Token
	Name      string
	FilePath  string
	IsStd     bool

	rowNumber int
}

func NewIncludeDirectiveNode(directive models.Token, name string, isStd bool) IncludeDirectiveNode {
	matches := directive.Type.Regex.FindAllStringSubmatch(directive.Text, -1)
	return IncludeDirectiveNode{
		Directive: directive,
		Name:      name,
		FilePath:  matches[0][1],
		IsStd:     isStd,
		rowNumber: directive.Row,
	}
}

func (n IncludeDirectiveNode) String() string {
	return n.Directive.String()
}

func (n IncludeDirectiveNode) RowNumber() int {
	return n.rowNumber
}
