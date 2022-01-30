package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
)

func (l *LoopStmt) Evaluate(ctx common.Context, inFunction, inLoop bool) StmtResult {
	if l.Body == nil {
		panic("unreachable")
	}

	if l.RangeBasedLoop != nil {
		return l.RangeBasedLoop.Evaluate(ctx, l.Body, inFunction, inLoop)
	}

	return l.ConditionalLoop.Evaluate(ctx, l.Body, inFunction, inLoop)
}

func (l *RangeBasedLoop) Evaluate(ctx common.Context, body *BlockStmts, inFunction, inLoop bool) StmtResult {
	leftBound, err := getBound(ctx, l.LeftBound, "ліва")
	if err != nil {
		return StmtResult{Err: err}
	}

	rightBound, err := getBound(ctx, l.RightBound, "права")
	if err != nil {
		return StmtResult{Err: err}
	}

	for leftBound < rightBound {
		ctx.PushScope(Scope{l.Variable: types.NewIntegerInstance(leftBound)})
		result := body.Evaluate(ctx, inFunction, true)
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

func (l *ConditionalLoop) Evaluate(ctx common.Context, body *BlockStmts, inFunction, inLoop bool) StmtResult {
	for {
		condition, err := l.Condition.Evaluate(ctx, nil)
		if err != nil {
			return StmtResult{Err: err}
		}

		if !condition.AsBool(ctx) {
			break
		}

		ctx.PushScope(Scope{})
		result := body.Evaluate(ctx, inFunction, true)
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

func getBound(ctx common.Context, bound *Expression, boundName string) (int64, error) {
	return mustInt(
		ctx, bound, func(t common.Type) string {
			return fmt.Sprintf("%s межа має бути цілого типу, отримано %s", boundName, t.GetTypeName())
		},
	)
}
