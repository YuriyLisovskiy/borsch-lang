package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

func (node *Throw) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
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

		return StmtResult{State: StmtThrown, Value: expression, Err: errors.New(message)}
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

func (node *Try) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	result := node.Stmts.Evaluate(state, inFunction, inLoop)
	if result.Err != nil {
		if result.State != StmtThrown {
			return StmtResult{Err: result.Err}
		}

		for _, catch := range node.CatchBlocks {
			blockResult, caught := catch.Evaluate(state, result.Value, inFunction, inLoop)
			if blockResult.Err != nil {
				return blockResult
			}

			switch blockResult.State {
			case StmtForceReturn, StmtBreak:
				return blockResult
			}

			if caught {
				result.Err = nil
				return result
			}
		}
	}

	return result
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
