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
		Err: utilities.RuntimeError(
			fmt.Sprintf(
				"помилки мають наслідувати клас '%s'",
				builtin.ErrorClass.Name,
			),
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
				node.trace(state, &blockResult)
			}

			return blockResult
		}

		if caught {
			return blockResult
		}
	}

	return result
}

func (node *Unsafe) trace(state common.State, result *StmtResult) {
	var stmt utilities.ErrorStatement = node
	if err, ok := result.Err.(utilities.RuntimeStatementError); ok {
		stmt = err.Statement()
	}

	state.GetInterpreter().Trace(stmt.Position(), "", stmt.String())
}

func (node *Catch) Evaluate(state common.State, exception common.Value, inFunction, inLoop bool) (StmtResult, bool) {
	errorClass, err := state.GetContext().GetClass(node.ErrorType)
	if err != nil {
		return StmtResult{Err: err}, false
	}

	targetErrorClass := exception.(types.ObjectInstance).GetClass()
	if targetErrorClass == errorClass || targetErrorClass.HasBase(errorClass.(*types.Class)) {
		ctx := state.GetContext()
		ctx.PushScope(Scope{node.ErrorVar: exception})
		result := node.Stmts.Evaluate(state, inFunction, inLoop)
		if result.Err != nil {
			return result, false
		}

		ctx.PopScope()
		return result, true
	}

	return StmtResult{}, false
}
