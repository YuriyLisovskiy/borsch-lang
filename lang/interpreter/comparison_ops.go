package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/lang/ast"
	"github.com/YuriyLisovskiy/borsch/lang/builtin/ops"
	"github.com/YuriyLisovskiy/borsch/lang/builtin/types"
	"github.com/YuriyLisovskiy/borsch/lang/util"
)

func (i *Interpreter) executeComparisonOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType ops.Operator,
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
			return nil, util.RuntimeError(fmt.Sprintf(err.Error(), opType.Description()))
		}

		switch opType {
		case ops.EqualsOp:
			return types.BoolType{Value: res == 0}, nil
		case ops.NotEqualsOp:
			return types.BoolType{Value: res != 0}, nil
		case ops.GreaterOp, ops.GreaterOrEqualsOp, ops.LessOp, ops.LessOrEqualsOp:
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
			return nil, util.RuntimeError(fmt.Sprintf(err.Error(), opType.Description()))
		}

		switch opType {
		case ops.EqualsOp:
			return types.BoolType{Value: res == 0}, nil
		case ops.NotEqualsOp:
			return types.BoolType{Value: res != 0}, nil
		case ops.GreaterOp:
			return types.BoolType{Value: res == 1}, nil
		case ops.GreaterOrEqualsOp:
			return types.BoolType{Value: res == 0 || res == 1}, nil
		case ops.LessOp:
			return types.BoolType{Value: res == -1}, nil
		case ops.LessOrEqualsOp:
			return types.BoolType{Value: res == 0 || res == -1}, nil
		default:
			return nil, util.RuntimeError("невідомий оператор")
		}
	}
}
