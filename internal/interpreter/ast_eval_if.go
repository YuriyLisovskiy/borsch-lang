package interpreter

import (
	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

func (node *IfStmt) Evaluate(state State, inFunction, inLoop bool) StmtResult {
	if node.Condition != nil {
		condition, err := node.Condition.Evaluate(state, nil)
		if err != nil {
			return StmtResult{Err: err}
		}

		ctx := state.Context()
		conditionValue, err := types2.ToBool(state.Context(), condition)
		if err != nil {
			return StmtResult{Err: err}
		}

		if conditionValue.(types2.Bool) {
			ctx.PushScope(Scope{})
			result := node.Body.Evaluate(state, inFunction, inLoop)
			if result.Err != nil {
				return result
			}

			ctx.PopScope()
			return result
		}

		if len(node.ElseIfStmts) != 0 {
			gotResult := false
			var result StmtResult
			for _, stmt := range node.ElseIfStmts {
				ctx.PushScope(Scope{})
				gotResult, result = stmt.Evaluate(state, inFunction, inLoop)
				ctx.PopScope()
				if result.Interrupt() {
					return result
				}

				if gotResult {
					break
				}
			}

			if gotResult {
				return result
			}
		}

		if node.Else != nil {
			ctx.PushScope(Scope{})
			result := node.Else.Evaluate(state, inFunction, inLoop)
			if result.Err != nil {
				return result
			}

			ctx.PopScope()
			return result
		}

		return StmtResult{}
	}

	panic("unreachable")
}

func (node *ElseIfStmt) Evaluate(state State, inFunction, inLoop bool) (bool, StmtResult) {
	condition, err := node.Condition.Evaluate(state, nil)
	if err != nil {
		return false, StmtResult{Err: err}
	}

	conditionValue, err := types2.ToBool(state.Context(), condition)
	if err != nil {
		return false, StmtResult{Err: err}
	}

	if conditionValue.(types2.Bool) {
		ctx := state.Context()
		ctx.PushScope(Scope{})
		result := node.Body.Evaluate(state, inFunction, inLoop)
		if result.Err != nil {
			return false, result
		}

		ctx.PopScope()
		return true, result
	}

	return false, StmtResult{}
}
