package interpreter

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

type StmtState uint8

func (s StmtState) String() string {
	switch s {
	case StmtNone:
		return "StmtNone"
	case StmtBreak:
		return "StmtBreak"
	case StmtForceReturn:
		return "StmtForceReturn"
	case StmtThrow:
		return "StmtThrown"
	default:
		return ""
	}
}

const (
	StmtNone StmtState = iota
	StmtBreak
	StmtForceReturn
	StmtThrow
)

type StmtResult struct {
	State StmtState
	Value types.Object
	Err   error
}

// Interrupt returns true when statement result contains
// and error, or has StmtForceReturn or StmtBreak state.
func (r StmtResult) Interrupt() bool {
	if r.Err != nil {
		return true
	}

	switch r.State {
	case StmtForceReturn, StmtBreak:
		return true
	}

	return false
}

// Evaluate executes statement.
// Returns (result value, force stop flag, error)
func (node *Stmt) Evaluate(state State, inFunction, inLoop bool) StmtResult {
	switch {
	case node.Throw != nil:
		return node.Throw.Evaluate(state)
	case node.Unsafe != nil:
		return node.Unsafe.Evaluate(state, inFunction, inLoop)
	case node.IfStmt != nil:
		return node.IfStmt.Evaluate(state, inFunction, inLoop)
	case node.LoopStmt != nil:
		return node.LoopStmt.Evaluate(state, inFunction, inLoop)
	case node.Block != nil:
		ctx := state.Context()
		ctx.PushScope(Scope{})
		blockResult := node.Block.Evaluate(state, inFunction, inLoop)
		if blockResult.Err != nil {
			return blockResult
		}

		ctx.PopScope()
		return blockResult
	case node.FunctionDef != nil:
		function, err := node.FunctionDef.Evaluate(state, state.Package().(*types.Package), false, nil)
		if err != nil {
			return StmtResult{Err: err}
		}

		return StmtResult{Value: function}
	case node.ClassDef != nil:
		class, err := node.ClassDef.Evaluate(state)
		if err != nil {
			return StmtResult{Err: err}
		}

		return StmtResult{Value: class}
	case node.ReturnStmt != nil:
		if !inFunction {
			return StmtResult{Err: errors.New("'повернути' за межами функції")}
		}

		result, err := node.ReturnStmt.Evaluate(state)
		return StmtResult{Value: result, State: StmtForceReturn, Err: err}
	case node.BreakStmt:
		if !inLoop {
			return StmtResult{Err: errors.New("'перервати' за межами циклу")}
		}

		return StmtResult{State: StmtBreak}
	case node.Assignment != nil:
		result, err := node.Assignment.Evaluate(state)
		return StmtResult{Value: result, Err: err}
	case node.Empty:
		return StmtResult{}
	default:
		panic("unreachable")
	}
}
