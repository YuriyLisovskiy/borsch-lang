package interpreter

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type StmtState uint8

const (
	StmtNone StmtState = iota
	StmtBreak
	StmtForceReturn
)

type StmtResult struct {
	State StmtState
	Value common.Value
	Err   error
}

// Evaluate executes statement.
// Returns (result value, force stop flag, error)
func (s *Stmt) Evaluate(state common.State, inFunction, inLoop bool) StmtResult {
	switch {
	case s.IfStmt != nil:
		return s.IfStmt.Evaluate(state, inFunction, inLoop)
	case s.LoopStmt != nil:
		return s.LoopStmt.Evaluate(state, inFunction, inLoop)
	case s.Block != nil:
		ctx := state.GetContext()
		ctx.PushScope(Scope{})
		blockResult := s.Block.Evaluate(state, inFunction, inLoop)
		if blockResult.Err != nil {
			return blockResult
		}

		ctx.PopScope()
		return blockResult
	case s.FunctionDef != nil:
		function, err := s.FunctionDef.Evaluate(state, state.GetCurrentPackage().(*types.PackageInstance), nil)
		if err != nil {
			return StmtResult{Err: err}
		}

		return StmtResult{Value: function}
	case s.ClassDef != nil:
		class, err := s.ClassDef.Evaluate(state)
		if err != nil {
			return StmtResult{Err: err}
		}

		return StmtResult{Value: class}
	case s.ReturnStmt != nil:
		if !inFunction {
			return StmtResult{Err: errors.New("'повернути' за межами функції")}
		}

		result, err := s.ReturnStmt.Evaluate(state)
		return StmtResult{Value: result, State: StmtForceReturn, Err: err}
	case s.BreakStmt:
		if !inLoop {
			return StmtResult{Err: errors.New("'перервати' за межами циклу")}
		}

		return StmtResult{State: StmtBreak}
	case s.Assignment != nil:
		result, err := s.Assignment.Evaluate(state)
		return StmtResult{Value: result, Err: err}
	case s.Empty:
		return StmtResult{}
	default:
		panic("unreachable")
	}
}
