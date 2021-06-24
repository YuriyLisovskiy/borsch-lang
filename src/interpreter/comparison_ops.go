package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)
//
//func compareNils(left, right types.NilType, opType Operator) (types.ValueType, error) {
//	switch opType {
//	case equalsOp:
//		return types.BoolType{Value: true}, nil
//	case notEqualsOp:
//		return types.BoolType{Value: false}, nil
//	case greaterOp, greaterOrEqualsOp, lessOp, lessOrEqualsOp:
//		return nil, util.RuntimeError(fmt.Sprintf(
//			"оператор %s невизначений для значень типів '%s' та '%s'",
//			opType.Description(), left.TypeName(), right.TypeName(),
//		))
//	default:
//		return nil, util.RuntimeError("невідомий оператор")
//	}
//}
//
//func compareReals(left, right types.RealType, opType Operator) (types.ValueType, error) {
//	switch opType {
//	case equalsOp:
//		return types.BoolType{Value: left.Value == right.Value}, nil
//	case notEqualsOp:
//		return types.BoolType{Value: left.Value != right.Value}, nil
//	case greaterOp:
//		return types.BoolType{Value: left.Value > right.Value}, nil
//	case greaterOrEqualsOp:
//		return types.BoolType{Value: left.Value >= right.Value}, nil
//	case lessOp:
//		return types.BoolType{Value: left.Value < right.Value}, nil
//	case lessOrEqualsOp:
//		return types.BoolType{Value: left.Value <= right.Value}, nil
//	default:
//		return nil, util.RuntimeError("невідомий оператор")
//	}
//}
//
//func compareIntegers(left, right types.IntegerType, opType Operator) (types.ValueType, error) {
//	switch opType {
//	case equalsOp:
//		return types.BoolType{Value: left.Value == right.Value}, nil
//	case notEqualsOp:
//		return types.BoolType{Value: left.Value != right.Value}, nil
//	case greaterOp:
//		return types.BoolType{Value: left.Value > right.Value}, nil
//	case greaterOrEqualsOp:
//		return types.BoolType{Value: left.Value >= right.Value}, nil
//	case lessOp:
//		return types.BoolType{Value: left.Value < right.Value}, nil
//	case lessOrEqualsOp:
//		return types.BoolType{Value: left.Value <= right.Value}, nil
//	default:
//		return nil, util.RuntimeError("невідомий оператор")
//	}
//}
//
//func compareStrings(left, right types.StringType, opType Operator) (types.ValueType, error) {
//	switch opType {
//	case equalsOp:
//		return types.BoolType{Value: left.Value == right.Value}, nil
//	case notEqualsOp:
//		return types.BoolType{Value: left.Value != right.Value}, nil
//	case greaterOp:
//		return types.BoolType{Value: left.Value > right.Value}, nil
//	case greaterOrEqualsOp:
//		return types.BoolType{Value: left.Value >= right.Value}, nil
//	case lessOp:
//		return types.BoolType{Value: left.Value < right.Value}, nil
//	case lessOrEqualsOp:
//		return types.BoolType{Value: left.Value <= right.Value}, nil
//	default:
//		return nil, util.RuntimeError("невідомий оператор")
//	}
//}
//
//func compareBooleans(left, right types.BoolType, opType Operator) (types.ValueType, error) {
//	switch opType {
//	case equalsOp:
//		return types.BoolType{Value: left.Value == right.Value}, nil
//	case notEqualsOp:
//		return types.BoolType{Value: left.Value != right.Value}, nil
//	case greaterOp, greaterOrEqualsOp, lessOp, lessOrEqualsOp:
//		return nil, util.RuntimeError(fmt.Sprintf(
//			"оператор %s невизначений для значень типів '%s' та '%s'",
//			opType.Description(), left.TypeName(), right.TypeName(),
//		))
//	default:
//		return nil, util.RuntimeError("невідомий оператор")
//	}
//}

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

	//if left.TypeHash() != right.TypeHash() {
	//	if left.TypeHash() != types.NilTypeHash && right.TypeHash() != types.NilTypeHash {
	//		return nil, util.RuntimeError(
	//			fmt.Sprintf(
	//				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
	//				opTypeNames[opType], left.TypeName(), right.TypeName(),
	//			),
	//		)
	//	}
	//}
	//
	switch left.(type) {
	case types.NilType, types.BoolType:
		res, err := left.CompareTo(right)
		if err != nil {
			return nil, util.RuntimeError(fmt.Sprintf(err.Error(), opTypeNames[opType]))
		}

		switch opType {
		case equalsOp:
			return types.BoolType{Value: res == 0}, nil
		case notEqualsOp:
			return types.BoolType{Value: res != 0}, nil
		case greaterOp, greaterOrEqualsOp, lessOp, lessOrEqualsOp:
			return nil, util.RuntimeError(fmt.Sprintf(
				"оператор %s невизначений для значень типів '%s' та '%s'",
				opType.Description(), left.TypeName(), right.TypeName(),
			))
		default:
			return nil, util.RuntimeError("невідомий оператор")
		}
	//case types.RealType:
	//	return compareReals(leftV, right.(types.RealType), opType)
	//case types.IntegerType:
	//	return compareIntegers(leftV, right.(types.IntegerType), opType)
	//case types.StringType:
	//	return compareStrings(leftV, right.(types.StringType), opType)
	//case types.BoolType:
	//	res, err := leftV.CompareTo(right)
	//	if err != nil {
	//		return nil, util.RuntimeError(fmt.Sprintf(err.Error(), opTypeNames[opType]))
	//	}
	//
	//	switch opType {
	//	case equalsOp:
	//		return types.BoolType{Value: res == 0}, nil
	//	case notEqualsOp:
	//		return types.BoolType{Value: res != 0}, nil
	//	case greaterOp, greaterOrEqualsOp, lessOp, lessOrEqualsOp:
	//		return nil, util.RuntimeError(fmt.Sprintf(
	//			"оператор %s невизначений для значень типів '%s' та '%s'",
	//			opType.Description(), left.TypeName(), right.TypeName(),
	//		))
	//	default:
	//		return nil, util.RuntimeError("невідомий оператор")
	//	}
	default:
		res, err := left.CompareTo(right)
		if err != nil {
			return nil, util.RuntimeError(fmt.Sprintf(err.Error(), opTypeNames[opType]))
		}

		switch opType {
		case equalsOp:
			return types.BoolType{Value: res == 0}, nil
		case notEqualsOp:
			return types.BoolType{Value: res != 0}, nil
		case greaterOp:
			return types.BoolType{Value: res == 1}, nil
		case greaterOrEqualsOp:
			return types.BoolType{Value: res == 0 || res == 1}, nil
		case lessOp:
			return types.BoolType{Value: res == -1}, nil
		case lessOrEqualsOp:
			return types.BoolType{Value: res == 0 || res == -1}, nil
		default:
			return nil, util.RuntimeError("невідомий оператор")
		}
	}

	//return nil, util.RuntimeError(fmt.Sprintf(
	//	"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
	//	opTypeNames[opType], left.TypeName(), right.TypeName(),
	//))
}
