package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Scope map[string]common.Type

func (p *Package) Evaluate(ctx common.Context) (common.Type, error) {
	ctx.PushScope(Scope{})
	for _, stmt := range p.Stmts {
		state := stmt.Evaluate(ctx, false, false)
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

	if err := ctx.BuildPackage(); err != nil {
		return nil, err
	}

	ctx.PopScope()
	return ctx.GetPackage(), nil
}

// Evaluate executes block of statements.
// Returns (result value, force stop flag, error)
func (b *BlockStmts) Evaluate(ctx common.Context, inFunction, inLoop bool) StmtResult {
	for _, stmt := range b.Stmts {
		result := stmt.Evaluate(ctx, inFunction, inLoop)
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

func (e *Expression) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	if e.LogicalAnd != nil {
		return e.LogicalAnd.Evaluate(ctx, valueToSet)
	}

	panic("unreachable")
}

func (a *Assignment) Evaluate(ctx common.Context) (common.Type, error) {
	if len(a.Next) == 0 {
		return a.Expression[0].Evaluate(ctx, nil)
	}

	return unpack(ctx, a.Expression, a.Next)
}

// Evaluate executes LogicalAnd operation.
// If `valueToSet` is nil, return variable or value from context,
// set a new value or return an error otherwise.
func (a *LogicalAnd) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.AndOp.Caption(), a.LogicalOr, a.Next)
}

func (a *LogicalOr) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.OrOp.Caption(), a.LogicalNot, a.Next)
}

func (a *LogicalNot) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	if a.Comparison != nil {
		return a.Comparison.Evaluate(ctx, valueToSet)
	}

	if a.Next != nil {
		value, err := a.Next.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		return callMethod(ctx, value, ops.NotOp.Caption(), &[]common.Type{}, nil)
	}

	panic("unreachable")
}

func (a *Comparison) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case ">=":
		return evalBinaryOperator(ctx, valueToSet, ops.GreaterOrEqualsOp.Caption(), a.BitwiseOr, a.Next)
	case ">":
		return evalBinaryOperator(ctx, valueToSet, ops.GreaterOp.Caption(), a.BitwiseOr, a.Next)
	case "<=":
		return evalBinaryOperator(ctx, valueToSet, ops.LessOrEqualsOp.Caption(), a.BitwiseOr, a.Next)
	case "<":
		return evalBinaryOperator(ctx, valueToSet, ops.LessOp.Caption(), a.BitwiseOr, a.Next)
	case "==":
		return evalBinaryOperator(ctx, valueToSet, ops.EqualsOp.Caption(), a.BitwiseOr, a.Next)
	case "!=":
		return evalBinaryOperator(ctx, valueToSet, ops.NotEqualsOp.Caption(), a.BitwiseOr, a.Next)
	default:
		return a.BitwiseOr.Evaluate(ctx, valueToSet)
	}
}

func (a *BitwiseOr) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.BitwiseOrOp.Caption(), a.BitwiseXor, a.Next)
}

func (a *BitwiseXor) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.BitwiseXorOp.Caption(), a.BitwiseAnd, a.Next)
}

func (a *BitwiseAnd) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.BitwiseAndOp.Caption(), a.BitwiseShift, a.Next)
}

func (a *BitwiseShift) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case "<<":
		return evalBinaryOperator(ctx, valueToSet, ops.BitwiseLeftShiftOp.Caption(), a.Addition, a.Next)
	case ">>":
		return evalBinaryOperator(ctx, valueToSet, ops.BitwiseRightShiftOp.Caption(), a.Addition, a.Next)
	default:
		return a.Addition.Evaluate(ctx, valueToSet)
	}
}

func (a *Addition) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case "+":
		return evalBinaryOperator(ctx, valueToSet, ops.AddOp.Caption(), a.MultiplicationOrMod, a.Next)
	case "-":
		return evalBinaryOperator(ctx, valueToSet, ops.SubOp.Caption(), a.MultiplicationOrMod, a.Next)
	default:
		return a.MultiplicationOrMod.Evaluate(ctx, valueToSet)
	}
}

func (a *MultiplicationOrMod) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case "/":
		return evalBinaryOperator(ctx, valueToSet, ops.DivOp.Caption(), a.Unary, a.Next)
	case "*":
		return evalBinaryOperator(ctx, valueToSet, ops.MulOp.Caption(), a.Unary, a.Next)
	case "%":
		return evalBinaryOperator(ctx, valueToSet, ops.ModuloOp.Caption(), a.Unary, a.Next)
	default:
		return a.Unary.Evaluate(ctx, valueToSet)
	}
}

func (a *Unary) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	switch a.Op {
	case "+":
		return evalUnaryOperator(ctx, ops.UnaryPlus.Caption(), a.Next)
	case "-":
		return evalUnaryOperator(ctx, ops.UnaryMinus.Caption(), a.Next)
	case "~":
		return evalUnaryOperator(ctx, ops.UnaryBitwiseNotOp.Caption(), a.Next)
	default:
		return a.Exponent.Evaluate(ctx, valueToSet)
	}
}

func (a *Exponent) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.PowOp.Caption(), a.Primary, a.Next)
}

func (a *Primary) Evaluate(ctx common.Context, valueToSet common.Type) (common.Type, error) {
	if a.SubExpression != nil {
		if valueToSet != nil {
			// TODO: change to normal description
			return nil, errors.New("unable to set to subexpression evaluation")
		}

		return a.SubExpression.Evaluate(ctx, valueToSet)
	}

	if a.Constant != nil {
		if valueToSet != nil {
			// TODO: change to normal description
			return nil, errors.New("unable to set to constant")
		}

		return a.Constant.Evaluate(ctx)
	}

	if a.AttributeAccess != nil {
		return a.AttributeAccess.Evaluate(ctx, valueToSet, nil)
	}

	if a.LambdaDef != nil {
		return a.LambdaDef.Evaluate(ctx)
	}

	panic("unreachable")
}

func (c *Constant) Evaluate(ctx common.Context) (common.Type, error) {
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
			value, err := expr.Evaluate(ctx, nil)
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
			key, value, err := entry.Evaluate(ctx)
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

func (d *DictionaryEntry) Evaluate(ctx common.Context) (common.Type, common.Type, error) {
	key, err := d.Key.Evaluate(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	value, err := d.Value.Evaluate(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	return key, value, nil
}

func (a *AttributeAccess) Evaluate(ctx common.Context, valueToSet, prevValue common.Type) (common.Type, error) {
	if a.Slicing == nil {
		panic("unreachable")
	}

	if valueToSet != nil {
		// set
		var currentValue common.Type
		var err error
		if a.AttributeAccess != nil {
			currentValue, err = a.Slicing.Evaluate(ctx, nil, prevValue)
			if err != nil {
				return nil, err
			}

			currentValue, err = a.AttributeAccess.Evaluate(ctx, valueToSet, currentValue)
		} else {
			currentValue, err = a.Slicing.Evaluate(ctx, valueToSet, prevValue)
		}

		if err != nil {
			return nil, err
		}

		return currentValue, nil
	}

	// get
	currentValue, err := a.Slicing.Evaluate(ctx, valueToSet, prevValue)
	if err != nil {
		return nil, err
	}

	if a.AttributeAccess != nil {
		return a.AttributeAccess.Evaluate(ctx, valueToSet, currentValue)
	}

	return currentValue, err
}

func (s *SlicingOrSubscription) Evaluate(
	ctx common.Context,
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

			variable, err = s.callFunction(ctx, prevValue)
			if err != nil {
				return nil, err
			}
		} else if s.Ident != nil {
			variable, err = setCurrentValue(ctx, prevValue, *s.Ident, valueToSet)
			if err != nil {
				return nil, err
			}
		} else {
			panic("unreachable")
		}

		if len(s.Ranges) != 0 {
			return evalSlicingOperation(ctx, variable, s.Ranges, valueToSet)
		}

		return variable, nil
	}

	// get
	var variable common.Type
	var err error = nil
	if s.Call != nil {
		variable, err = s.callFunction(ctx, prevValue)
		if err != nil {
			return nil, err
		}
	} else if s.Ident != nil {
		variable, err = getCurrentValue(ctx, prevValue, *s.Ident)
		if err != nil {
			return nil, err
		}
	} else {
		panic("unreachable")
	}

	if len(s.Ranges) != 0 {
		return evalSlicingOperation(ctx, variable, s.Ranges, nil)
	}

	return variable, nil
}

func (s *SlicingOrSubscription) callFunction(ctx common.Context, prevValue common.Type) (common.Type, error) {
	variable, err := getCurrentValue(ctx, prevValue, s.Call.Ident)
	if err != nil {
		return nil, err
	}

	variable, err = s.Call.Evaluate(ctx, variable, prevValue)
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

func (l *LambdaDef) Evaluate(ctx common.Context) (common.Type, error) {
	arguments, err := l.ParametersSet.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(ctx, l.ReturnTypes)
	if err != nil {
		return nil, err
	}

	return types.NewFunctionInstance(
		"",
		arguments,
		func(ctx common.Context, _ *[]common.Type, kwargs *map[string]common.Type) (common.Type, error) {
			return l.Body.Evaluate(ctx)
		},
		returnTypes,
		false,
		ctx.GetPackage().(*types.PackageInstance),
		"", // TODO: add doc
	), nil
}
