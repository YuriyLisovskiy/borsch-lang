package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func boolToInt(v bool) int64 {
	if v {
		return 1
	}

	return 0
}

func (i *Interpreter) executeArithmeticOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType Operator, rootDir string, currentFile string,
) (types.ValueType, error) {
	left, err := i.executeNode(leftNode, rootDir, currentFile)
	if err != nil {
		return types.NoneType{}, err
	}

	right, err := i.executeNode(rightNode, rootDir, currentFile)
	if err != nil {
		return types.NoneType{}, err
	}

	if left.TypeHash() != right.TypeHash() {
		return types.NoneType{}, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				opType.Description(), left.TypeName(), right.TypeName(),
			),
		)
	}

	switch opType {
	case sumOp:
		switch leftVal := left.(type) {
		case types.RealType:
			return types.RealType{
				Value: leftVal.Value + right.(types.RealType).Value,
			}, nil
		case types.IntegerType:
			return types.IntegerType{
				Value: leftVal.Value + right.(types.IntegerType).Value,
			}, nil
		case types.StringType:
			return types.StringType{
				Value: leftVal.Value + right.(types.StringType).Value,
			}, nil
		case types.BoolType:
			return types.IntegerType{
				Value: boolToInt(leftVal.Value) + boolToInt(right.(types.BoolType).Value),
			}, nil
		case types.ListType:
			return types.ListType{
				Values: append(leftVal.Values, right.(types.ListType).Values...),
			}, nil
		}

	case subOp:
		switch leftVal := left.(type) {
		case types.RealType:
			return types.RealType{
				Value: leftVal.Value - right.(types.RealType).Value,
			}, nil
		case types.IntegerType:
			return types.IntegerType{
				Value: leftVal.Value - right.(types.IntegerType).Value,
			}, nil
		case types.BoolType:
			return types.IntegerType{
				Value: boolToInt(leftVal.Value) - boolToInt(right.(types.BoolType).Value),
			}, nil
		}
	case mulOp:
		switch leftVal := left.(type) {
		case types.RealType:
			return types.RealType{
				Value: leftVal.Value * right.(types.RealType).Value,
			}, nil
		case types.IntegerType:
			return types.IntegerType{
				Value: leftVal.Value * right.(types.IntegerType).Value,
			}, nil
		case types.BoolType:
			return types.IntegerType{
				Value: boolToInt(leftVal.Value) * boolToInt(right.(types.BoolType).Value),
			}, nil
		}
	case divOp:
		switch leftVal := left.(type) {
		case types.RealType:
			rightVal := right.(types.RealType).Value
			if rightVal == 0 {
				return types.NoneType{}, util.RuntimeError("ділення на нуль")
			}

			return types.RealType{
				Value: leftVal.Value / right.(types.RealType).Value,
			}, nil
		case types.IntegerType:
			rightVal := right.(types.IntegerType).Value
			if rightVal == 0 {
				return types.NoneType{}, util.RuntimeError("ділення на нуль")
			}

			return types.RealType{
				Value: float64(leftVal.Value) / float64(right.(types.IntegerType).Value),
			}, nil
		case types.BoolType:
			rightVal := right.(types.BoolType).Value
			if !rightVal {
				return types.NoneType{}, util.RuntimeError("ділення на нуль")
			}

			return types.RealType{
				Value: float64(boolToInt(leftVal.Value) / boolToInt(right.(types.BoolType).Value)),
			}, nil
		}

	default:
		return types.NoneType{}, util.RuntimeError("невідомий оператор")
	}

	return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opTypeNames[opType], left.TypeName(), right.TypeName(),
	))
}
