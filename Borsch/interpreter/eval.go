package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Scope map[string]common.Type

func (p *Package) Evaluate(state common.State) (common.Type, error) {
	state.GetContext().PushScope(Scope{})
	for _, stmt := range p.Stmts {
		state := stmt.Evaluate(state, false, false)
		if state.Err != nil {
			pos := stmt.Pos
			return nil, errors.New(
				fmt.Sprintf(
					"  Файл \"%s\", рядок %d, позиція %d,\n    %s\n%s",
					pos.Filename, pos.Line, pos.Column, stmt.String(), state.Err.Error(),
				),
			)
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
		// if result.State == StmtForceReturn || stmt.ReturnStmt != nil {
		// 	return result
		// }
	}

	return StmtResult{Value: types.NewNilInstance()}
}

func (e *Expression) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	if e.LogicalAnd != nil {
		return e.LogicalAnd.Evaluate(state, valueToSet)
	}

	panic("unreachable")
}

func (a *Assignment) Evaluate(state common.State) (common.Type, error) {
	if len(a.Next) == 0 {
		return a.Expression[0].Evaluate(state, nil)
	}

	return unpack(state, a.Expression, a.Next)
}

// Evaluate executes LogicalAnd operation.
// If `valueToSet` is nil, return variable or value from context,
// set a new value or return an error otherwise.
func (a *LogicalAnd) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(state, valueToSet, ops.AndOp.Name(), a.LogicalOr, a.Next)
}

func (a *LogicalOr) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(state, valueToSet, ops.OrOp.Name(), a.LogicalNot, a.Next)
}

func (a *LogicalNot) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	if a.Comparison != nil {
		return a.Comparison.Evaluate(state, valueToSet)
	}

	if a.Next != nil {
		value, err := a.Next.Evaluate(state, nil)
		if err != nil {
			return nil, err
		}

		return builtin.CallByName(state, value, ops.NotOp.Name(), &[]common.Type{}, nil, true)
	}

	panic("unreachable")
}

func (a *Comparison) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case ">=":
		return evalBinaryOperator(state, valueToSet, ops.GreaterOrEqualsOp.Name(), a.BitwiseOr, a.Next)
	case ">":
		return evalBinaryOperator(state, valueToSet, ops.GreaterOp.Name(), a.BitwiseOr, a.Next)
	case "<=":
		return evalBinaryOperator(state, valueToSet, ops.LessOrEqualsOp.Name(), a.BitwiseOr, a.Next)
	case "<":
		return evalBinaryOperator(state, valueToSet, ops.LessOp.Name(), a.BitwiseOr, a.Next)
	case "==":
		return evalBinaryOperator(state, valueToSet, ops.EqualsOp.Name(), a.BitwiseOr, a.Next)
	case "!=":
		return evalBinaryOperator(state, valueToSet, ops.NotEqualsOp.Name(), a.BitwiseOr, a.Next)
	default:
		return a.BitwiseOr.Evaluate(state, valueToSet)
	}
}

func (a *BitwiseOr) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(state, valueToSet, ops.BitwiseOrOp.Name(), a.BitwiseXor, a.Next)
}

func (a *BitwiseXor) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(state, valueToSet, ops.BitwiseXorOp.Name(), a.BitwiseAnd, a.Next)
}

func (a *BitwiseAnd) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(state, valueToSet, ops.BitwiseAndOp.Name(), a.BitwiseShift, a.Next)
}

func (a *BitwiseShift) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case "<<":
		return evalBinaryOperator(state, valueToSet, ops.BitwiseLeftShiftOp.Name(), a.Addition, a.Next)
	case ">>":
		return evalBinaryOperator(state, valueToSet, ops.BitwiseRightShiftOp.Name(), a.Addition, a.Next)
	default:
		return a.Addition.Evaluate(state, valueToSet)
	}
}

func (a *Addition) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case "+":
		return evalBinaryOperator(state, valueToSet, ops.AddOp.Name(), a.MultiplicationOrMod, a.Next)
	case "-":
		return evalBinaryOperator(state, valueToSet, ops.SubOp.Name(), a.MultiplicationOrMod, a.Next)
	default:
		return a.MultiplicationOrMod.Evaluate(state, valueToSet)
	}
}

func (a *MultiplicationOrMod) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case "/":
		return evalBinaryOperator(state, valueToSet, ops.DivOp.Name(), a.Unary, a.Next)
	case "*":
		return evalBinaryOperator(state, valueToSet, ops.MulOp.Name(), a.Unary, a.Next)
	case "%":
		return evalBinaryOperator(state, valueToSet, ops.ModuloOp.Name(), a.Unary, a.Next)
	default:
		return a.Unary.Evaluate(state, valueToSet)
	}
}

func (a *Unary) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case "+":
		return evalUnaryOperator(state, ops.UnaryPlus.Name(), a.Next)
	case "-":
		return evalUnaryOperator(state, ops.UnaryMinus.Name(), a.Next)
	case "~":
		return evalUnaryOperator(state, ops.UnaryBitwiseNotOp.Name(), a.Next)
	default:
		return a.Exponent.Evaluate(state, valueToSet)
	}
}

func (a *Exponent) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(state, valueToSet, ops.PowOp.Name(), a.Primary, a.Next)
}

func (a *Primary) Evaluate(state common.State, valueToSet common.Type) (common.Type, error) {
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

func (c *Constant) Evaluate(state common.State) (common.Type, error) {
	if c.Integer != nil {
		return types.NewIntegerInstance(*c.Integer), nil
	}

	if c.Real != nil {
		return types.NewRealInstance(*c.Real), nil
	}

	if c.Bool != nil {
		return types.NewBoolInstance(bool(*c.Bool)), nil
	}

	if c.String != nil {
		return types.NewStringInstance(*c.String), nil
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

func (d *DictionaryEntry) Evaluate(state common.State) (common.Type, common.Type, error) {
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

func (a *AttributeAccess) Evaluate(state common.State, valueToSet, prevValue common.Type) (common.Type, error) {
	if a.Slicing == nil {
		panic("unreachable")
	}

	if valueToSet != nil {
		// set
		var currentValue common.Type
		var err error
		if a.AttributeAccess != nil {
			currentValue, err = a.Slicing.Evaluate(state, nil, prevValue)
			if err != nil {
				return nil, err
			}

			currentValue, err = a.AttributeAccess.Evaluate(state, valueToSet, currentValue)
		} else {
			currentValue, err = a.Slicing.Evaluate(state, valueToSet, prevValue)
		}

		if err != nil {
			return nil, err
		}

		return currentValue, nil
	}

	// get
	currentValue, err := a.Slicing.Evaluate(state, valueToSet, prevValue)
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
	valueToSet common.Type,
	prevValue common.Type,
) (common.Type, error) {
	if valueToSet != nil {
		// set
		var variable common.Type
		var err error = nil
		rangesLen := len(s.Ranges)
		if rangesLen != 0 && s.Ranges[rangesLen-1].RightBound != nil {
			return nil, util.RuntimeError("неможливо присвоїти значення у зріз")
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
	var variable common.Type
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

func (s *SlicingOrSubscription) callFunction(state common.State, prevValue common.Type) (common.Type, error) {
	ctx := state.GetContext()
	variable, err := getCurrentValue(ctx, prevValue, s.Call.Ident)
	if err != nil {
		return nil, err
	}

	variable, err = s.Call.Evaluate(state, variable, prevValue)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"  Файл \"%s\", рядок %d, позиція %d\n    %s\n%s",
				s.Call.Pos.Filename, s.Call.Pos.Line, s.Call.Pos.Column, "TODO", err.Error(),
			),
		)
	}

	return variable, nil
}

func (l *LambdaDef) Evaluate(state common.State) (common.Type, error) {
	arguments, err := l.ParametersSet.Evaluate(state)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(state, l.ReturnTypes)
	if err != nil {
		return nil, err
	}

	return types.NewFunctionInstance(
		"",
		arguments,
		func(state common.State, _ *[]common.Type, kwargs *map[string]common.Type) (common.Type, error) {
			return l.Body.Evaluate(state)
		},
		returnTypes,
		false,
		state.GetCurrentPackage().(*types.PackageInstance),
		"", // TODO: add doc
	), nil
}
