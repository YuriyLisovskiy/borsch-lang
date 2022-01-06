package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeLogicalOp(
	leftOperand types.Type,
	rightOperand types.Type,
	opType ops.Operator,
) (types.Type, error) {
	var res types.Type
	switch opType {
	case ops.AndOp, ops.OrOp:
		operatorFunc, err := leftOperand.GetAttribute(opType.Caption())
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}

		switch operator := operatorFunc.(type) {
		case *types.FunctionInstance:
			args := []types.Type{leftOperand, rightOperand}
			kwargs := map[string]types.Type{
				operator.Arguments[0].Name: leftOperand,
				operator.Arguments[1].Name: rightOperand,
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
		panic("fatal: invalid binary operator")
	}

	if res != nil {
		return res, nil
	}

	return nil, util.OperatorError(opType.Description(), leftOperand.GetTypeName(), rightOperand.GetTypeName())
}
