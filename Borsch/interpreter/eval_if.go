package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (s *IfStmt) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	if s.Condition != nil {
		condition, err := s.Condition.Evaluate(state, nil)
		if err != nil {
			return StmtResult{Err: err}
		}

		ctx := state.GetContext()
		if condition.AsBool(state) {
			ctx.PushScope(Scope{})
			result := s.Body.Evaluate(state, inFunction, inLoop)
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
				gotResult, result = stmt.Evaluate(state, inFunction, inLoop)
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
			result := s.Else.Evaluate(state, inFunction, inLoop)
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

func (s *ElseIfStmt) Evaluate(state common.State, inFunction, inLoop bool) (bool, StmtResult) {
	condition, err := s.Condition.Evaluate(state, nil)
	if err != nil {
		return false, StmtResult{Err: err}
	}

	if condition.AsBool(state) {
		ctx := state.GetContext()
		ctx.PushScope(Scope{})
		result := s.Body.Evaluate(state, inFunction, inLoop)
		if result.Err != nil {
			return false, result
		}

		ctx.PopScope()
		return true, result
	}

	return false, StmtResult{}
}
