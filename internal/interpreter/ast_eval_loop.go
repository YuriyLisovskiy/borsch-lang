package interpreter

import (
	"errors"
	"fmt"

	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

func (node *LoopStmt) Evaluate(state State, inFunction, inLoop bool) StmtResult {
	if node.RangeBasedLoop != nil {
		return node.RangeBasedLoop.Evaluate(state, node.Body, inFunction)
	} else if node.ConditionalLoop != nil {
		return node.ConditionalLoop.Evaluate(state, node.Body, inFunction)
	}

	return evalInfiniteLoop(state, node.Body, inFunction)
}

func (node *RangeBasedLoop) Evaluate(state State, body *BlockStmts, inFunction bool) StmtResult {
	leftBound, err := getBound(state, node.LeftBound, "ліва")
	if err != nil {
		return StmtResult{Err: err}
	}

	rightBound, err := getBound(state, node.RightBound, "права")
	if err != nil {
		return StmtResult{Err: err}
	}

	ctx := state.Context()
	for leftBound < rightBound {
		ctx.PushScope(Scope{node.Variable.String(): types2.Int(leftBound)})
		result := body.Evaluate(state, inFunction, true)
		ctx.PopScope()
		if result.Interrupt() {
			if result.State == StmtBreak {
				result.State = StmtNone
			}

			return result
		}

		leftBound += 1
	}

	return StmtResult{}
}

func (node *ConditionalLoop) Evaluate(state State, body *BlockStmts, inFunction bool) StmtResult {
	ctx := state.Context()
	for {
		condition, err := node.Condition.Evaluate(state, nil)
		if err != nil {
			return StmtResult{Err: err}
		}

		conditionValue, err := types2.ToBool(state.Context(), condition)
		if err != nil {
			return StmtResult{Err: err}
		}

		if !conditionValue.(types2.Bool) {
			break
		}

		ctx.PushScope(Scope{})
		result := body.Evaluate(state, inFunction, true)
		ctx.PopScope()
		if result.Interrupt() {
			if result.State == StmtBreak {
				result.State = StmtNone
			}

			return result
		}
	}

	return StmtResult{}
}

func getBound(state State, bound *Expression, boundName string) (types2.Int, error) {
	return mustInt(
		state, bound, func(t types2.Object) error {
			return errors.New(fmt.Sprintf("%s межа має бути цілого типу, отримано %s", boundName, t.Class().Name))
		},
	)
}

func evalInfiniteLoop(state State, body *BlockStmts, inFunction bool) StmtResult {
	ctx := state.Context()
	for {
		ctx.PushScope(Scope{})
		result := body.Evaluate(state, inFunction, true)
		ctx.PopScope()
		if result.Interrupt() {
			if result.State == StmtBreak {
				result.State = StmtNone
			}

			return result
		}
	}
}
