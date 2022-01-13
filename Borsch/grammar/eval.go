package grammar

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type Scope map[string]types.Type

type OperatorEvaluatable interface {
	Evaluate(*Context, types.Type) (types.Type, error)
}

func (p *Package) Evaluate(ctx *Context) (*types.PackageInstance, error) {
	ctx.PushScope(Scope{})
	for _, stmt := range p.Stmts {
		_, err := stmt.Evaluate(ctx)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf(
					"  Файл \"%s\", рядок %d, позиція %d,\n    %s\n%s",
					p.Pos.Filename, stmt.Pos.Line, stmt.Pos.Column, "TODO", err.Error(),
				),
			)
		}
	}

	ctx.package_.Attributes = ctx.scopes[len(ctx.scopes)-1]
	ctx.PopScope()
	return ctx.package_, nil
}

func (s *WhileStmt) Evaluate(ctx *Context) (types.Type, error) {
	// TODO:
	panic("unreachable")
}

func (s *IfStmt) Evaluate(ctx *Context) (types.Type, error) {
	if s.Condition != nil {
		condition, err := s.Condition.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		var args []types.Type
		kwargs := map[string]types.Type{}
		conditionResult, err := callMethod(condition, ops.BoolOperatorName, &args, &kwargs)
		if err != nil {
			return nil, err
		}

		val, err := mustBool(conditionResult)
		if err != nil {
			return nil, err
		}

		if val.Value {
			ctx.PushScope(Scope{})
			result, err := s.Body.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

			ctx.PopScope()
			return result, nil
		}

		if len(s.ElseIfStmts) != 0 {
			gotResult := false
			var result types.Type = nil
			var err error = nil
			for _, stmt := range s.ElseIfStmts {
				ctx.PushScope(Scope{})
				gotResult, result, err = stmt.Evaluate(ctx)
				if err != nil {
					return nil, err
				}

				ctx.PopScope()
				if gotResult {
					break
				}
			}

			if gotResult {
				return result, nil
			}
		}

		if s.Else != nil {
			ctx.PushScope(Scope{})
			result, err := s.Else.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

			ctx.PopScope()
			return result, nil
		}

		return nil, nil
	}

	return nil, errors.New("interpreter: condition is nil")
}

func (s *ElseIfStmt) Evaluate(ctx *Context) (bool, types.Type, error) {
	condition, err := s.Condition.Evaluate(ctx)
	if err != nil {
		return false, nil, err
	}

	boolCondition, err := mustBool(condition)
	if err != nil {
		return false, nil, err
	}

	if boolCondition.Value {
		ctx.PushScope(Scope{})
		result, err := s.Body.Evaluate(ctx)
		if err != nil {
			return false, nil, err
		}

		ctx.PopScope()
		return true, result, nil
	}

	return false, nil, nil
}

func (b *BlockStmts) Evaluate(ctx *Context) (types.Type, error) {
	for _, stmt := range b.Stmts {
		result, err := stmt.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		if stmt.ReturnStmt != nil {
			return result, nil
		}
	}

	return nil, nil
}

func (s *Stmt) Evaluate(ctx *Context) (types.Type, error) {
	if s.IfStmt != nil {
		return s.IfStmt.Evaluate(ctx)
	} else if s.WhileStmt != nil {
		return s.WhileStmt.Evaluate(ctx)
	} else if s.Block != nil {
		ctx.PushScope(Scope{})
		result, err := s.Block.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		ctx.PopScope()
		return result, nil
	} else if s.FunctionDef != nil {
		function, err := s.FunctionDef.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		return function, ctx.setVar(s.FunctionDef.Name, function)
	} else if s.ReturnStmt != nil {
		resultCount := len(s.ReturnStmt.Expressions)
		switch {
		case resultCount == 1:
			return s.ReturnStmt.Expressions[0].Evaluate(ctx)
		case resultCount > 1:
			result := types.NewListInstance()
			for _, expression := range s.ReturnStmt.Expressions {
				value, err := expression.Evaluate(ctx)
				if err != nil {
					return nil, err
				}

				result.Values = append(result.Values, value)
			}

			return result, nil
		}
	} else if s.Expression != nil {
		return s.Expression.Evaluate(ctx)
	} else if s.Empty {
		return nil, nil
	}

	panic("unreachable")
}

func (b *FunctionBody) Evaluate(ctx *Context) (types.Type, error) {
	return b.Stmts.Evaluate(ctx)
}

func (f *FunctionDef) Evaluate(ctx *Context) (types.Type, error) {
	var arguments []types.FunctionArgument
	for _, parameter := range f.Parameters {
		arguments = append(
			arguments, types.FunctionArgument{
				TypeHash:   types.GetTypeHash(parameter.Type), // TODO: get type hash with package name
				Name:       parameter.Name,
				IsVariadic: false,
				IsNullable: parameter.IsNullable,
			},
		)
	}

	var returnTypes []types.FunctionReturnType
	for _, returnType := range f.ReturnTypes {
		returnTypes = append(
			returnTypes, types.FunctionReturnType{
				TypeHash:   types.GetTypeHash(*returnType),
				IsNullable: false, // TODO: add it in grammar as '?'
			},
		)
	}

	return types.NewFunctionInstance(
		f.Name,
		arguments,
		func(_ *[]types.Type, kwargs *map[string]types.Type) (types.Type, error) {
			ctx.PushScope(*kwargs)
			result, err := f.Body.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

			ctx.PopScope()
			return result, nil
		},
		returnTypes,
		false,
		ctx.package_,
		"", // TODO: add doc
	), nil
}

func (e *Expression) Evaluate(ctx *Context) (types.Type, error) {
	if e.Assignment != nil {
		return e.Assignment.Evaluate(ctx)
	}

	panic("unreachable")
}

func (c *Constant) Evaluate(ctx *Context) (types.Type, error) {
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
			value, err := expr.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

			list.Values = append(list.Values, value)
		}

		return list, nil
	}

	panic("unreachable")
}

func (a *Assignment) Evaluate(ctx *Context) (types.Type, error) {
	var value types.Type = nil
	if a.Next != nil {
		var err error
		value, err = a.Next.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	return a.LogicalAnd.Evaluate(ctx, value)
}

// Evaluate executes LogicalAnd operation.
// If `valueToSet` is nil, return variable or value from context,
// set a new value or return an error otherwise.
func (a *LogicalAnd) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.AndOp.Caption(), a.LogicalOr, a.Next)
}

func (a *LogicalOr) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.OrOp.Caption(), a.LogicalNot, a.Next)
}

func (a *LogicalNot) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	if a.Comparison != nil {
		return a.Comparison.Evaluate(ctx, valueToSet)
	}

	if a.Next != nil {
		value, err := a.Next.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		return callMethod(value, ops.NotOp.Caption(), &[]types.Type{}, nil)
	}

	panic("unreachable")
}

func (a *Comparison) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
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

func (a *BitwiseOr) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.BitwiseOrOp.Caption(), a.BitwiseXor, a.Next)
}

func (a *BitwiseXor) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.BitwiseXorOp.Caption(), a.BitwiseAnd, a.Next)
}

func (a *BitwiseAnd) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.BitwiseAndOp.Caption(), a.BitwiseShift, a.Next)
}

func (a *BitwiseShift) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	switch a.Op {
	case "<<":
		return evalBinaryOperator(ctx, valueToSet, ops.BitwiseLeftShiftOp.Caption(), a.Addition, a.Next)
	case ">>":
		return evalBinaryOperator(ctx, valueToSet, ops.BitwiseRightShiftOp.Caption(), a.Addition, a.Next)
	default:
		return a.Addition.Evaluate(ctx, valueToSet)
	}
}

func (a *Addition) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	switch a.Op {
	case "+":
		return evalBinaryOperator(ctx, valueToSet, ops.AddOp.Caption(), a.MultiplicationOrMod, a.Next)
	case "-":
		return evalBinaryOperator(ctx, valueToSet, ops.SubOp.Caption(), a.MultiplicationOrMod, a.Next)
	default:
		return a.MultiplicationOrMod.Evaluate(ctx, valueToSet)
	}
}

func (a *MultiplicationOrMod) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
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

func (a *Unary) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
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

func (a *Exponent) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	return evalBinaryOperator(ctx, valueToSet, ops.PowOp.Caption(), a.Primary, a.Next)
}

func (a *Primary) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	if a.Constant != nil {
		if valueToSet != nil {
			// TODO: change to normal description
			return nil, errors.New("unable to set to constant")
		}

		return a.Constant.Evaluate(ctx)
	}

	if a.RandomAccess != nil {
		return a.RandomAccess.Evaluate(ctx, valueToSet)
	}

	if a.Ident != nil {
		if valueToSet != nil {
			err := ctx.setVar(*a.Ident, valueToSet)
			return valueToSet, err
		}

		return ctx.getVar(*a.Ident)
	}

	if a.SubExpression != nil {
		if valueToSet != nil {
			// TODO: change to normal description
			return nil, errors.New("unable to set to subexpression evaluation")
		}

		return a.SubExpression.Evaluate(ctx)
	}

	if a.CallFunc != nil {
		if valueToSet != nil {
			// TODO: change to normal description
			return nil, errors.New("unable to set to function call")
		}

		result, err := a.CallFunc.Evaluate(ctx)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf(
					"  Файл \"%s\", рядок %d, позиція %d\n    %s\n%s",
					a.Pos.Filename, a.CallFunc.Pos.Line, a.CallFunc.Pos.Column, "TODO", err.Error(),
				),
			)
		}

		return result, nil
	}

	panic("unreachable")
}

func (a *RandomAccess) Evaluate(ctx *Context, valueToSet types.Type) (types.Type, error) {
	variable, err := ctx.getVar(a.Ident)
	if err != nil {
		return nil, err
	}

	if valueToSet != nil {
		variable, err = evalSingleSetByIndexOperation(ctx, variable, a.Index, valueToSet)
		if err != nil {
			return nil, err
		}

		return variable, ctx.setVar(a.Ident, variable)
	}

	for _, indexExpression := range a.Index {
		index, err := indexExpression.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		variable, err = evalSingleGetByIndexOperation(variable, index)
		if err != nil {
			return nil, err
		}
	}

	return variable, nil
}

func (a *CallFunc) Evaluate(ctx *Context) (types.Type, error) {
	variable, err := ctx.getVar(a.Ident)
	if err != nil {
		return nil, err
	}

	switch object := variable.(type) {
	case *types.Class:
		callable, err := object.GetAttribute(ops.ConstructorName)
		if err != nil {
			return nil, err
		}

		switch constructor := callable.(type) {
		case *types.FunctionInstance:
			instance, err := object.GetEmptyInstance()
			if err != nil {
				return nil, err
			}

			args := []types.Type{instance}
			kwargs := map[string]types.Type{constructor.Arguments[0].Name: instance}
			for i, expressionArgument := range a.Arguments {
				arg, err := expressionArgument.Evaluate(ctx)
				if err != nil {
					return nil, err
				}

				args = append(args, arg)
				kwargs[constructor.Arguments[i+1].Name] = arg
			}

			if err := types.CheckFunctionArguments(constructor, &args, &kwargs); err != nil {
				return nil, err
			}

			ctx.PushScope(kwargs)
			_, err = constructor.Call(&args, &kwargs)
			if err != nil {
				return nil, err
			}

			ctx.PopScope()
			return args[0], nil
		default:
			return nil, util.ObjectIsNotCallable(a.Ident, callable.GetTypeName())
		}
	case *types.FunctionInstance:
		var args []types.Type
		kwargs := map[string]types.Type{}
		for i, expressionArgument := range a.Arguments {
			arg, err := expressionArgument.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

			args = append(args, arg)
			kwargs[object.Arguments[i].Name] = arg
		}

		if err := types.CheckFunctionArguments(object, &args, &kwargs); err != nil {
			return nil, err
		}

		ctx.PushScope(kwargs)
		res, err := object.Call(&args, &kwargs)
		if err != nil {
			return nil, err
		}

		ctx.PopScope()
		return res, nil
	case types.Instance:
		callable, err := object.GetClass().GetAttribute(ops.CallOperatorName)
		if err != nil {
			return nil, err
		}

		switch callOperator := callable.(type) {
		case *types.FunctionInstance:
			args := []types.Type{variable}
			kwargs := map[string]types.Type{callOperator.Arguments[0].Name: variable}
			for i, expressionArgument := range a.Arguments {
				arg, err := expressionArgument.Evaluate(ctx)
				if err != nil {
					return nil, err
				}

				args = append(args, arg)
				kwargs[callOperator.Arguments[i+1].Name] = arg
			}

			if err := types.CheckFunctionArguments(callOperator, &args, &kwargs); err != nil {
				return nil, err
			}

			ctx.PushScope(kwargs)
			res, err := callOperator.Call(&args, &kwargs)
			if err != nil {
				return nil, err
			}

			ctx.PushScope(kwargs)
			return res, nil
		default:
			return nil, util.ObjectIsNotCallable(a.Ident, callOperator.GetTypeName())
		}
	default:
		return nil, util.ObjectIsNotCallable(a.Ident, object.GetTypeName())
	}
}

func callMethod(object types.Type, funcName string, args *[]types.Type, kwargs *map[string]types.Type) (
	types.Type,
	error,
) {
	attribute, err := object.GetAttribute(funcName)
	if err != nil {
		return nil, util.RuntimeError(err.Error())
	}

	switch function := attribute.(type) {
	case *types.FunctionInstance:
		if len(function.Arguments) == 0 {
			return nil, errors.New(fmt.Sprintf("%s is not a method", function.Representation()))
		}

		*args = append([]types.Type{object}, *args...)
		if kwargs == nil {
			kwargs = &map[string]types.Type{}
		}

		argsLen := len(*args)
		for i := 0; i < argsLen; i++ {
			(*kwargs)[function.Arguments[i].Name] = (*args)[i]
		}

		if err := types.CheckFunctionArguments(function, args, kwargs); err != nil {
			return nil, err
		}

		res, err := function.Call(args, kwargs)
		if err != nil {
			return nil, util.RuntimeError(fmt.Sprintf(err.Error(), funcName))
		}

		return res, nil
	default:
		return nil, util.ObjectIsNotCallable(funcName, attribute.GetTypeName())
	}
}

func mustBool(result types.Type) (types.BoolInstance, error) {
	switch value := result.(type) {
	case types.BoolInstance:
		return value, nil
	default:
		return types.BoolInstance{}, errors.New(
			fmt.Sprintf(
				"'%s' має повертати результат з типом 'логічний', отримано '%s'",
				ops.BoolOperatorName,
				result.GetTypeName(),
			),
		)
	}
}

func evalBinaryOperator(
	ctx *Context,
	valueToSet types.Type,
	operatorName string,
	current OperatorEvaluatable,
	next OperatorEvaluatable,
) (types.Type, error) {
	left, err := current.Evaluate(ctx, valueToSet)
	if err != nil {
		return nil, err
	}

	if !reflect.ValueOf(next).IsNil() {
		right, err := next.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		return callMethod(left, operatorName, &[]types.Type{right}, nil)
	}

	return left, nil
}

func evalUnaryOperator(ctx *Context, operatorName string, operator OperatorEvaluatable) (types.Type, error) {
	if operator != nil {
		value, err := operator.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		return callMethod(value, operatorName, &[]types.Type{}, nil)
	}

	panic("unreachable")
}

func evalSingleGetByIndexOperation(variable types.Type, index types.Type) (types.Type, error) {
	switch iterable := variable.(type) {
	case types.SequentialType:
		switch integerIndex := index.(type) {
		case types.IntegerInstance:
			return iterable.GetElement(integerIndex.Value)
		default:
			return nil, util.RuntimeError("індекси мають бути цілого типу")
		}
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор довільного доступу до об'єкта з типом '%s'",
				variable.GetTypeName(),
			),
		)
	}
}

func evalSingleSetByIndexOperation(
	ctx *Context,
	variable types.Type,
	indices []*Expression,
	value types.Type,
) (types.Type, error) {
	switch iterable := variable.(type) {
	case types.SequentialType:
		index, err := indices[0].Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		switch integerIndex := index.(type) {
		case types.IntegerInstance:
			if len(indices) == 1 {
				return iterable.SetElement(integerIndex.Value, value)
			}

			element, err := iterable.GetElement(integerIndex.Value)
			if err != nil {
				return nil, err
			}

			element, err = evalSingleSetByIndexOperation(ctx, element, indices[1:], value)
			if err != nil {
				return nil, err
			}

			return iterable.SetElement(integerIndex.Value, element)
		default:
			return nil, util.RuntimeError("індекси мають бути цілого типу")
		}
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор довільного доступу до об'єкта з типом '%s'",
				variable.GetTypeName(),
			),
		)
	}
}
