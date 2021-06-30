package interpreter

import (
	"github.com/YuriyLisovskiy/borsch/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch/Borsch/util"
)

func (i *Interpreter) executeBitwiseOp(
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

	var res types.ValueType
	switch opType {
	case ops.BitwiseLeftShiftOp:
		res, err = left.Add(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	case ops.BitwiseRightShiftOp:
		res, err = left.Sub(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	case ops.BitwiseAndOp:
		res, err = left.Mul(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	case ops.BitwiseXorOp:
		res, err = left.Div(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	case ops.BitwiseOrOp:
		res, err = left.Pow(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	default:
		panic("fatal: invalid bitwise operator")
	}

	if res != nil {
		return res, nil
	}

	return nil, util.OperatorError(opType.Description(), left.TypeName(), right.TypeName())
}
