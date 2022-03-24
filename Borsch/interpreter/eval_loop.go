package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (node *LoopStmt) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	if node.RangeBasedLoop != nil {
		return node.RangeBasedLoop.Evaluate(state, node.Body, inFunction)
	} else if node.ConditionalLoop != nil {
		return node.ConditionalLoop.Evaluate(state, node.Body, inFunction)
	}

	return evalInfiniteLoop(state, node.Body, inFunction)
}

func (node *RangeBasedLoop) Evaluate(state common.State, body *BlockStmts, inFunction bool) StmtResult {
	leftBound, err := getBound(state, node.LeftBound, "ліва")
	if err != nil {
		return StmtResult{Err: err}
	}

	rightBound, err := getBound(state, node.RightBound, "права")
	if err != nil {
		return StmtResult{Err: err}
	}

	ctx := state.GetContext()
	for leftBound < rightBound {
		ctx.PushScope(Scope{node.Variable.String(): types.Int(leftBound)})
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

func (node *ConditionalLoop) Evaluate(state common.State, body *BlockStmts, inFunction bool) StmtResult {
	ctx := state.GetContext()
	for {
		condition, err := node.Condition.Evaluate(state, nil)
		if err != nil {
			return StmtResult{Err: err}
		}

		conditionValue, err := types.ToBool(state.GetContext(), condition)
		if err != nil {
			return StmtResult{Err: err}
		}

		if !conditionValue.(types.Bool) {
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

func getBound(state common.State, bound *Expression, boundName string) (int64, error) {
	return mustInt(
		state, bound, func(t types.Object) string {
			return fmt.Sprintf("%s межа має бути цілого типу, отримано %s", boundName, t.Class().Name)
		},
	)
}

func evalInfiniteLoop(state common.State, body *BlockStmts, inFunction bool) StmtResult {
	ctx := state.GetContext()
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
