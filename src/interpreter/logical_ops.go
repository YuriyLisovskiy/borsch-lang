package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (e *Interpreter) executeLogicalOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType Operator, rootDir string,
) (builtin.ValueType, error) {
	left, err := e.executeNode(leftNode, rootDir)
	if err != nil {
		return builtin.NoneType{}, err
	}

	right, err := e.executeNode(rightNode, rootDir)
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
