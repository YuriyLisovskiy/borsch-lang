package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeArithmeticOp(
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
	case ops.AddOp, ops.SubOp, ops.MulOp, ops.DivOp, ops.PowOp, ops.ModuloOp:
		operatorFunc, err := left.GetAttribute(opType.Caption())
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}

		switch operator := operatorFunc.(type) {
		case types.FunctionType:
			res, err = operator.Callable(
				[]types.Type{left, right},
				map[string]types.Type{
					"я": left,
					"інший": right,
				},
			)
		default:
			// TODO: повернути повідомлення, що атрибут не callable!
			panic("NOT CALLABLE!")
		}


	// TODO: remove commented code!
	// case ops.AddOp:
	// 	addFunc, err := left.GetAttr("__додати__")
	// 	if err != nil {
	// 		return nil, util.RuntimeError(err.Error())
	// 	}
	//
	// 	res, err = addFunc.(types.FunctionType).Callable(
	// 		[]types.ValueType{left, right},
	// 		map[string]types.ValueType{
	// 			"я": left,
	// 			"інший": right,
	// 		},
	// 	)
	//
	// 	// res, err = left.Add(right)
	// 	if err != nil {
	// 		return nil, util.RuntimeError(err.Error())
	// 	}
	// case ops.SubOp:
	// 	res, err = left.Sub(right)
	// 	if err != nil {
	// 		return nil, util.RuntimeError(err.Error())
	// 	}
	// case ops.MulOp:
	// 	res, err = left.Mul(right)
	// 	if err != nil {
	// 		return nil, util.RuntimeError(err.Error())
	// 	}
	// case ops.DivOp:
	// 	res, err = left.Div(right)
	// 	if err != nil {
	// 		return nil, util.RuntimeError(err.Error())
	// 	}
	// case ops.PowOp:
	// 	res, err = left.Pow(right)
	// 	if err != nil {
	// 		return nil, util.RuntimeError(err.Error())
	// 	}
	// case ops.ModuloOp:
	// 	res, err = left.Mod(right)
	// 	if err != nil {
	// 		return nil, util.RuntimeError(err.Error())
	// 	}


	default:
		panic("fatal: invalid arithmetic operator")
	}
	
	if res != nil {
		return res, nil
	}

	return nil, util.OperatorError(opType.Description(), left.GetTypeName(), right.GetTypeName())
}
