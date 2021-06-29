package interpreter

import (
	"github.com/YuriyLisovskiy/borsch/lang/ast"
	"github.com/YuriyLisovskiy/borsch/lang/builtin/ops"
	"github.com/YuriyLisovskiy/borsch/lang/builtin/types"
	"github.com/YuriyLisovskiy/borsch/lang/util"
)

func (i *Interpreter) executeArithmeticOp(
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
	case ops.AddOp:
		res, err = left.Add(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	case ops.SubOp:
		res, err = left.Sub(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	case ops.MulOp:
		res, err = left.Mul(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	case ops.DivOp:
		res, err = left.Div(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	case ops.PowOp:
		res, err = left.Pow(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	case ops.ModuloOp:
		res, err = left.Mod(right)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}
	default:
		return nil, util.RuntimeError("невідомий оператор")
	}
	
	if res != nil {
		return res, nil
	}

	return nil, util.OperatorError(opType.Description(), left.TypeName(), right.TypeName())
}
