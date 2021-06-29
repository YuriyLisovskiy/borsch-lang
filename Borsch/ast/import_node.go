package ast

import (
	"github.com/YuriyLisovskiy/borsch/Borsch/models"
)

type ImportNode struct {
	Directive models.Token
	Name      string
	FilePath  string
	IsStd     bool

	rowNumber int
}

func NewImportNode(directive models.Token, name string, isStd bool) ImportNode {
	matches := directive.Type.Regex.FindAllStringSubmatch(directive.Text, -1)
	return ImportNode{
		Directive: directive,
		Name:      name,
		FilePath:  matches[0][1],
		IsStd:     isStd,
		rowNumber: directive.Row,
	}
}

func (n ImportNode) String() string {
	return n.Directive.String()
}

func (n ImportNode) RowNumber() int {
	return n.rowNumber
}
