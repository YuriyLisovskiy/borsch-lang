package interpreter

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Scope map[string]common.Value

func (p *Package) Evaluate(state common.State) (common.Value, error) {
	state.GetContext().PushScope(Scope{})
	for _, stmt := range p.Stmts {
		stmtState := stmt.Evaluate(state, false, false)
		if stmtState.Err != nil {
			err := stmtState.Err
			state.GetInterpreter().Trace(p.Pos, "<пакет>", stmt.String())
			return nil, err
		}
	}

	return nil, nil
}

// Evaluate executes block of statements.
// Returns (result value, force stop flag, error)
func (b *BlockStmts) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	for _, stmt := range b.Stmts {
		result := stmt.Evaluate(state, inFunction, inLoop)
		if result.Err != nil {
			return result
		}

		switch result.State {
		case StmtForceReturn, StmtBreak:
			return result
		}
	}

	return StmtResult{Value: types.NewNilInstance()}
}

func (e *Expression) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	if e.LogicalAnd != nil {
		return e.LogicalAnd.Evaluate(state, valueToSet)
	}

	panic("unreachable")
}

func (a *Assignment) Evaluate(state common.State) (common.Value, error) {
	if len(a.Next) == 0 {
		return a.Expressions[0].Evaluate(state, nil)
	}

	return unpack(state, a.Expressions, a.Next)
}

// Evaluate executes LogicalAnd operation.
// If `valueToSet` is nil, return variable or value from context,
// set a new value or return an error otherwise.
func (a *LogicalAnd) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.AndOp.Name(), a.LogicalOr, a.Next)
}

func (o *LogicalOr) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.OrOp.Name(), o.LogicalNot, o.Next)
}

func (a *LogicalNot) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	if a.Comparison != nil {
		return a.Comparison.Evaluate(state, valueToSet)
	}

	if a.Next != nil {
		value, err := a.Next.Evaluate(state, nil)
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

func (a *Comparison) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch a.Op {
	case ">=":
		return evalBinaryOperator(state, valueToSet, common.GreaterOrEqualsOp.Name(), a.BitwiseOr, a.Next)
	case ">":
		return evalBinaryOperator(state, valueToSet, common.GreaterOp.Name(), a.BitwiseOr, a.Next)
	case "<=":
		return evalBinaryOperator(state, valueToSet, common.LessOrEqualsOp.Name(), a.BitwiseOr, a.Next)
	case "<":
		return evalBinaryOperator(state, valueToSet, common.LessOp.Name(), a.BitwiseOr, a.Next)
	case "==":
		return evalBinaryOperator(state, valueToSet, common.EqualsOp.Name(), a.BitwiseOr, a.Next)
	case "!=":
		return evalBinaryOperator(state, valueToSet, common.NotEqualsOp.Name(), a.BitwiseOr, a.Next)
	default:
		return a.BitwiseOr.Evaluate(state, valueToSet)
	}
}

func (a *BitwiseOr) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.BitwiseOrOp.Name(), a.BitwiseXor, a.Next)
}

func (a *BitwiseXor) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.BitwiseXorOp.Name(), a.BitwiseAnd, a.Next)
}

func (a *BitwiseAnd) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.BitwiseAndOp.Name(), a.BitwiseShift, a.Next)
}

func (a *BitwiseShift) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch a.Op {
	case "<<":
		return evalBinaryOperator(state, valueToSet, common.BitwiseLeftShiftOp.Name(), a.Addition, a.Next)
	case ">>":
		return evalBinaryOperator(state, valueToSet, common.BitwiseRightShiftOp.Name(), a.Addition, a.Next)
	default:
		return a.Addition.Evaluate(state, valueToSet)
	}
}

func (a *Addition) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch a.Op {
	case "+":
		return evalBinaryOperator(state, valueToSet, common.AddOp.Name(), a.MultiplicationOrMod, a.Next)
	case "-":
		return evalBinaryOperator(state, valueToSet, common.SubOp.Name(), a.MultiplicationOrMod, a.Next)
	default:
		return a.MultiplicationOrMod.Evaluate(state, valueToSet)
	}
}

func (a *MultiplicationOrMod) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch a.Op {
	case "/":
		return evalBinaryOperator(state, valueToSet, common.DivOp.Name(), a.Unary, a.Next)
	case "*":
		return evalBinaryOperator(state, valueToSet, common.MulOp.Name(), a.Unary, a.Next)
	case "%":
		return evalBinaryOperator(state, valueToSet, common.ModuloOp.Name(), a.Unary, a.Next)
	default:
		return a.Unary.Evaluate(state, valueToSet)
	}
}

func (a *Unary) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	switch a.Op {
	case "+":
		return evalUnaryOperator(state, common.UnaryPlus.Name(), a.Next)
	case "-":
		return evalUnaryOperator(state, common.UnaryMinus.Name(), a.Next)
	case "~":
		return evalUnaryOperator(state, common.UnaryBitwiseNotOp.Name(), a.Next)
	default:
		return a.Exponent.Evaluate(state, valueToSet)
	}
}

func (a *Exponent) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	return evalBinaryOperator(state, valueToSet, common.PowOp.Name(), a.Primary, a.Next)
}

func (a *Primary) Evaluate(state common.State, valueToSet common.Value) (common.Value, error) {
	if a.SubExpression != nil {
		if valueToSet != nil {
			// TODO: change to normal description
			return nil, errors.New("unable to set to subexpression evaluation")
		}

		return a.SubExpression.Evaluate(state, valueToSet)
	}

	if a.Constant != nil {
		if valueToSet != nil {
			// TODO: change to normal description
			return nil, errors.New("unable to set to constant")
		}

		return a.Constant.Evaluate(state)
	}

	if a.AttributeAccess != nil {
		return a.AttributeAccess.Evaluate(state, valueToSet, nil)
	}

	if a.LambdaDef != nil {
		return a.LambdaDef.Evaluate(state)
	}

	panic("unreachable")
}

func (c *Constant) Evaluate(state common.State) (common.Value, error) {
	if c.Integer != nil {
		return types.NewIntegerInstance(*c.Integer), nil
	}

	if c.Real != nil {
		return types.NewRealInstance(*c.Real), nil
	}

	if c.Bool != nil {
		return types.NewBoolInstance(bool(*c.Bool)), nil
	}

	if c.StringValue != nil {
		return types.NewStringInstance(*c.StringValue), nil
	}

	if c.List != nil {
		list := types.NewListInstance()
		for _, expr := range c.List {
			value, err := expr.Evaluate(state, nil)
			if err != nil {
				return nil, err
			}

			list.Values = append(list.Values, value)
		}

		return list, nil
	}

	if c.EmptyList {
		return types.NewListInstance(), nil
	}

	if c.Dictionary != nil {
		dict := types.NewDictionaryInstance()
		for _, entry := range c.Dictionary {
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

	if c.EmptyDictionary {
		return types.NewDictionaryInstance(), nil
	}

	panic("unreachable")
}

func (d *DictionaryEntry) Evaluate(state common.State) (common.Value, common.Value, error) {
	key, err := d.Key.Evaluate(state, nil)
	if err != nil {
		return nil, nil, err
	}

	value, err := d.Value.Evaluate(state, nil)
	if err != nil {
		return nil, nil, err
	}

	return key, value, nil
}

func (a *AttributeAccess) Evaluate(state common.State, valueToSet, prevValue common.Value) (common.Value, error) {
	if a.SlicingOrSubscription == nil {
		panic("unreachable")
	}

	if valueToSet != nil {
		// set
		var currentValue common.Value
		var err error
		if a.AttributeAccess != nil {
			currentValue, err = a.SlicingOrSubscription.Evaluate(state, nil, prevValue)
			if err != nil {
				return nil, err
			}

			currentValue, err = a.AttributeAccess.Evaluate(state, valueToSet, currentValue)
		} else {
			currentValue, err = a.SlicingOrSubscription.Evaluate(state, valueToSet, prevValue)
		}

		if err != nil {
			return nil, err
		}

		return currentValue, nil
	}

	// get
	currentValue, err := a.SlicingOrSubscription.Evaluate(state, valueToSet, prevValue)
	if err != nil {
		return nil, err
	}

	if a.AttributeAccess != nil {
		return a.AttributeAccess.Evaluate(state, valueToSet, currentValue)
	}

	return currentValue, err
}

func (s *SlicingOrSubscription) Evaluate(
	state common.State,
	valueToSet common.Value,
	prevValue common.Value,
) (common.Value, error) {
	if valueToSet != nil {
		// set
		var variable common.Value
		var err error = nil
		rangesLen := len(s.Ranges)
		if rangesLen != 0 && s.Ranges[rangesLen-1].RightBound != nil {
			return nil, util.RuntimeError("неможливо присвоїти значення зрізу")
		}

		if s.Call != nil {
			if len(s.Ranges) == 0 {
				return nil, util.RuntimeError("неможливо присвоїти значення виклику функції")
			}

			variable, err = s.callFunction(state, prevValue)
			if err != nil {
				return nil, err
			}
		} else if s.Ident != nil {
			if len(s.Ranges) != 0 {
				variable, err = getCurrentValue(state.GetContext(), prevValue, *s.Ident)
			} else {
				variable, err = setCurrentValue(state.GetContext(), prevValue, *s.Ident, valueToSet)
			}

			if err != nil {
				return nil, err
			}
		} else {
			panic("unreachable")
		}

		if len(s.Ranges) != 0 {
			return evalSlicingOperation(state, variable, s.Ranges, valueToSet)
		}

		return variable, nil
	}

	// get
	var variable common.Value
	var err error = nil
	if s.Call != nil {
		variable, err = s.callFunction(state, prevValue)
		if err != nil {
			return nil, err
		}
	} else if s.Ident != nil {
		variable, err = getCurrentValue(state.GetContext(), prevValue, *s.Ident)
		if err != nil {
			return nil, err
		}
	} else {
		panic("unreachable")
	}

	if len(s.Ranges) != 0 {
		return evalSlicingOperation(state, variable, s.Ranges, nil)
	}

	return variable, nil
}

func (s *SlicingOrSubscription) callFunction(state common.State, prevValue common.Value) (common.Value, error) {
	ctx := state.GetContext()
	variable, err := getCurrentValue(ctx, prevValue, s.Call.Ident)
	if err != nil {
		return nil, err
	}

	isLambda := false
	variable, err = s.Call.Evaluate(state, variable, prevValue, &isLambda)
	if err != nil {
		funcName := s.Call.Ident
		if isLambda {
			funcName = common.LambdaSignature
		}

		state.GetInterpreter().Trace(s.Call.Pos, funcName, s.Call.String())
		return nil, err
	}

	return variable, nil
}

func (l *LambdaDef) Evaluate(state common.State) (common.Value, error) {
	arguments, err := l.ParametersSet.Evaluate(state)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(state, l.ReturnTypes)
	if err != nil {
		return nil, err
	}

	lambda := types.NewFunctionInstance(
		common.LambdaSignature,
		arguments,
		func(state common.State, _ *[]common.Value, kwargs *map[string]common.Value) (common.Value, error) {
			return l.Body.Evaluate(state)
		},
		returnTypes,
		false,
		state.GetCurrentPackage().(*types.PackageInstance),
		"", // TODO: add doc
	)

	if l.InstantCall {
		return l.evalInstantCall(state, lambda)
	}

	return lambda, nil
}

func (l *LambdaDef) evalInstantCall(state common.State, function *types.FunctionInstance) (common.Value, error) {
	var args []common.Value
	if len(l.InstantCallArguments) != 0 {
		if err := updateArgs(state, l.InstantCallArguments, &args); err != nil {
			return nil, err
		}
	}

	return types.Call(state, function, &args, nil)
}
