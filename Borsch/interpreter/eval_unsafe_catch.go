package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
	"github.com/alecthomas/participle/v2/lexer"
)

func (node *Throw) Position() lexer.Position {
	return node.Pos
}

func (node *Throw) Evaluate(state common.State) StmtResult {
	expression, err := node.Expression.Evaluate(state, nil)
	if err != nil {
		return StmtResult{Err: err}
	}

	expressionClass := expression.(types.ObjectInstance).GetClass()
	if expressionClass == builtin.ErrorClass || expressionClass.HasBase(builtin.ErrorClass) {
		message, err := expression.String(state)
		if err != nil {
			return StmtResult{Err: err}
		}

		return StmtResult{State: StmtThrow, Value: expression, Err: utilities.NewRuntimeStatementError(message, node)}
	}

	return StmtResult{
		Err: state.RuntimeError(
			fmt.Sprintf(
				"помилки мають наслідувати клас '%s'",
				builtin.ErrorClass.Name,
			),
			node,
		),
	}
}

func (node *Unsafe) Position() lexer.Position {
	return node.Pos
}

func (node *Unsafe) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	result := node.Stmts.Evaluate(state, inFunction, inLoop)
	if result.State != StmtThrow {
		return result
	}

	for _, catchBlock := range node.CatchBlocks {
		blockResult, caught := catchBlock.Evaluate(state, result.Value, inFunction, inLoop)
		if blockResult.Interrupt() {
			if blockResult.State == StmtThrow {
				err := blockResult.Err.(utilities.RuntimeStatementError)
				blockResult.Err = state.RuntimeError(err.Error(), err.Statement())
			}

			return blockResult
		}

		if caught {
			return blockResult
		}
	}

	result.Err = state.RuntimeError(result.Err.Error(), node.Stmts.GetCurrentStmt())
	return result
}

func (node *Catch) Evaluate(state common.State, exception common.Value, inFunction, inLoop bool) (StmtResult, bool) {
	errorToCatch, err := node.ErrorType.Evaluate(state, nil, nil)
	if err != nil {
		return StmtResult{Err: err}, false
	}

	if _, ok := errorToCatch.(*types.Class); !ok {
		str, err := errorToCatch.String(state)
		if err != nil {
			return StmtResult{Err: err}, false
		}

		return StmtResult{Err: state.RuntimeError(fmt.Sprintf("об'єкт '%s' не є класом", str), node)}, false
	}

	generatedErrorClass := exception.(types.ObjectInstance).GetClass()
	if generatedErrorClass == errorToCatch || generatedErrorClass.HasBase(errorToCatch.(*types.Class)) {
		ctx := state.GetContext()
		ctx.PushScope(Scope{node.ErrorVar: exception})
		result := node.Stmts.Evaluate(state, inFunction, inLoop)
		if result.Err != nil {
			return result, false
		}

		ctx.PopScope()
		return result, true
	}

	if !errorToCatch.(*types.Class).HasBase(builtin.ErrorClass) {
		return StmtResult{
			Err: state.RuntimeError(
				fmt.Sprintf(
					"перехоплення помилок, які не наслідують клас '%s' заборонено",
					builtin.ErrorClass.Name,
				),
				node,
			),
		}, false
	}

	return StmtResult{}, false
}
