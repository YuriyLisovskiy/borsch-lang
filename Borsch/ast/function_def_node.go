package ast

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch/Borsch/models"
	"strings"
)

type FunctionDefNode struct {
	Name       models.Token
	Arguments  []types.FunctionArgument
	ReturnType types.FunctionReturnType
	Body       []models.Token

	rowNumber int
}

func NewFunctionDefNode(
	name models.Token, args []types.FunctionArgument, retType types.FunctionReturnType, body []models.Token,
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
