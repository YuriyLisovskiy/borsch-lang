package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func compareNones(left, right types.NoneType, opType Operator) (types.ValueType, error) {
	switch opType {
	case equalsOp:
		return types.BoolType{Value: true}, nil
	case notEqualsOp:
		return types.BoolType{Value: false}, nil
	case greaterOp, greaterOrEqualsOp, lessOp, lessOrEqualsOp:
		return nil, util.RuntimeError(fmt.Sprintf(
			"оператор %s невизначений для значень типів '%s' та '%s'",
			opType.Description(), left.TypeName(), right.TypeName(),
		))
	default:
		return nil, util.RuntimeError("невідомий оператор")
	}
}

func compareReals(left, right types.RealType, opType Operator) (types.ValueType, error) {
	switch opType {
	case equalsOp:
		return types.BoolType{Value: left.Value == right.Value}, nil
	case notEqualsOp:
		return types.BoolType{Value: left.Value != right.Value}, nil
	case greaterOp:
		return types.BoolType{Value: left.Value > right.Value}, nil
	case greaterOrEqualsOp:
		return types.BoolType{Value: left.Value >= right.Value}, nil
	case lessOp:
		return types.BoolType{Value: left.Value < right.Value}, nil
	case lessOrEqualsOp:
		return types.BoolType{Value: left.Value <= right.Value}, nil
	default:
		return nil, util.RuntimeError("невідомий оператор")
	}
}

func compareIntegers(left, right types.IntegerType, opType Operator) (types.ValueType, error) {
	switch opType {
	case equalsOp:
		return types.BoolType{Value: left.Value == right.Value}, nil
	case notEqualsOp:
		return types.BoolType{Value: left.Value != right.Value}, nil
	case greaterOp:
		return types.BoolType{Value: left.Value > right.Value}, nil
	case greaterOrEqualsOp:
		return types.BoolType{Value: left.Value >= right.Value}, nil
	case lessOp:
		return types.BoolType{Value: left.Value < right.Value}, nil
	case lessOrEqualsOp:
		return types.BoolType{Value: left.Value <= right.Value}, nil
	default:
		return nil, util.RuntimeError("невідомий оператор")
	}
}

func compareStrings(left, right types.StringType, opType Operator) (types.ValueType, error) {
	switch opType {
	case equalsOp:
		return types.BoolType{Value: left.Value == right.Value}, nil
	case notEqualsOp:
		return types.BoolType{Value: left.Value != right.Value}, nil
	case greaterOp:
		return types.BoolType{Value: left.Value > right.Value}, nil
	case greaterOrEqualsOp:
		return types.BoolType{Value: left.Value >= right.Value}, nil
	case lessOp:
		return types.BoolType{Value: left.Value < right.Value}, nil
	case lessOrEqualsOp:
		return types.BoolType{Value: left.Value <= right.Value}, nil
	default:
		return nil, util.RuntimeError("невідомий оператор")
	}
}

func compareBooleans(left, right types.BoolType, opType Operator) (types.ValueType, error) {
	switch opType {
	case equalsOp:
		return types.BoolType{Value: left.Value == right.Value}, nil
	case notEqualsOp:
		return types.BoolType{Value: left.Value != right.Value}, nil
	case greaterOp, greaterOrEqualsOp, lessOp, lessOrEqualsOp:
		return nil, util.RuntimeError(fmt.Sprintf(
			"оператор %s невизначений для значень типів '%s' та '%s'",
			opType.Description(), left.TypeName(), right.TypeName(),
		))
	default:
		return nil, util.RuntimeError("невідомий оператор")
	}
}

func (i *Interpreter) executeComparisonOp(
	leftNode ast.ExpressionNode, rightNode ast.ExpressionNode, opType Operator,
	rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	left, _, err := i.executeNode(leftNode, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	right, _, err := i.executeNode(rightNode, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	if left.TypeHash() != right.TypeHash() {
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				opTypeNames[opType], left.TypeName(), right.TypeName(),
			),
		)
	}

	switch leftV := left.(type) {
	case types.NoneType:
		return compareNones(leftV, right.(types.NoneType), opType)
	case types.RealType:
		return compareReals(leftV, right.(types.RealType), opType)
	case types.IntegerType:
		return compareIntegers(leftV, right.(types.IntegerType), opType)
	case types.StringType:
		return compareStrings(leftV, right.(types.StringType), opType)
	case types.BoolType:
		return compareBooleans(leftV, right.(types.BoolType), opType)
	}

	return nil, util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opTypeNames[opType], left.TypeName(), right.TypeName(),
	))
}
