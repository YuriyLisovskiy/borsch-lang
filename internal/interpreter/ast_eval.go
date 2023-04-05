package interpreter

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/internal/builtin"
	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/internal/utilities"
)

type Scope map[string]types2.Object

func (node *Package) Evaluate(state State) (types2.Object, error) {
	state.Context().PushScope(map[string]types2.Object{})
	result := node.Stmts.Evaluate(state, false, false)
	if result.Err != nil {
		state.Trace(node.Stmts, "<пакет>")
	}

	return result.Value, result.Err
}

// Evaluate executes block of statements.
// Returns (result value, force stop flag, error)
func (node *BlockStmts) Evaluate(state State, inFunction, inLoop bool) StmtResult {
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

	return StmtResult{Value: types2.Nil}
}

func (node *Expression) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	if node.LogicalAnd != nil {
		return node.LogicalAnd.Evaluate(state, valueToSet)
	}

	panic("unreachable")
}

func (node *Assignment) Evaluate(state State) (types2.Object, error) {
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
func (node *LogicalAnd) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	return evalBinaryOperator(state, valueToSet, types2.And, node.LogicalOr, node.Next)
}

func (node *LogicalOr) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	return evalBinaryOperator(state, valueToSet, types2.Or, node.LogicalNot, node.Next)
}

func (node *LogicalNot) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	if node.Comparison != nil {
		return node.Comparison.Evaluate(state, valueToSet)
	}

	if node.Next != nil {
		value, err := node.Next.Evaluate(state, nil)
		if err != nil {
			return nil, err
		}

		return types2.Not(state.Context(), value)
	}

	panic("unreachable")
}

func (node *Comparison) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	switch node.Op {
	case ">=":
		return evalBinaryOperator(state, valueToSet, types2.GreaterOrEquals, node.BitwiseOr, node.Next)
	case ">":
		return evalBinaryOperator(state, valueToSet, types2.Greater, node.BitwiseOr, node.Next)
	case "<=":
		return evalBinaryOperator(state, valueToSet, types2.LessOrEquals, node.BitwiseOr, node.Next)
	case "<":
		return evalBinaryOperator(state, valueToSet, types2.Less, node.BitwiseOr, node.Next)
	case "==":
		return evalBinaryOperator(state, valueToSet, types2.Equals, node.BitwiseOr, node.Next)
	case "!=":
		return evalBinaryOperator(state, valueToSet, types2.NotEquals, node.BitwiseOr, node.Next)
	default:
		return node.BitwiseOr.Evaluate(state, valueToSet)
	}
}

func (node *BitwiseOr) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	return evalBinaryOperator(state, valueToSet, types2.BitwiseOr, node.BitwiseXor, node.Next)
}

func (node *BitwiseXor) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	return evalBinaryOperator(state, valueToSet, types2.BitwiseXor, node.BitwiseAnd, node.Next)
}

func (node *BitwiseAnd) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	return evalBinaryOperator(state, valueToSet, types2.BitwiseAnd, node.BitwiseShift, node.Next)
}

func (node *BitwiseShift) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	switch node.Op {
	case "<<":
		return evalBinaryOperator(state, valueToSet, types2.ShiftLeft, node.Addition, node.Next)
	case ">>":
		return evalBinaryOperator(state, valueToSet, types2.ShiftRight, node.Addition, node.Next)
	default:
		return node.Addition.Evaluate(state, valueToSet)
	}
}

func (node *Addition) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	switch node.Op {
	case "+":
		return evalBinaryOperator(state, valueToSet, types2.Add, node.MultiplicationOrMod, node.Next)
	case "-":
		return evalBinaryOperator(state, valueToSet, types2.Sub, node.MultiplicationOrMod, node.Next)
	default:
		return node.MultiplicationOrMod.Evaluate(state, valueToSet)
	}
}

func (node *MultiplicationOrMod) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	switch node.Op {
	case "/":
		return evalBinaryOperator(state, valueToSet, types2.Div, node.Unary, node.Next)
	case "*":
		return evalBinaryOperator(state, valueToSet, types2.Mul, node.Unary, node.Next)
	case "%":
		return evalBinaryOperator(state, valueToSet, types2.Mod, node.Unary, node.Next)
	default:
		return node.Unary.Evaluate(state, valueToSet)
	}
}

func (node *Unary) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	switch node.Op {
	case "+":
		return evalUnaryOperator(state, types2.Positive, node.Next)
	case "-":
		return evalUnaryOperator(state, types2.Negate, node.Next)
	case "~":
		return evalUnaryOperator(state, types2.Invert, node.Next)
	default:
		return node.Exponent.Evaluate(state, valueToSet)
	}
}

func (node *Exponent) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	return evalBinaryOperator(state, valueToSet, types2.Pow, node.Primary, node.Next)
}

func (node *Primary) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
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

		return node.Literal.Evaluate(state, valueToSet)
	}

	if node.AttributeAccess != nil {
		return node.AttributeAccess.Evaluate(state, valueToSet, nil)
	}

	if node.LambdaDef != nil {
		return node.LambdaDef.Evaluate(state)
	}

	panic("unreachable")
}

func (node *Literal) Evaluate(state State, valueToSet types2.Object) (types2.Object, error) {
	if node.Nil {
		return types2.Nil, nil
	}

	if node.Integer != nil {
		return types2.IntFromString(*node.Integer, 0)
	}

	if node.Real != nil {
		return types2.RealFromString(*node.Real)
	}

	if node.Bool != nil {
		return types2.NewBool(bool(*node.Bool)), nil
	}

	if node.StringValue != nil {
		return types2.String(*node.StringValue), nil
	}

	if node.MultilineString != nil {
		return types2.String(*node.MultilineString), nil
	}

	if node.List != nil {
		list := types2.NewList()
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
		return types2.NewList(), nil
	}

	// if node.SubExpression != nil {
	// 	if valueToSet != nil {
	// 		return nil, utilities.SyntaxError("неможливо записати значення у вираз")
	// 	}
	//
	// 	return node.SubExpression.Evaluate(state, valueToSet)
	// }

	// if node.Dictionary != nil {
	// 	dict := types.NewDictionaryInstance()
	// 	for _, entry := range node.Dictionary {
	// 		key, value, err := entry.Evaluate(state)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	//
	// 		if err := dict.SetElement(key, value); err != nil {
	// 			return nil, err
	// 		}
	// 	}
	//
	// 	return dict, nil
	// }

	// if node.EmptyDictionary {
	// 	return types.NewDictionaryInstance(), nil
	// }

	panic("unreachable")
}

func (node *DictionaryEntry) Evaluate(state State) (types2.Object, types2.Object, error) {
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

func (node *AttributeAccess) Evaluate(state State, valueToSet, prevValue types2.Object) (
	types2.Object,
	error,
) {
	if node.IdentOrCall == nil {
		panic("unreachable")
	}

	if valueToSet != nil {
		// set
		var currentValue types2.Object
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

func (node *IdentOrCall) Evaluate(state State, valueToSet types2.Object, prevValue types2.Object) (
	types2.Object,
	error,
) {
	if valueToSet != nil {
		// set
		var variable types2.Object
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
				variable, err = getCurrentValue(state.Context(), prevValue, node.Ident.String())
			} else {
				variable, err = setCurrentValue(state.Context(), prevValue, node.Ident.String(), valueToSet)
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
				return setCurrentValue(state.Context(), prevValue, node.Ident.String(), variable)
			}
		}

		return variable, nil
	}

	// get
	var variable types2.Object
	var err error = nil
	if node.Call != nil {
		variable, err = node.callFunction(state, prevValue)
		if err != nil {
			return nil, err
		}
	} else if node.Ident != nil {
		variable, err = getCurrentValue(state.Context(), prevValue, node.Ident.String())
		if err != nil {
			state.Trace(node, "")
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

func (node *IdentOrCall) callFunction(state State, prevValue types2.Object) (types2.Object, error) {
	ctx := state.Context()
	variable, err := getCurrentValue(ctx, prevValue, node.Call.Ident.String())
	if err != nil {
		state.Trace(node, "")
		return nil, err
	}

	variable, err = node.Call.Evaluate(state, variable)
	if err != nil {
		if _, ok := err.(utilities.CallError); !ok {
			err = utilities.NewCallError(err, string(node.Call.Ident))
		}

		return nil, err
	}

	return variable, nil
}

func (node *SlicingOrSubscription) Evaluate(
	state State,
	variable types2.Object,
	valueToSet types2.Object,
) (types2.Object, error) {
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

func (node *LambdaDef) Evaluate(state State) (types2.Object, error) {
	arguments, err := node.ParametersSet.Evaluate(state)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(state, node.ReturnTypes)
	if err != nil {
		return nil, err
	}

	lambda := types2.LambdaNew(
		builtin.LambdaSignature,
		state.Package().(*types2.Package),
		arguments,
		returnTypes,
		func(ctx types2.Context, args types2.Tuple, kwargs types2.StringDict) (types2.Object, error) {
			return node.Body.Evaluate(state.NewChild().WithContext(ctx))
		},
	)

	if node.InstantCall {
		return node.evalInstantCall(state, lambda)
	}

	return lambda, nil
}

func (node *LambdaDef) evalInstantCall(state State, function *types2.Method) (types2.Object, error) {
	var args types2.Tuple
	if len(node.InstantCallArguments) != 0 {
		if err := updateArgs(state, node.InstantCallArguments, &args); err != nil {
			return nil, err
		}
	}

	return types2.Call(state.Context(), function, args)
}
