package interpreter

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type Scope map[string]common.Value

func (node *Package) Evaluate(state common.State) (common.Value, error) {
	state.GetContext().PushScope(map[string]common.Value{})
	result := node.Stmts.Evaluate(state, false, false)
	if result.Err != nil {
		state.Trace(node.Stmts, "<пакет>")
	}

	return result.Value, result.Err
}

// Evaluate executes block of statements.
// Returns (result value, force stop flag, error)
func (node *BlockStmts) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	node.stmtPos = 0
	for _, stmt := range node.Stmts {
		result := stmt.Evaluate(state, inFunction, inLoop)
		if result.Interrupt() {
			if callErr, ok := result.Err.(utilities.CallError); ok {
				state.Trace(stmt, callErr.Function())
				result.Err = callErr.Original()
			}

			return result
		}

		node.stmtPos++
	}

	return StmtResult{Value: types.NewNilInstance()}
}

func (node *Expression) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	if node.LogicalAnd != nil {
		return node.LogicalAnd.Evaluate(state, valueToSet)
	}

	panic("unreachable")
}

func (node *Assignment) Evaluate(state common.State) (common.Value, error) {
	if len(node.Next) == 0 {
		return node.Expressions[0].Evaluate(state, nil)
	}

	value, err := unpack(state, node.Expressions, node.Next)
	if err != nil {
		if _, ok := err.(utilities.CallError); !ok {
			state.Trace(node, "")
		}

		return nil, err
	}

	return value, nil
}

// Evaluate executes LogicalAnd operation.
// If `valueToSet` is nil, return variable or value from context,
// set a new value or return an error otherwise.
func (node *LogicalAnd) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.AndOp.Name(), node.LogicalOr, node.Next)
}

func (node *LogicalOr) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.OrOp.Name(), node.LogicalNot, node.Next)
}

func (node *LogicalNot) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	if node.Comparison != nil {
		return node.Comparison.Evaluate(state, valueToSet)
	}

	if node.Next != nil {
		value, err := node.Next.Evaluate(state, nil)
		if err != nil {
			return nil, err
		}

		opName := common.NotOp.Name()
		operatorFunc, err := value.GetOperator(opName)
		if err != nil {
			return nil, err
		}

		return types.CallAttribute(state, value, operatorFunc, opName, nil, nil, true)
	}

	panic("unreachable")
}

func (node *Comparison) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch node.Op {
	case ">=":
		return evalBinaryOperator(state, valueToSet, common.GreaterOrEqualsOp.Name(), node.BitwiseOr, node.Next)
	case ">":
		return evalBinaryOperator(state, valueToSet, common.GreaterOp.Name(), node.BitwiseOr, node.Next)
	case "<=":
		return evalBinaryOperator(state, valueToSet, common.LessOrEqualsOp.Name(), node.BitwiseOr, node.Next)
	case "<":
		return evalBinaryOperator(state, valueToSet, common.LessOp.Name(), node.BitwiseOr, node.Next)
	case "==":
		return evalBinaryOperator(state, valueToSet, common.EqualsOp.Name(), node.BitwiseOr, node.Next)
	case "!=":
		return evalBinaryOperator(state, valueToSet, common.NotEqualsOp.Name(), node.BitwiseOr, node.Next)
	default:
		return node.BitwiseOr.Evaluate(state, valueToSet)
	}
}

func (node *BitwiseOr) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.BitwiseOrOp.Name(), node.BitwiseXor, node.Next)
}

func (node *BitwiseXor) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.BitwiseXorOp.Name(), node.BitwiseAnd, node.Next)
}

func (node *BitwiseAnd) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.BitwiseAndOp.Name(), node.BitwiseShift, node.Next)
}

func (node *BitwiseShift) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch node.Op {
	case "<<":
		return evalBinaryOperator(state, valueToSet, common.BitwiseLeftShiftOp.Name(), node.Addition, node.Next)
	case ">>":
		return evalBinaryOperator(state, valueToSet, common.BitwiseRightShiftOp.Name(), node.Addition, node.Next)
	default:
		return node.Addition.Evaluate(state, valueToSet)
	}
}

func (node *Addition) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch node.Op {
	case "+":
		return evalBinaryOperator(state, valueToSet, common.AddOp.Name(), node.MultiplicationOrMod, node.Next)
	case "-":
		return evalBinaryOperator(state, valueToSet, common.SubOp.Name(), node.MultiplicationOrMod, node.Next)
	default:
		return node.MultiplicationOrMod.Evaluate(state, valueToSet)
	}
}

func (node *MultiplicationOrMod) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch node.Op {
	case "/":
		return evalBinaryOperator(state, valueToSet, common.DivOp.Name(), node.Unary, node.Next)
	case "*":
		return evalBinaryOperator(state, valueToSet, common.MulOp.Name(), node.Unary, node.Next)
	case "%":
		return evalBinaryOperator(state, valueToSet, common.ModuloOp.Name(), node.Unary, node.Next)
	default:
		return node.Unary.Evaluate(state, valueToSet)
	}
}

func (node *Unary) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch node.Op {
	case "+":
		return evalUnaryOperator(state, common.UnaryPlus.Name(), node.Next)
	case "-":
		return evalUnaryOperator(state, common.UnaryMinus.Name(), node.Next)
	case "~":
		return evalUnaryOperator(state, common.UnaryBitwiseNotOp.Name(), node.Next)
	default:
		return node.Exponent.Evaluate(state, valueToSet)
	}
}

func (node *Exponent) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.PowOp.Name(), node.Primary, node.Next)
}

func (node *Primary) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	if node.SubExpression != nil {
		if valueToSet != nil {
			return nil, utilities.SyntaxError("неможливо записати значення у вираз")
		}

		return node.SubExpression.Evaluate(state, valueToSet)
	}

	if node.Literal != nil {
		if valueToSet != nil {
			return nil, utilities.SyntaxError("неможливо встановити значення у літерал")
		}

		return node.Literal.Evaluate(state)
	}

	if node.AttributeAccess != nil {
		return node.AttributeAccess.Evaluate(state, valueToSet, nil)
	}

	if node.LambdaDef != nil {
		return node.LambdaDef.Evaluate(state)
	}

	panic("unreachable")
}

func (node *Literal) Evaluate(state common.State) (common.Value, error) {
	if node.Nil {
		return types.NewNilInstance(), nil
	}

	if node.Integer != nil {
		return types.NewIntegerInstance(*node.Integer), nil
	}

	if node.Real != nil {
		return types.NewRealInstance(*node.Real), nil
	}

	if node.Bool != nil {
		return types.NewBoolInstance(bool(*node.Bool)), nil
	}

	if node.StringValue != nil {
		return types.NewStringInstance(*node.StringValue), nil
	}

	if node.MultilineString != nil {
		return types.NewStringInstance(*node.MultilineString), nil
	}

	if node.List != nil {
		list := types.NewListInstance()
		for _, expr := range node.List {
			value, err := expr.Evaluate(state, nil)
			if err != nil {
				return nil, err
			}

			list.Values = append(list.Values, value)
		}

		return list, nil
	}

	if node.EmptyList {
		return types.NewListInstance(), nil
	}

	if node.Dictionary != nil {
		dict := types.NewDictionaryInstance()
		for _, entry := range node.Dictionary {
			key, value, err := entry.Evaluate(state)
			if err != nil {
				return nil, err
			}

			if err := dict.SetElement(key, value); err != nil {
				return nil, err
			}
		}

		return dict, nil
	}

	if node.EmptyDictionary {
		return types.NewDictionaryInstance(), nil
	}

	panic("unreachable")
}

func (node *DictionaryEntry) Evaluate(state common.State) (common.Value, common.Value, error) {
	key, err := node.Key.Evaluate(state, nil)
	if err != nil {
		return nil, nil, err
	}

	value, err := node.Value.Evaluate(state, nil)
	if err != nil {
		return nil, nil, err
	}

	return key, value, nil
}

func (node *AttributeAccess) Evaluate(state common.State, valueToSet, prevValue common.Value) (common.Value, error) {
	if node.IdentOrCall == nil {
		panic("unreachable")
	}

	if valueToSet != nil {
		// set
		var currentValue common.Value
		var err error
		if node.AttributeAccess != nil {
			currentValue, err = node.IdentOrCall.Evaluate(state, nil, prevValue)
			if err != nil {
				return nil, err
			}

			currentValue, err = node.AttributeAccess.Evaluate(state, valueToSet, currentValue)
		} else {
			currentValue, err = node.IdentOrCall.Evaluate(state, valueToSet, prevValue)
		}

		if err != nil {
			return nil, err
		}

		return currentValue, nil
	}

	// get
	currentValue, err := node.IdentOrCall.Evaluate(state, valueToSet, prevValue)
	if err != nil {
		return nil, err
	}

	if node.AttributeAccess != nil {
		return node.AttributeAccess.Evaluate(state, valueToSet, currentValue)
	}

	return currentValue, err
}

func (node *IdentOrCall) Evaluate(state common.State, valueToSet common.Value, prevValue common.Value) (
	common.Value,
	error,
) {
	if valueToSet != nil {
		// set
		var variable common.Value
		var err error = nil
		if node.Call != nil {
			if node.SlicingOrSubscription == nil {
				return nil, errors.New("неможливо присвоїти значення виклику функції")
			}

			variable, err = node.callFunction(state, prevValue)
			if err != nil {
				return nil, err
			}
		} else if node.Ident != nil {
			if node.SlicingOrSubscription != nil {
				variable, err = getCurrentValue(state.GetContext(), prevValue, node.Ident.String())
			} else {
				variable, err = setCurrentValue(state.GetContext(), prevValue, node.Ident.String(), valueToSet)
			}

			if err != nil {
				return nil, err
			}
		} else {
			panic("unreachable")
		}

		if node.SlicingOrSubscription != nil {
			variable, err = node.SlicingOrSubscription.Evaluate(state, variable, valueToSet)
			if err != nil {
				return nil, err
			}

			if node.Ident != nil {
				return setCurrentValue(state.GetContext(), prevValue, node.Ident.String(), variable)
			}
		}

		return variable, nil
	}

	// get
	var variable common.Value
	var err error = nil
	if node.Call != nil {
		variable, err = node.callFunction(state, prevValue)
		if err != nil {
			return nil, err
		}
	} else if node.Ident != nil {
		variable, err = getCurrentValue(state.GetContext(), prevValue, node.Ident.String())
		if err != nil {
			return nil, err
		}
	} else {
		panic("unreachable")
	}

	if node.SlicingOrSubscription != nil {
		return node.SlicingOrSubscription.Evaluate(state, variable, nil)
	}

	return variable, nil
}

func (node *IdentOrCall) callFunction(state common.State, prevValue common.Value) (common.Value, error) {
	ctx := state.GetContext()
	variable, err := getCurrentValue(ctx, prevValue, node.Call.Ident.String())
	if err != nil {
		return nil, err
	}

	isLambda := false
	variable, err = node.Call.Evaluate(state, variable, prevValue, &isLambda)
	if err != nil {
		return nil, err
	}

	return variable, nil
}

func (node *SlicingOrSubscription) Evaluate(
	state common.State,
	variable common.Value,
	valueToSet common.Value,
) (common.Value, error) {
	if valueToSet != nil {
		// set
		rangesLen := len(node.Ranges)
		if rangesLen != 0 && node.Ranges[rangesLen-1].RightBound != nil {
			return nil, errors.New("неможливо присвоїти значення зрізу")
		}

		if len(node.Ranges) != 0 {
			return evalSlicingOperation(state, variable, node.Ranges, valueToSet)
		}

		return variable, nil
	}

	// get
	if len(node.Ranges) != 0 {
		return evalSlicingOperation(state, variable, node.Ranges, nil)
	}

	return variable, nil
}

func (node *LambdaDef) Evaluate(state common.State) (common.Value, error) {
	arguments, err := node.ParametersSet.Evaluate(state)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(state, node.ReturnTypes)
	if err != nil {
		return nil, err
	}

	lambda := types.NewFunctionInstance(
		common.LambdaSignature,
		arguments,
		func(state common.State, _ *[]common.Value, kwargs *map[string]common.Value) (common.Value, error) {
			return node.Body.Evaluate(state)
		},
		returnTypes,
		false,
		state.GetCurrentPackage().(*types.PackageInstance),
		"",
	)

	if node.InstantCall {
		return node.evalInstantCall(state, lambda)
	}

	return lambda, nil
}

func (node *LambdaDef) evalInstantCall(state common.State, function *types.FunctionInstance) (common.Value, error) {
	var args []common.Value
	if len(node.InstantCallArguments) != 0 {
		if err := updateArgs(state, node.InstantCallArguments, &args); err != nil {
			return nil, err
		}
	}

	return types.Call(state, function, &args, nil)
}
