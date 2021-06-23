package ast

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/models"
	"strings"
)

type FunctionDefNode struct {
	Name       models.Token
	Arguments  []types.FunctionArgument
	ReturnType int
	Body       []models.Token

	rowNumber int
}

func NewFunctionDefNode(
	name models.Token, args []types.FunctionArgument, retType int, body []models.Token,
) FunctionDefNode {
	return FunctionDefNode{
		Name:       name,
		Arguments:  args,
		ReturnType: retType,
		Body:       body,
		rowNumber:  name.Row,
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
		params = append(params, strParam)
	}

	return "функція " + n.Name.Text + "(" + strings.Join(params, ", ") + ") -> " + types.GetTypeName(n.ReturnType)
}

func (n FunctionDefNode) RowNumber() int {
	return n.rowNumber
}
