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
) (types.ValueType, error) {
	operand, _, err := i.executeNode(node.Operand, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	var operator ops.Operator
	var res types.ValueType
	switch node.Operator.Type.Name {
	case models.Add:
		operator = ops.UnaryPlus
		res, err = operand.Plus()
		if err != nil {
			return nil, err
		}
	case models.Sub:
		operator = ops.UnaryMinus
		res, err = operand.Minus()
		if err != nil {
			return nil, err
		}
	case models.BitwiseNotOp:
		operator = ops.UnaryBitwiseNotOp
		res, err = operand.BitwiseNot()
		if err != nil {
			return nil, err
		}
	case models.NotOp:
		operator = ops.NotOp
		res, err = operand.Not()
		if err != nil {
			return nil, err
		}
	default:
		return nil, util.RuntimeError("невідомий унарний оператор")
	}

	if res != nil {
		return res, nil
	}

	return nil, util.RuntimeError(fmt.Sprintf(
		"непідтримуваний тип операнда для унарного оператора %s: '%s'",
		operator.Description(), operand.TypeName(),
	))
}
