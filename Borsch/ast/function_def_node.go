package ast

import (
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
)

type FunctionDefNode struct {
	Name        models.Token
	Arguments   []types.FunctionArgument
	ReturnType  types.FunctionReturnType
	Body        []models.Token
	IsAnonymous bool

	rowNumber int
}

func NewFunctionDefNode(
	rowNumber int,
	name models.Token,
	args []types.FunctionArgument,
	retType types.FunctionReturnType,
	body []models.Token,
) FunctionDefNode {
	return FunctionDefNode{
		Name:        name,
		Arguments:   args,
		ReturnType:  retType,
		Body:        body,
		IsAnonymous: name.Text != "",
		rowNumber:   rowNumber,
	}
}

func (n FunctionDefNode) String() string {
	var params []string
	for _, param := range n.Arguments {
		strParam := fmt.Sprintf("%s: ", param.Name)
		if param.IsVariadic {
			strParam += "..."
		}

		strParam += param.TypeName()
		if param.IsNullable {
			strParam += "?"
		}

		params = append(params, strParam)
	}

	return "функція " + n.Name.Text + "(" + strings.Join(params, ", ") + ") -> " + n.ReturnType.String()
}

func (n FunctionDefNode) RowNumber() int {
	return n.rowNumber
}
