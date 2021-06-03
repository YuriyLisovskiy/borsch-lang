package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func compareNones(left, right builtin.NoneType, opType Operator) (builtin.ValueType, error) {
	switch opType {
	case equalsOp:
		return builtin.BoolType{Value: true}, nil
	case notEqualsOp:
		return builtin.BoolType{Value: false}, nil
	case greaterOp, greaterOrEqualsOp, lessOp, lessOrEqualsOp:
		return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"оператор %s невизначений для значень типів '%s' та '%s'",
			opType.Description(), left.TypeName(), right.TypeName(),
		))
	default:
		return builtin.NoneType{}, util.RuntimeError("невідомий оператор")
	}
}

func compareReals(left, right builtin.RealNumberType, opType Operator) (builtin.ValueType, error) {
	switch opType {
	case equalsOp:
		return builtin.BoolType{Value: left.Value == right.Value}, nil
	case notEqualsOp:
		return builtin.BoolType{Value: left.Value != right.Value}, nil
	case greaterOp:
		return builtin.BoolType{Value: left.Value > right.Value}, nil
	case greaterOrEqualsOp:
		return builtin.BoolType{Value: left.Value >= right.Value}, nil
	case lessOp:
		return builtin.BoolType{Value: left.Value < right.Value}, nil
	case lessOrEqualsOp:
		return builtin.BoolType{Value: left.Value <= right.Value}, nil
	default:
		return builtin.NoneType{}, util.RuntimeError("невідомий оператор")
	}
}

func compareIntegers(left, right builtin.IntegerNumberType, opType Operator) (builtin.ValueType, error) {
	switch opType {
	case equalsOp:
		return builtin.BoolType{Value: left.Value == right.Value}, nil
	case notEqualsOp:
		return builtin.BoolType{Value: left.Value != right.Value}, nil
	case greaterOp:
		return builtin.BoolType{Value: left.Value > right.Value}, nil
	case greaterOrEqualsOp:
		return builtin.BoolType{Value: left.Value >= right.Value}, nil
	case lessOp:
		return builtin.BoolType{Value: left.Value < right.Value}, nil
	case lessOrEqualsOp:
		return builtin.BoolType{Value: left.Value <= right.Value}, nil
	default:
		return builtin.NoneType{}, util.RuntimeError("невідомий оператор")
	}
}

func compareStrings(left, right builtin.StringType, opType Operator) (builtin.ValueType, error) {
	switch opType {
	case equalsOp:
		return builtin.BoolType{Value: left.Value == right.Value}, nil
	case notEqualsOp:
		return builtin.BoolType{Value: left.Value != right.Value}, nil
	case greaterOp:
		return builtin.BoolType{Value: left.Value > right.Value}, nil
	case greaterOrEqualsOp:
		return builtin.BoolType{Value: left.Value >= right.Value}, nil
	case lessOp:
		return builtin.BoolType{Value: left.Value < right.Value}, nil
	case lessOrEqualsOp:
		return builtin.BoolType{Value: left.Value <= right.Value}, nil
	default:
		return builtin.NoneType{}, util.RuntimeError("невідомий оператор")
	}
}

func compareBooleans(left, right builtin.BoolType, opType Operator) (builtin.ValueType, error) {
	switch opType {
	case equalsOp:
		return builtin.BoolType{Value: left.Value == right.Value}, nil
	case notEqualsOp:
		return builtin.BoolType{Value: left.Value != right.Value}, nil
	case greaterOp, greaterOrEqualsOp, lessOp, lessOrEqualsOp:
		return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"оператор %s невизначений для значень типів '%s' та '%s'",
			opType.Description(), left.TypeName(), right.TypeName(),
		))
	default:
		return builtin.NoneType{}, util.RuntimeError("невідомий оператор")
	}
}

func (e *Interpreter) executeComparisonOp(
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
				opTypeNames[opType], left.TypeName(), right.TypeName(),
			),
		)
	}

	switch leftV := left.(type) {
	case builtin.NoneType:
		return compareNones(leftV, right.(builtin.NoneType), opType)
	case builtin.RealNumberType:
		return compareReals(leftV, right.(builtin.RealNumberType), opType)
	case builtin.IntegerNumberType:
		return compareIntegers(leftV, right.(builtin.IntegerNumberType), opType)
	case builtin.StringType:
		return compareStrings(leftV, right.(builtin.StringType), opType)
	case builtin.BoolType:
		return compareBooleans(leftV, right.(builtin.BoolType), opType)
	}

	return builtin.NoneType{}, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opTypeNames[opType], left.TypeName(), right.TypeName(),
	))
}
