package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (s *IfStmt) Evaluate(ctx common.Context, inFunction, inLoop bool) StmtResult {
	if s.Condition != nil {
		condition, err := s.Condition.Evaluate(ctx, nil)
		if err != nil {
			return StmtResult{Err: err}
		}

		if condition.AsBool(ctx) {
			ctx.PushScope(Scope{})
			result := s.Body.Evaluate(ctx, inFunction, inLoop)
			if err != nil {
				return result
			}

			ctx.PopScope()
			return result
		}

		if len(s.ElseIfStmts) != 0 {
			gotResult := false
			var result StmtResult
			var err error = nil
			for _, stmt := range s.ElseIfStmts {
				ctx.PushScope(Scope{})
				gotResult, result = stmt.Evaluate(ctx, inFunction, inLoop)
				if err != nil {
					return result
				}

				ctx.PopScope()
				switch result.State {
				case StmtForceReturn, StmtBreak:
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

		if s.Else != nil {
			ctx.PushScope(Scope{})
			result := s.Else.Evaluate(ctx, inFunction, inLoop)
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

func (s *ElseIfStmt) Evaluate(ctx common.Context, inFunction, inLoop bool) (bool, StmtResult) {
	condition, err := s.Condition.Evaluate(ctx, nil)
	if err != nil {
		return false, StmtResult{Err: err}
	}

	if condition.AsBool(ctx) {
		ctx.PushScope(Scope{})
		result := s.Body.Evaluate(ctx, inFunction, inLoop)
		if result.Err != nil {
			return false, result
		}

		ctx.PopScope()
		return true, result
	}

	return false, StmtResult{}
}
