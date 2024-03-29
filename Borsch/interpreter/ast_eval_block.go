package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

func (node *Throw) Evaluate(state State) StmtResult {
	expressionObj, err := node.Expression.Evaluate(state, nil)
	if err != nil {
		return StmtResult{Err: err}
	}

	expressionClass := expressionObj.Class()
	if expressionClass == types.ErrorClass || types.ErrorClass.IsBaseOf(expressionClass) {
		state.Trace(node, "")
		stmtResult := StmtResult{
			State: StmtThrow,
			Value: expressionObj,
			// Err:   utilities.NewRuntimeStatementError(message, node),
		}

		if execErr, ok := expressionObj.(types.LangException); ok {
			stmtResult.Err = execErr
		} else {
			// TODO: remove this branch in future!

			message, err := types.ToGoString(state.Context(), expressionObj)
			if err != nil {
				return StmtResult{Err: err}
			}

			stmtResult.Err = utilities.NewRuntimeStatementError(message, node)
		}

		return stmtResult
	}

	return StmtResult{
		Err: state.RuntimeError(fmt.Sprintf("помилки мають наслідувати клас '%s'", types.ErrorClass.Name), node),
	}
}

func (node *Block) Evaluate(state State, inFunction, inLoop bool) StmtResult {
	result := node.Stmts.Evaluate(state, inFunction, inLoop)
	if result.State != StmtThrow {
		if result.Err == nil {
			return result
		}

		langErr, ok := result.Err.(types.LangException)
		if !ok {
			return result
		}

		result.Value = langErr
	}

	if len(node.CatchBlocks) > 0 {
		for _, catchBlock := range node.CatchBlocks {
			blockResult, caught := catchBlock.Evaluate(state, result.Value, inFunction, inLoop)
			if blockResult.Interrupt() {
				return blockResult
			}

			if caught {
				return blockResult
			}
		}
	}

	state.Trace(node.Stmts.GetCurrentStmt(), "")
	return result
}

func (node *Catch) Evaluate(state State, exception types.Object, inFunction, inLoop bool) (
	StmtResult,
	bool,
) {
	errorToCatch, err := node.ErrorType.Evaluate(state, nil, nil)
	if err != nil {
		return StmtResult{Err: err}, false
	}

	if _, ok := errorToCatch.(*types.Class); !ok {
		str, err := types.ToGoString(state.Context(), errorToCatch)
		if err != nil {
			return StmtResult{Err: err}, false
		}

		return StmtResult{Err: state.RuntimeError(fmt.Sprintf("об'єкт '%s' не є класом", str), node)}, false
	}

	generatedErrorClass := exception.Class()
	errorToCatchClass := errorToCatch.(*types.Class)
	if shouldCatch(generatedErrorClass, errorToCatchClass) {
		// TODO: check if stacktrace is ok when the line below is not used!
		// state.PopTrace()
		return node.catch(state, exception, inFunction, inLoop)
	}

	if !types.ErrorClass.IsBaseOf(errorToCatchClass) {
		return StmtResult{
			Err: state.RuntimeError(
				fmt.Sprintf(
					"перехоплення помилок, які не наслідують клас '%s' заборонено",
					types.ErrorClass.Name,
				), node,
			),
		}, false
	}

	return StmtResult{}, false
}

func (node *Catch) catch(state State, err types.Object, inFunction, inLoop bool) (StmtResult, bool) {
	ctx := state.Context()
	ctx.PushScope(Scope{node.ErrorVar.String(): err})
	result := node.Stmts.Evaluate(state, inFunction, inLoop)
	if result.Err != nil {
		return result, false
	}

	ctx.PopScope()
	return result, true
}

func shouldCatch(generated, toCatch *types.Class) bool {
	return generated == toCatch || toCatch.IsBaseOf(generated)
}
