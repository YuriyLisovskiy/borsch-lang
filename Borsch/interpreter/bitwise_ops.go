package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeBitwiseOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType ops.Operator,
	rootDir string, thisPackage, parentPackage string,
) (types.Type, error) {
	left, _, err := i.executeNode(leftNode, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	right, _, err := i.executeNode(rightNode, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	var res types.Type
	switch opType {
	case ops.BitwiseLeftShiftOp, ops.BitwiseRightShiftOp, ops.BitwiseAndOp, ops.BitwiseXorOp, ops.BitwiseOrOp:
		operatorFunc, err := left.GetAttribute(opType.Caption())
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}

		switch operator := operatorFunc.(type) {
		case types.CallableType:
			res, err = operator.Call(
				[]types.Type{left, right}, map[string]types.Type{
					"я":     left,
					"інший": right,
				},
			)
			if err != nil {
				return nil, util.RuntimeError(err.Error())
			}
		default:
			return nil, util.ObjectIsNotCallable(opType.Caption(), operatorFunc.GetTypeName())
		}
	default:
		panic("fatal: invalid bitwise operator")
	}

	if res != nil {
		return res, nil
	}

	return nil, util.OperatorError(opType.Description(), left.GetTypeName(), right.GetTypeName())
}
