package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeLogicalOp(
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
	case ops.AndOp:
		res, err = left.And(right)
		if err != nil {
			return nil, err
		}
	case ops.OrOp:
		res, err = left.Or(right)
		if err != nil {
			return nil, err
		}
	default:
		panic("fatal: invalid binary operator")
	}

	if res != nil {
		return res, nil
	}

	return nil, util.OperatorError(opType.Description(), left.TypeName(), right.TypeName())
}
