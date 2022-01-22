package grammar

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (s *IfStmt) Evaluate(ctx common.Context, inFunction bool) (
	common.Type,
	bool,
	error,
) {
	if s.Condition != nil {
		condition, err := s.Condition.Evaluate(ctx, nil)
		if err != nil {
			return nil, false, err
		}

		if condition.AsBool(ctx) {
			ctx.PushScope(Scope{})
			result, forceReturn, err := s.Body.Evaluate(ctx, inFunction)
			if err != nil {
				return nil, false, err
			}

			ctx.PopScope()
			return result, forceReturn, nil
		}

		if len(s.ElseIfStmts) != 0 {
			gotResult := false
			var result common.Type = nil
			var err error = nil
			for _, stmt := range s.ElseIfStmts {
				ctx.PushScope(Scope{})
				var forceReturn bool
				gotResult, result, forceReturn, err = stmt.Evaluate(ctx, inFunction)
				if err != nil {
					return nil, false, err
				}

				ctx.PopScope()
				if forceReturn {
					return result, true, nil
				}

				if gotResult {
					break
				}
			}

			if gotResult {
				return result, false, nil
			}
		}

		if s.Else != nil {
			ctx.PushScope(Scope{})
			result, forceReturn, err := s.Else.Evaluate(ctx, inFunction)
			if err != nil {
				return nil, false, err
			}

			ctx.PopScope()
			return result, forceReturn, nil
		}

		return nil, false, nil
	}

	return nil, false, errors.New("interpreter: condition is nil")
}

func (s *ElseIfStmt) Evaluate(ctx common.Context, inFunction bool) (
	bool,
	common.Type,
	bool,
	error,
) {
	condition, err := s.Condition.Evaluate(ctx, nil)
	if err != nil {
		return false, nil, false, err
	}

	if condition.AsBool(ctx) {
		ctx.PushScope(Scope{})
		result, forceReturn, err := s.Body.Evaluate(ctx, inFunction)
		if err != nil {
			return false, nil, false, err
		}

		ctx.PopScope()
		return true, result, forceReturn, nil
	}

	return false, nil, false, nil
}
