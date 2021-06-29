package ast

import (
	"fmt"
)

type ReturnNode struct {
	Value ExpressionNode

	rowNumber int
}

func NewReturnNode(value ExpressionNode, rowNumber int) ReturnNode {
	return ReturnNode{
		Value:     value,
		rowNumber: rowNumber,
	}
}

func (n ReturnNode) String() string {
	strValue := ""
	if n.Value != nil {
		strValue = n.Value.String()
	}

	if strValue != "" {
		strValue = " " + strValue
	}

	return fmt.Sprintf("повернути%s;", strValue)
}

func (n ReturnNode) RowNumber() int {
	return n.rowNumber
}
