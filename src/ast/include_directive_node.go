package ast

import (
	"github.com/YuriyLisovskiy/borsch/src/models"
)

type IncludeDirectiveNode struct {
	FilePath string
}

func NewIncludeDirectiveNode(directive models.Token) IncludeDirectiveNode {
	matches := directive.Type.Regex.FindAllStringSubmatch(directive.Text, -1)
	return IncludeDirectiveNode{
		FilePath: matches[0][1],
	}
}

func (n IncludeDirectiveNode) String() string {
	return "'" + n.FilePath + "'"
}
