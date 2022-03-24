package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (node *Throw) Evaluate(state common.State) StmtResult {
	// expression, err := node.Expression.Evaluate(state, nil)
	// if err != nil {
	// 	return StmtResult{Err: err}
	// }

	// expressionClass := expression.Class()
	// if expressionClass == builtin.ErrorClass || expressionClass.HasBase(builtin.ErrorClass) {
	// 	message, err := expression.String(state)
	// 	if err != nil {
	// 		return StmtResult{Err: err}
	// 	}
	//
	// 	state.Trace(node, "")
	// 	return StmtResult{State: StmtThrow, Value: expression, Err: utilities.NewRuntimeStatementError(message, node)}
	// }

	return StmtResult{
		Err: state.RuntimeError(
			fmt.Sprintf(
				"помилки мають наслідувати клас '%s'",
				// builtin.ErrorClass.Name,
			),
			node,
		),
	}
}

func (node *Unsafe) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	result := node.Stmts.Evaluate(state, inFunction, inLoop)
	if result.State != StmtThrow {
		return result
	}

	for _, catchBlock := range node.CatchBlocks {
		blockResult, caught := catchBlock.Evaluate(state, result.Value, inFunction, inLoop)
		if blockResult.Interrupt() {
			return blockResult
		}

		if caught {
			return blockResult
		}
	}

	state.Trace(node.Stmts.GetCurrentStmt(), "")
	return result
}

func (node *Catch) Evaluate(state common.State, exception types.Object, inFunction, inLoop bool) (StmtResult, bool) {
	errorToCatch, err := node.ErrorType.Evaluate(state, nil, nil)
	if err != nil {
		return StmtResult{Err: err}, false
	}

	if _, ok := errorToCatch.(*types.Class); !ok {
		str, err := types.ToGoString(state.GetContext(), errorToCatch)
		if err != nil {
			return StmtResult{Err: err}, false
		}

		return StmtResult{Err: state.RuntimeError(fmt.Sprintf("об'єкт '%s' не є класом", str), node)}, false
	}

	generatedErrorClass := exception.Class()
	errorToCatchClass := errorToCatch.(*types.Class)
	if shouldCatch(generatedErrorClass, errorToCatchClass) {
		state.PopTrace()
		return node.catch(state, exception, inFunction, inLoop)
	}

	// if !errorToCatchClass.HasBase(builtin.ErrorClass) {
	// 	return StmtResult{
	// 		Err: state.RuntimeError(
	// 			fmt.Sprintf(
	// 				"перехоплення помилок, які не наслідують клас '%s' заборонено",
	// 				builtin.ErrorClass.Name,
	// 			),
	// 			node,
	// 		),
	// 	}, false
	// }

	return StmtResult{}, false
}

func (node *Catch) catch(state common.State, err types.Object, inFunction, inLoop bool) (StmtResult, bool) {
	ctx := state.GetContext()
	ctx.PushScope(Scope{node.ErrorVar.String(): err})
	result := node.Stmts.Evaluate(state, inFunction, inLoop)
	if result.Err != nil {
		return result, false
	}

	ctx.PopScope()
	return result, true
}

func shouldCatch(generated, toCatch *types.Class) bool {
	return generated == toCatch || generated.HasBase(toCatch)
}
