package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeLogicalOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType Operator, rootDir string, currentFile string,
) (types.ValueType, error) {
	left, err := i.executeNode(leftNode, rootDir, currentFile)
	if err != nil {
		return types.NoneType{}, err
	}

	right, err := i.executeNode(rightNode, rootDir, currentFile)
	if err != nil {
		return types.NoneType{}, err
	}

	if left.TypeHash() != right.TypeHash() {
		return types.NoneType{}, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				opType.Description(), left.TypeName(), right.TypeName(),
			),
		)
	}

	switch opType {
	case andOp:
		switch leftVal := left.(type) {
		case types.BoolType:
			return types.BoolType{
				Value: leftVal.Value && right.(types.BoolType).Value,
			}, nil
		}
	case orOp:
		switch leftVal := left.(type) {
		case types.BoolType:
			return types.BoolType{
				Value: leftVal.Value || right.(types.BoolType).Value,
			}, nil
		}

	default:
		return types.NoneType{}, util.RuntimeError("невідомий оператор")
	}

	return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opType.Description(), left.TypeName(), right.TypeName(),
	))
}
