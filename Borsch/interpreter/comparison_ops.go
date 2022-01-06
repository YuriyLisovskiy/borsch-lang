package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeComparisonOp(
	leftOperand types.Type,
	rightOperand types.Type,
	opType ops.Operator,
) (types.Type, error) {
	switch leftOperand.(type) {
	case types.NilInstance, types.BoolInstance:
		switch opType {
		case ops.EqualsOp, ops.NotEqualsOp:
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

				res, err := operator.Call(&args, &kwargs)
				if err != nil {
					return nil, util.RuntimeError(fmt.Sprintf(err.Error(), opType.Description()))
				}

				return res, nil
			default:
				return nil, util.ObjectIsNotCallable(opType.Caption(), operatorFunc.GetTypeName())
			}
		case ops.GreaterOp, ops.GreaterOrEqualsOp, ops.LessOp, ops.LessOrEqualsOp:
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"оператор %s невизначений для значень типів '%s' та '%s'",
					opType.Description(), leftOperand.GetTypeName(), rightOperand.GetTypeName(),
				),
			)
		default:
			return nil, util.RuntimeError("невідомий оператор")
		}
	default:
		switch opType {
		case ops.EqualsOp, ops.NotEqualsOp, ops.GreaterOp, ops.GreaterOrEqualsOp, ops.LessOp, ops.LessOrEqualsOp:
			if operatorFunc, err := leftOperand.GetAttribute(opType.Caption()); err == nil {
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

					res, err := operator.Call(&args, &kwargs)
					if err != nil {
						return nil, util.RuntimeError(fmt.Sprintf(err.Error(), opType.Description()))
					}

					return res, nil
				default:
					return nil, util.ObjectIsNotCallable(opType.Caption(), operatorFunc.GetTypeName())
				}
			}
		default:
			return nil, util.RuntimeError("невідомий оператор")
		}
	}

	return nil, util.OperatorError(opType.Description(), leftOperand.GetTypeName(), rightOperand.GetTypeName())
}
