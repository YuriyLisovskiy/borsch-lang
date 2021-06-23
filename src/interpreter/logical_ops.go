package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeLogicalOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType Operator,
	rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	left, err := i.executeNode(leftNode, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	right, err := i.executeNode(rightNode, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	if left.TypeHash() != right.TypeHash() {
		return nil, util.RuntimeError(
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
		return nil, util.RuntimeError("невідомий оператор")
	}

	return nil, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opType.Description(), left.TypeName(), right.TypeName(),
	))
}
