package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (l *LoopStmt) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	if l.Body == nil {
		panic("unreachable")
	}

	if l.RangeBasedLoop != nil {
		return l.RangeBasedLoop.Evaluate(state, l.Body, inFunction, inLoop)
	}

	return l.ConditionalLoop.Evaluate(state, l.Body, inFunction, inLoop)
}

func (l *RangeBasedLoop) Evaluate(state common.State, body *BlockStmts, inFunction, inLoop bool) StmtResult {
	leftBound, err := getBound(state, l.LeftBound, "ліва")
	if err != nil {
		return StmtResult{Err: err}
	}

	rightBound, err := getBound(state, l.RightBound, "права")
	if err != nil {
		return StmtResult{Err: err}
	}

	ctx := state.GetContext()
	for leftBound < rightBound {
		ctx.PushScope(Scope{l.Variable: types.NewIntegerInstance(leftBound)})
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

func (l *ConditionalLoop) Evaluate(state common.State, body *BlockStmts, inFunction, inLoop bool) StmtResult {
	ctx := state.GetContext()
	for {
		condition, err := l.Condition.Evaluate(state, nil)
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
