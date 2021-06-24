package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeComparisonOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType Operator,
	rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	left, _, err := i.executeNode(leftNode, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	right, _, err := i.executeNode(rightNode, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	switch left.(type) {
	case types.NilType, types.BoolType:
		res, err := left.CompareTo(right)
		if err != nil {
			return nil, util.RuntimeError(fmt.Sprintf(err.Error(), opTypeNames[opType]))
		}

		switch opType {
		case equalsOp:
			return types.BoolType{Value: res == 0}, nil
		case notEqualsOp:
			return types.BoolType{Value: res != 0}, nil
		case greaterOp, greaterOrEqualsOp, lessOp, lessOrEqualsOp:
			return nil, util.RuntimeError(fmt.Sprintf(
				"оператор %s невизначений для значень типів '%s' та '%s'",
				opType.Description(), left.TypeName(), right.TypeName(),
			))
		default:
			return nil, util.RuntimeError("невідомий оператор")
		}
	default:
		res, err := left.CompareTo(right)
		if err != nil {
			return nil, util.RuntimeError(fmt.Sprintf(err.Error(), opTypeNames[opType]))
		}

		switch opType {
		case equalsOp:
			return types.BoolType{Value: res == 0}, nil
		case notEqualsOp:
			return types.BoolType{Value: res != 0}, nil
		case greaterOp:
			return types.BoolType{Value: res == 1}, nil
		case greaterOrEqualsOp:
			return types.BoolType{Value: res == 0 || res == 1}, nil
		case lessOp:
			return types.BoolType{Value: res == -1}, nil
		case lessOrEqualsOp:
			return types.BoolType{Value: res == 0 || res == -1}, nil
		default:
			return nil, util.RuntimeError("невідомий оператор")
		}
	}
}
