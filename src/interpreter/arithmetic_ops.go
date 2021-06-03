package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (e *Interpreter) executeArithmeticOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType Operator, rootDir string, currentFile string,
) (builtin.ValueType, error) {
	left, err := e.executeNode(leftNode, rootDir, currentFile)
	if err != nil {
		return builtin.NoneType{}, err
	}

	right, err := e.executeNode(rightNode, rootDir, currentFile)
	if err != nil {
		return builtin.NoneType{}, err
	}

	if left.TypeHash() != right.TypeHash() {
		return builtin.NoneType{}, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				opType.Description(), left.TypeName(), right.TypeName(),
			),
		)
	}

	switch opType {
	case sumOp:
		switch leftVal := left.(type) {
		case builtin.RealNumberType:
			return builtin.RealNumberType{
				Value: leftVal.Value + right.(builtin.RealNumberType).Value,
			}, nil
		case builtin.IntegerNumberType:
			return builtin.IntegerNumberType{
				Value: leftVal.Value + right.(builtin.IntegerNumberType).Value,
			}, nil
		case builtin.StringType:
			return builtin.StringType{
				Value: leftVal.Value + right.(builtin.StringType).Value,
			}, nil
		}

	case subOp:
		switch leftVal := left.(type) {
		case builtin.RealNumberType:
			return builtin.RealNumberType{
				Value: leftVal.Value - right.(builtin.RealNumberType).Value,
			}, nil
		case builtin.IntegerNumberType:
			return builtin.IntegerNumberType{
				Value: leftVal.Value - right.(builtin.IntegerNumberType).Value,
			}, nil
		}
	case mulOp:
		switch leftVal := left.(type) {
		case builtin.RealNumberType:
			return builtin.RealNumberType{
				Value: leftVal.Value * right.(builtin.RealNumberType).Value,
			}, nil
		case builtin.IntegerNumberType:
			return builtin.IntegerNumberType{
				Value: leftVal.Value * right.(builtin.IntegerNumberType).Value,
			}, nil
		}
	case divOp:
		switch leftVal := left.(type) {
		case builtin.RealNumberType:
			rightVal := right.(builtin.RealNumberType).Value
			if rightVal == 0 {
				return builtin.NoneType{}, util.RuntimeError("ділення на нуль")
			}

			return builtin.RealNumberType{
				Value: leftVal.Value / right.(builtin.RealNumberType).Value,
			}, nil
		case builtin.IntegerNumberType:
			rightVal := right.(builtin.IntegerNumberType).Value
			if rightVal == 0 {
				return builtin.NoneType{}, util.RuntimeError("ділення на нуль")
			}

			return builtin.RealNumberType{
				Value: float64(leftVal.Value) / right.(builtin.RealNumberType).Value,
			}, nil
		}

	default:
		return builtin.NoneType{}, util.RuntimeError("невідомий оператор")
	}

	return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opTypeNames[opType], left.TypeName(), right.TypeName(),
	))
}
