package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeUnaryOp(
	node *ast.UnaryOperationNode, rootDir string, thisPackage, parentPackage string,
) (types.Type, error) {
	operand, _, err := i.executeNode(node.Operand, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	var op ops.Operator
	var res types.Type
	switch node.Operator.Type.Name {
	case models.Add:
		op = ops.UnaryPlus
	case models.Sub:
		op = ops.UnaryMinus
	case models.BitwiseNotOp:
		op = ops.UnaryBitwiseNotOp
	case models.NotOp:
		op = ops.NotOp
	default:
		return nil, util.RuntimeError("невідомий унарний оператор")
	}

	operatorFunc, err := operand.GetAttribute(op.Caption())
	if err != nil {
		return nil, util.RuntimeError(err.Error())
	}

	switch operator := operatorFunc.(type) {
	case types.CallableType:
		res, err = operator.Call([]types.Type{operand}, map[string]types.Type{"я": operand})
	default:
		return nil, util.ObjectIsNotCallable(op.Caption(), operatorFunc.GetTypeName())
	}

	if res != nil {
		return res, nil
	}

	return nil, util.RuntimeError(
		fmt.Sprintf(
			"непідтримуваний тип операнда для унарного оператора %s: '%s'",
			op.Description(), operand.GetTypeName(),
		),
	)
}
