package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeArithmeticOp(
	ctx *Context,
	leftNode ast.ExpressionNode,
	rightNode ast.ExpressionNode,
	opType ops.Operator,
) (types.Type, error) {
	left, _, err := i.executeNode(ctx, leftNode)
	if err != nil {
		return nil, err
	}

	right, _, err := i.executeNode(ctx, rightNode)
	if err != nil {
		return nil, err
	}

	var res types.Type
	switch opType {
	case ops.AddOp, ops.SubOp, ops.MulOp, ops.DivOp, ops.PowOp, ops.ModuloOp:
		operatorFunc, err := left.GetAttribute(opType.Caption())
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}

		switch operator := operatorFunc.(type) {
		case *types.FunctionInstance:
			args := []types.Type{left, right}
			kwargs := map[string]types.Type{
				"я":     left,
				"інший": right,
			}
			if err := types.CheckFunctionArguments(operator, &args, &kwargs); err != nil {
				return nil, err
			}

			res, err = operator.Call(&args, &kwargs)
			if err != nil {
				return nil, util.RuntimeError(err.Error())
			}
		default:
			return nil, util.ObjectIsNotCallable(opType.Caption(), operatorFunc.GetTypeName())
		}
	default:
		panic("fatal: invalid arithmetic operator")
	}

	if res != nil {
		return res, nil
	}

	return nil, util.OperatorError(opType.Description(), left.GetTypeName(), right.GetTypeName())
}
