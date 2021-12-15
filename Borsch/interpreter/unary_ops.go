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

	var operator ops.Operator
	var res types.Type
	switch node.Operator.Type.Name {
	case models.Add:
		operator = ops.UnaryPlus
	case models.Sub:
		operator = ops.UnaryMinus
	case models.BitwiseNotOp:
		operator = ops.UnaryBitwiseNotOp
	case models.NotOp:
		operator = ops.NotOp
	default:
		return nil, util.RuntimeError("невідомий унарний оператор")
	}

	operatorFunc, err := operand.GetAttribute(operator.Caption())
	if err != nil {
		return nil, util.RuntimeError(err.Error())
	}

	switch operator := operatorFunc.(type) {
	case types.FunctionType:
		res, err = operator.Callable([]types.Type{operand}, map[string]types.Type{"я": operand})
	default:
		// TODO: повернути повідомлення, що атрибут не callable!
		panic("NOT CALLABLE!")
	}

	if res != nil {
		return res, nil
	}

	return nil, util.RuntimeError(
		fmt.Sprintf(
			"непідтримуваний тип операнда для унарного оператора %s: '%s'",
			operator.Description(), operand.GetTypeName(),
		),
	)
}
