package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
)

func (l *LoopStmt) Evaluate(ctx common.Context, inFunction bool) (common.Type, bool, error) {
	if l.RangeBasedLoop == nil || l.Body == nil {
		panic("unreachable")
	}

	return l.RangeBasedLoop.Evaluate(ctx, l.Body, inFunction)
}

func (s *RangeBasedLoop) Evaluate(ctx common.Context, body *BlockStmts, inFunction bool) (common.Type, bool, error) {
	leftBound, err := getBound(ctx, s.LeftBound, "ліва")
	if err != nil {
		return nil, false, err
	}

	rightBound, err := getBound(ctx, s.RightBound, "права")
	if err != nil {
		return nil, false, err
	}

	for leftBound < rightBound {
		ctx.PushScope(Scope{s.Variable: types.NewIntegerInstance(leftBound)})
		result, forceReturn, err := body.Evaluate(ctx, inFunction)
		if err != nil {
			return nil, false, err
		}

		ctx.PopScope()
		if forceReturn {
			return result, forceReturn, nil
		}

		leftBound += 1
	}

	return nil, false, nil
}

func getBound(ctx common.Context, bound *Expression, boundName string) (int64, error) {
	return mustInt(
		ctx, bound, func(t common.Type) string {
			return fmt.Sprintf("%s межа має бути цілого типу, отримано %s", boundName, t.GetTypeName())
		},
	)
}
