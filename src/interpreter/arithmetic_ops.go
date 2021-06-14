package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"math"
)

func boolToInt(v bool) int64 {
	if v {
		return 1
	}

	return 0
}

func boolToFloat64(v bool) float64 {
	if v {
		return 1.0
	}

	return 0.0
}

func (i *Interpreter) executeArithmeticOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType Operator, rootDir string, currentFile string,
) (types.ValueType, error) {
	left, err := i.executeNode(leftNode, rootDir, currentFile)
	if err != nil {
		return nil, err
	}

	right, err := i.executeNode(rightNode, rootDir, currentFile)
	if err != nil {
		return nil, err
	}

	if left.TypeHash() != right.TypeHash() {
		return nil, util.RuntimeError(
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
				return nil, util.RuntimeError("ділення на нуль")
			}

			return types.RealType{
				Value: leftVal.Value / right.(types.RealType).Value,
			}, nil
		case types.IntegerType:
			rightVal := right.(types.IntegerType).Value
			if rightVal == 0 {
				return nil, util.RuntimeError("ділення на нуль")
			}

			return types.RealType{
				Value: float64(leftVal.Value) / float64(right.(types.IntegerType).Value),
			}, nil
		case types.BoolType:
			rightVal := right.(types.BoolType).Value
			if !rightVal {
				return nil, util.RuntimeError("ділення на нуль")
			}

			return types.RealType{
				Value: float64(boolToInt(leftVal.Value) / boolToInt(right.(types.BoolType).Value)),
			}, nil
		}
	case exponentOp:
		switch leftVal := left.(type) {
		case types.RealType:
			return types.RealType{
				Value: math.Pow(leftVal.Value, right.(types.RealType).Value),
			}, nil
		case types.IntegerType:
			return types.IntegerType{
				Value: int64(math.Pow(float64(leftVal.Value), float64(right.(types.IntegerType).Value))),
			}, nil
		case types.BoolType:
			return types.IntegerType{
				Value: int64(math.Pow(boolToFloat64(leftVal.Value), boolToFloat64(right.(types.BoolType).Value))),
			}, nil
		}
	case moduloOp:
		switch leftVal := left.(type) {
		case types.IntegerType:
			rightVal := right.(types.IntegerType).Value
			if rightVal == 0 {
				return nil, util.RuntimeError("ділення на нуль")
			}

			return types.IntegerType{Value: leftVal.Value % rightVal}, nil
		case types.BoolType:
			rightVal := right.(types.BoolType).Value
			if !rightVal {
				return nil, util.RuntimeError("ділення на нуль")
			}

			return types.IntegerType{Value: boolToInt(leftVal.Value) % boolToInt(rightVal)}, nil
		}
	default:
		return nil, util.RuntimeError("невідомий оператор")
	}

	return nil, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opTypeNames[opType], left.TypeName(), right.TypeName(),
	))
}
