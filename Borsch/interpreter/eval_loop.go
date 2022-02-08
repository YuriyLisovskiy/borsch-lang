package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (node *LoopStmt) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	if node.Body == nil {
		panic("unreachable")
	}

	if node.RangeBasedLoop != nil {
		return node.RangeBasedLoop.Evaluate(state, node.Body, inFunction, inLoop)
	} else if node.ConditionalLoop != nil {
		return node.ConditionalLoop.Evaluate(state, node.Body, inFunction, inLoop)
	}

	panic("unreachable")
}

func (node *RangeBasedLoop) Evaluate(state common.State, body *BlockStmts, inFunction, inLoop bool) StmtResult {
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
		ctx.PushScope(Scope{node.Variable: types.NewIntegerInstance(leftBound)})
		result := body.Evaluate(state, inFunction, true)
		if result.Err != nil {
			return result
		}

		ctx.PopScope()
		switch result.State {
		case StmtForceReturn:
			return result
		case StmtBreak:
			result.State = StmtNone
			return result
		}

		leftBound += 1
	}

	return StmtResult{}
}

func (node *ConditionalLoop) Evaluate(state common.State, body *BlockStmts, inFunction, inLoop bool) StmtResult {
	ctx := state.GetContext()
	for {
		condition, err := node.Condition.Evaluate(state, nil)
		if err != nil {
			return StmtResult{Err: err}
		}

		conditionValue, err := condition.AsBool(state)
		if err != nil {
			return StmtResult{Err: err}
		}

		if !conditionValue {
			break
		}

		ctx.PushScope(Scope{})
		result := body.Evaluate(state, inFunction, true)
		if result.Err != nil {
			return result
		}

		ctx.PopScope()
		switch result.State {
		case StmtForceReturn:
			return result
		case StmtBreak:
			result.State = StmtNone
			return result
		}
	}

	return StmtResult{}
}

func getBound(state common.State, bound *Expression, boundName string) (int64, error) {
	return mustInt(
		state, bound, func(t common.Value) string {
			return fmt.Sprintf("%s межа має бути цілого типу, отримано %s", boundName, t.GetTypeName())
		},
	)
}
