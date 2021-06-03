package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeLogicalOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType Operator, rootDir string, currentFile string,
) (builtin.ValueType, error) {
	left, err := i.executeNode(leftNode, rootDir, currentFile)
	if err != nil {
		return builtin.NoneType{}, err
	}

	right, err := i.executeNode(rightNode, rootDir, currentFile)
	if err != nil {
		return builtin.NoneType{}, err
	}

	if left.TypeHash() != right.TypeHash() {
		return builtin.NoneType{}, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				opType.Description(), left.TypeName(), right.TypeName(),
			),
		)
	}

	switch opType {
	case andOp:
		switch leftVal := left.(type) {
		case builtin.BoolType:
			return builtin.BoolType{
				Value: leftVal.Value && right.(builtin.BoolType).Value,
			}, nil
		}
	case orOp:
		switch leftVal := left.(type) {
		case builtin.BoolType:
			return builtin.BoolType{
				Value: leftVal.Value || right.(builtin.BoolType).Value,
			}, nil
		}

	default:
		return builtin.NoneType{}, util.RuntimeError("невідомий оператор")
	}

	return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opType.Description(), left.TypeName(), right.TypeName(),
	))
}
