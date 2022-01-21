package grammar

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
	"github.com/alecthomas/participle/v2/lexer"
)

type Scope map[string]common.Type

func (p *Package) Evaluate(ctx common.Context) (common.Type, error) {
	ctx.PushScope(Scope{})
	for _, stmt := range p.Stmts {
		_, _, err := stmt.Evaluate(ctx, false)
		if err != nil {
			pos := stmt.getPos()
			return nil, errors.New(
				fmt.Sprintf(
					"  Файл \"%s\", рядок %d, позиція %d,\n    %s\n%s",
					pos.Filename, pos.Line, pos.Column, stmt.String(), err.Error(),
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

func (s *WhileStmt) Evaluate(ctx common.Context) (common.Type, bool, error) {
	// TODO:
	panic("unreachable")
}

func (s *IfStmt) Evaluate(ctx common.Context, inFunction bool) (
	common.Type,
	bool,
	error,
) {
	if s.Condition != nil {
		condition, err := s.Condition.Evaluate(ctx, nil)
		if err != nil {
			return nil, false, err
		}

		if condition.AsBool(ctx) {
			ctx.PushScope(Scope{})
			result, forceReturn, err := s.Body.Evaluate(ctx, inFunction)
			if err != nil {
				return nil, false, err
			}

			ctx.PopScope()
			return result, forceReturn, nil
		}

		if len(s.ElseIfStmts) != 0 {
			gotResult := false
			var result common.Type = nil
			var err error = nil
			for _, stmt := range s.ElseIfStmts {
				ctx.PushScope(Scope{})
				var forceReturn bool
				gotResult, result, forceReturn, err = stmt.Evaluate(ctx, inFunction)
				if err != nil {
					return nil, false, err
				}

				ctx.PopScope()
				if forceReturn {
					return result, true, nil
				}

				if gotResult {
					break
				}
			}

			if gotResult {
				return result, false, nil
			}
		}

		if s.Else != nil {
			ctx.PushScope(Scope{})
			result, forceReturn, err := s.Else.Evaluate(ctx, inFunction)
			if err != nil {
				return nil, false, err
			}

			ctx.PopScope()
			return result, forceReturn, nil
		}

		return nil, false, nil
	}

	return nil, false, errors.New("interpreter: condition is nil")
}

func (s *ElseIfStmt) Evaluate(ctx common.Context, inFunction bool) (
	bool,
	common.Type,
	bool,
	error,
) {
	condition, err := s.Condition.Evaluate(ctx, nil)
	if err != nil {
		return false, nil, false, err
	}

	if condition.AsBool(ctx) {
		ctx.PushScope(Scope{})
		result, forceReturn, err := s.Body.Evaluate(ctx, inFunction)
		if err != nil {
			return false, nil, false, err
		}

		ctx.PopScope()
		return true, result, forceReturn, nil
	}

	return false, nil, false, nil
}

// Evaluate executes block of statements.
// Returns (result value, force stop flag, error)
func (b *BlockStmts) Evaluate(ctx common.Context, inFunction bool) (common.Type, bool, error) {
	for _, stmt := range b.Stmts {
		result, forceReturn, err := stmt.Evaluate(ctx, inFunction)
		if err != nil {
			return nil, false, err
		}

		if forceReturn || stmt.ReturnStmt != nil {
			return result, true, nil
		}
	}

	return types.NewNilInstance(), false, nil
}

// Evaluate executes statement.
// Returns (result value, force stop flag, error)
func (s *Stmt) Evaluate(ctx common.Context, inFunction bool) (
	common.Type,
	bool,
	error,
) {
	if s.IfStmt != nil {
		return s.IfStmt.Evaluate(ctx, inFunction)
	} else if s.WhileStmt != nil {
		return s.WhileStmt.Evaluate(ctx)
	} else if s.Block != nil {
		ctx.PushScope(Scope{})
		result, forceReturn, err := s.Block.Evaluate(ctx, inFunction)
		if err != nil {
			return nil, false, err
		}

		ctx.PopScope()
		return result, forceReturn, nil
	} else if s.FunctionDef != nil {
		function, err := s.FunctionDef.Evaluate(ctx, ctx.GetPackage().(*types.PackageInstance), nil)
		if err != nil {
			return nil, false, err
		}

		return function, false, err
	} else if s.ClassDef != nil {
		class, err := s.ClassDef.Evaluate(ctx, ctx.GetPackage().(*types.PackageInstance))
		if err != nil {
			return nil, false, err
		}

		return class, false, err
	} else if s.ReturnStmt != nil {
		if !inFunction {
			return nil, false, errors.New("'повернути' за межами функції")
		}

		result, err := s.ReturnStmt.Evaluate(ctx)
		return result, false, err
		// } else if s.Expression != nil {
		//	result, err := s.Expression.Evaluate(ctx, nil)
		//	return result, false, err
	} else if s.Assignment != nil {
		result, err := s.Assignment.Evaluate(ctx)
		return result, false, err
	} else if s.Empty {
		return nil, false, nil
	}

	panic("unreachable")
}

func (s *Stmt) getPos() lexer.Position {
	if s.IfStmt != nil {
		return s.IfStmt.Pos
	} else if s.WhileStmt != nil {
		return s.WhileStmt.Pos
	} else if s.Block != nil {
		return s.Block.Pos
	} else if s.FunctionDef != nil {
		return s.FunctionDef.Pos
	} else if s.ClassDef != nil {
		return s.ClassDef.Pos
	} else if s.ReturnStmt != nil {
		return s.ReturnStmt.Pos
	} else if s.Assignment != nil {
		return s.Assignment.Pos
	} else if s.Empty {
		return s.Pos
	}

	panic("unreachable")
}

func (s *Stmt) String() string {
	if s.IfStmt != nil {
		return "s.IfStmt."
	} else if s.WhileStmt != nil {
		return "s.WhileStmt."
	} else if s.Block != nil {
		return "s.Block."
	} else if s.FunctionDef != nil {
		return "s.FunctionDef."
	} else if s.ClassDef != nil {
		return "s.ClassDef."
	} else if s.ReturnStmt != nil {
		return "повернути ..."
	} else if s.Assignment != nil {
		return "s.Assignment."
	} else if s.Empty {
		return ";"
	}

	panic("unreachable")
}

func (b *FunctionBody) Evaluate(ctx common.Context) (common.Type, error) {
	result, _, err := b.Stmts.Evaluate(ctx, true)
	return result, err
}

func (f *FunctionDef) Evaluate(
	ctx common.Context,
	parentPackage *types.PackageInstance,
	check func([]types.FunctionArgument, []types.FunctionReturnType) error,
) (common.Type, error) {
	arguments, err := evalParameters(ctx, f.Parameters)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(ctx, f.ReturnTypes)
	if err != nil {
		return nil, err
	}

	if check != nil {
		if err := check(arguments, returnTypes); err != nil {
			return nil, err
		}
	}

	function := types.NewFunctionInstance(
		f.Name,
		arguments,
		func(ctx common.Context, _ *[]common.Type, kwargs *map[string]common.Type) (common.Type, error) {
			return f.Body.Evaluate(ctx)
		},
		returnTypes,
		parentPackage == nil,
		parentPackage,
		"", // TODO: add doc
	)
	return function, ctx.SetVar(f.Name, function)
}

func (p *Parameter) Evaluate(ctx common.Context) (*types.FunctionArgument, error) {
	class, err := ctx.GetClass(p.Type)
	if err != nil {
		return nil, err
	}

	return &types.FunctionArgument{
		Type:       class.(*types.Class),
		Name:       p.Name,
		IsVariadic: false,
		IsNullable: p.IsNullable,
	}, nil
}

func (t *ReturnType) Evaluate(ctx common.Context) (*types.FunctionReturnType, error) {
	class, err := ctx.GetClass(t.Name)
	if err != nil {
		return nil, err
	}

	return &types.FunctionReturnType{
		Type:       class.(*types.Class),
		IsNullable: t.IsNullable,
	}, nil
}

func (s *ReturnStmt) Evaluate(ctx common.Context) (common.Type, error) {
	resultCount := len(s.Expressions)
	switch {
	case resultCount == 1:
		return s.Expressions[0].Evaluate(ctx, nil)
	case resultCount > 1:
		result := types.NewListInstance()
		for _, expression := range s.Expressions {
			value, err := expression.Evaluate(ctx, nil)
			if err != nil {
				return nil, err
			}

			result.Values = append(result.Values, value)
		}

		return result, nil
	}

	panic("unreachable")
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

func (a *Unary) Evaluate(ctx common.Context, valueToSet common.Type) (
	common.Type,
	error,
) {
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

// func (a *Subscription) Evaluate(ctx common.Context, valueToSet common.Type, prevValue common.Type) (
// 	common.Type,
// 	error,
// ) {
// 	if a.Slicing == nil {
// 		panic("unreachable")
// 	}
//
// 	if valueToSet != nil {
// 		// set
// 		var variable common.Type
// 		var err error
// 		if len(a.Indices) != 0 {
// 			variable, err = a.Slicing.Evaluate(ctx, nil, prevValue)
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			variable, err = evalSingleSetByIndexOperation(ctx, variable, a.Indices, valueToSet)
// 		} else {
// 			variable, err = a.Slicing.Evaluate(ctx, valueToSet, prevValue)
// 		}
//
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		return variable, nil
// 	}
//
// 	// get
// 	variable, err := a.Slicing.Evaluate(ctx, valueToSet, prevValue)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for _, expression := range a.Indices {
// 		index, err := expression.Evaluate(ctx, nil)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		variable, err = evalSingleGetByIndexOperation(ctx, variable, index)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	return variable, nil
// }

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

		if s.CallFunc != nil {
			return nil, util.RuntimeError("неможливо присвоїти значення виклику функції")
		} else if s.Ident != nil {
			if *s.Ident == "нуль" {
				return nil, util.RuntimeError("неможливо встановити значення об'єкту 'нуль'")
			}

			if prevValue != nil {
				variable, err = prevValue.SetAttribute(*s.Ident, valueToSet)
			} else {
				variable = valueToSet
				err = ctx.SetVar(*s.Ident, valueToSet)
			}

			if err != nil {
				return nil, err
			}

		} else {
			panic("unreachable")
		}

		variable, err = evalSingleSetByIndexOperation(ctx, variable, s.Ranges, valueToSet)
		if err != nil {
			return nil, err
		}
	}

	// get
	var variable common.Type
	var err error = nil
	if s.CallFunc != nil {
		variable, err = getCurrentValue(ctx, prevValue, s.CallFunc.Ident)
		if err != nil {
			return nil, err
		}

		variable, err = s.CallFunc.Evaluate(ctx, variable, prevValue)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf(
					"  Файл \"%s\", рядок %d, позиція %d\n    %s\n%s",
					s.CallFunc.Pos.Filename, s.CallFunc.Pos.Line, s.CallFunc.Pos.Column, "TODO", err.Error(),
				),
			)
		}
	} else if s.Ident != nil {
		if *s.Ident == "нуль" {
			if prevValue == nil {
				return types.NewNilInstance(), nil
			} else {
				return nil, util.RuntimeError("'нуль' не є атрибутом")
			}
		}

		variable, err = getCurrentValue(ctx, prevValue, *s.Ident)
		if err != nil {
			return nil, err
		}
	} else {
		panic("unreachable")
	}

	for _, range_ := range s.Ranges {
		leftBound, err := range_.LeftBound.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		if range_.RightBound != nil {
			rightBound, err := range_.RightBound.Evaluate(ctx, nil)
			if err != nil {
				return nil, err
			}

			variable, err = evalSlicingOperation(ctx, variable, leftBound, rightBound)
			if err != nil {
				return nil, err
			}
		} else {
			variable, err = evalSingleGetByIndexOperation(ctx, variable, leftBound)
			if err != nil {
				return nil, err
			}
		}
	}

	return variable, nil
}

func (l *LambdaDef) Evaluate(ctx common.Context) (common.Type, error) {
	arguments, err := evalParameters(ctx, l.Parameters)
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

func (a *CallFunc) Evaluate(ctx common.Context, variable common.Type, selfInstance common.Type) (
	common.Type,
	error,
) {
	switch function := variable.(type) {
	case *types.Class:
		callable, err := function.GetAttribute(ops.ConstructorName)
		if err != nil {
			return nil, err
		}

		switch __constructor__ := callable.(type) {
		case *types.FunctionInstance:
			instance, err := function.GetEmptyInstance()
			if err != nil {
				return nil, err
			}

			args := []common.Type{instance}
			kwargs := map[string]common.Type{__constructor__.Arguments[0].Name: instance}

			// TODO: check if constructor returns nothing.
			_, err = a.evalFunction(ctx, __constructor__, &args, &kwargs, 1)
			if err != nil {
				return nil, err
			}

			return args[0], nil
		default:
			return nil, util.ObjectIsNotCallable(a.Ident, callable.GetTypeName())
		}
	case *types.FunctionInstance:
		var args []common.Type
		kwargs := map[string]common.Type{}
		argsShift := 0
		if selfInstance != nil {
			switch selfInstance.(type) {
			case *types.Class, *types.PackageInstance:
				// ignore
			case types.ObjectInstance:
				argsShift++
				args = append(args, selfInstance)
				kwargs[function.Arguments[0].Name] = selfInstance
			}
		}

		return a.evalFunction(ctx, function, &args, &kwargs, argsShift)
	case types.ObjectInstance:
		operator, err := function.GetPrototype().GetAttribute(ops.CallOperatorName)
		if err != nil {
			return nil, err
		}

		switch __call__ := operator.(type) {
		case *types.FunctionInstance:
			args := []common.Type{variable}
			kwargs := map[string]common.Type{__call__.Arguments[0].Name: variable}
			return a.evalFunction(ctx, __call__, &args, &kwargs, 1)
		default:
			return nil, util.ObjectIsNotCallable(a.Ident, operator.GetTypeName())
		}
	default:
		return nil, util.ObjectIsNotCallable(a.Ident, function.GetTypeName())
	}
}

func (a *CallFunc) evalFunction(
	ctx common.Context,
	function *types.FunctionInstance,
	args *[]common.Type,
	kwargs *map[string]common.Type,
	argsShift int,
) (common.Type, error) {
	variadicArgs := types.NewListInstance()
	variadicArgsIndex := -1
	for i, expressionArgument := range a.Arguments {
		arg, err := expressionArgument.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		*args = append(*args, arg)
		if variadicArgsIndex == -1 {
			if i+argsShift >= len(function.Arguments) {
				// TODO: return ukr error!
				return nil, util.RuntimeError("too many arguments")
			}

			if function.Arguments[i+argsShift].IsVariadic {
				variadicArgsIndex = i + argsShift
				variadicArgs.Values = append(variadicArgs.Values, arg)
			} else {
				(*kwargs)[function.Arguments[i+argsShift].Name] = arg
			}
		} else {
			variadicArgs.Values = append(variadicArgs.Values, arg)
		}
	}

	if variadicArgsIndex != -1 {
		(*kwargs)[function.Arguments[variadicArgsIndex].Name] = variadicArgs
	}

	if err := types.CheckFunctionArguments(ctx, function, args, kwargs); err != nil {
		return nil, err
	}

	ctx.PushScope(*kwargs)
	res, err := function.Call(ctx, args, kwargs)
	if err != nil {
		return nil, err
	}

	if err := types.CheckResult(ctx, res, function); err != nil {
		return nil, err
	}

	ctx.PopScope()
	return res, nil
}

func getCurrentValue(ctx common.Context, prevValue common.Type, identifier string) (common.Type, error) {
	if prevValue != nil {
		return prevValue.GetAttribute(identifier)
	}

	return ctx.GetVar(identifier)
}
