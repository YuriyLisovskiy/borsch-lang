package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
)

func (l *LoopStmt) Evaluate(ctx common.Context, inFunction, inLoop bool) StmtResult {
	if l.RangeBasedLoop == nil || l.Body == nil {
		panic("unreachable")
	}

	return l.RangeBasedLoop.Evaluate(ctx, l.Body, inFunction, inLoop)
}

func (s *RangeBasedLoop) Evaluate(ctx common.Context, body *BlockStmts, inFunction, inLoop bool) StmtResult {
	leftBound, err := getBound(ctx, s.LeftBound, "ліва")
	if err != nil {
		return StmtResult{Err: err}
	}

	rightBound, err := getBound(ctx, s.RightBound, "права")
	if err != nil {
		return StmtResult{Err: err}
	}

	for leftBound < rightBound {
		ctx.PushScope(Scope{s.Variable: types.NewIntegerInstance(leftBound)})
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

func getBound(ctx common.Context, bound *Expression, boundName string) (int64, error) {
	return mustInt(
		ctx, bound, func(t common.Type) string {
			return fmt.Sprintf("%s межа має бути цілого типу, отримано %s", boundName, t.GetTypeName())
		},
	)
}
