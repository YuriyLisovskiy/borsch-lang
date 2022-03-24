package interpreter

import (
	"fmt"
	"testing"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func makeThrowStmt(name *Ident) *Throw {
	return &Throw{
		Expression: &Expression{
			LogicalAnd: &LogicalAnd{
				LogicalOr: &LogicalOr{
					LogicalNot: &LogicalNot{
						Comparison: &Comparison{
							BitwiseOr: &BitwiseOr{
								BitwiseXor: &BitwiseXor{
									BitwiseAnd: &BitwiseAnd{
										BitwiseShift: &BitwiseShift{
											Addition: &Addition{
												MultiplicationOrMod: &MultiplicationOrMod{
													Unary: &Unary{
														Exponent: &Exponent{
															Primary: &Primary{
																AttributeAccess: &AttributeAccess{
																	IdentOrCall: &IdentOrCall{
																		Ident: name,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type testInterpreter struct {
}

func (i *testInterpreter) Import(common.State, string) (types.Object, error) {
	return nil, nil
}

func (i *testInterpreter) StackTrace() *common.StackTrace {
	st := &common.StackTrace{}
	st.Push(&common.TraceRow{})
	return st
}

func TestThrow_EvaluateSuccess(t *testing.T) {
	errorIdent := Ident("err")
	throwNode := makeThrowStmt(&errorIdent)
	errMessage := "This is an error"
	exc, _ := builtin.NewErrorInstance(errMessage)
	state := StateImpl{
		interpreter: &testInterpreter{},
		context: &ContextImpl{
			scopes: []map[string]types.Object{
				{errorIdent.String(): exc},
			},
		},
	}

	result := throwNode.Evaluate(&state)
	if result.State != StmtThrow {
		assertionFailed(t, result.State.String(), StmtThrow.String())
	}

	errMessage = fmt.Sprintf("%s: %s", builtin.ErrorClass.GetName(), errMessage)
	if result.Err.Error() != errMessage {
		assertionFailed(t, result.Err.Error(), errMessage)
	}

	if result.Value != exc {
		t.Error("value is not exception")
	}
}

func TestThrow_EvaluateFail_NotAnErrorInstance(t *testing.T) {
	errorIdent := Ident("err")
	throwNode := makeThrowStmt(&errorIdent)
	errMessage := "This is an error"
	state := StateImpl{
		interpreter: &testInterpreter{},
		context: &ContextImpl{
			scopes: []map[string]types.Object{
				{errorIdent.String(): types.NewStringInstance(errMessage)},
			},
		},
	}

	result := throwNode.Evaluate(&state)
	if result.State != StmtNone {
		assertionFailed(t, result.State.String(), StmtNone.String())
	}

	if result.Err == nil {
		t.Error("error is nil")
	}

	if result.Value != nil {
		t.Error("value is not nil")
	}
}

func TestUnsafe_EvaluateNoThrow(t *testing.T) {
	errorIdent := Ident("Error")
	unsafe := &Unsafe{
		Stmts: &BlockStmts{
			Stmts:   []*Stmt{},
			stmtPos: 0,
		},
		CatchBlocks: []*Catch{
			{
				ErrorVar: "e",
				ErrorType: &AttributeAccess{
					IdentOrCall: &IdentOrCall{
						Ident: &errorIdent,
					},
				},
				Stmts: &BlockStmts{
					Stmts:   []*Stmt{},
					stmtPos: 0,
				},
			},
		},
	}

	errMessage := "This is an error"
	err, _ := builtin.NewErrorInstance(errMessage)
	state := StateImpl{
		interpreter: &testInterpreter{},
		context: &ContextImpl{
			scopes: []map[string]types.Object{
				{errorIdent.String(): err},
			},
		},
	}
	result := unsafe.Evaluate(&state, false, false)
	if result.State == StmtThrow {
		t.Error("state is StmtThrow")
	}

	if result.Err != nil {
		t.Error("error is not nil")
	}

	if _, ok := result.Value.(types.NilInstance); !ok {
		t.Error("result value is not NilInstance")
	}
}

func TestUnsafe_EvaluateThrownAndNotCaught(t *testing.T) {
	errorClass := types.Class{
		Name:  "CustomError",
		Class: types.TypeClass,
		Bases: []*types.Class{builtin.ErrorClass},
	}
	errorClassName := Ident(errorClass.Name)
	errorIdent := Ident("Error")
	unsafe := &Unsafe{
		Stmts: &BlockStmts{
			Stmts:   []*Stmt{{Throw: makeThrowStmt(&errorIdent)}},
			stmtPos: 0,
		},
		CatchBlocks: []*Catch{
			{
				ErrorVar: "e",
				ErrorType: &AttributeAccess{
					IdentOrCall: &IdentOrCall{
						Ident: &errorClassName,
					},
				},
				Stmts: &BlockStmts{
					Stmts:   []*Stmt{},
					stmtPos: 0,
				},
			},
		},
	}

	errMessage := "This is an error"
	err, _ := builtin.NewErrorInstance(errMessage)
	state := StateImpl{
		interpreter: &testInterpreter{},
		context: &ContextImpl{
			scopes: []map[string]types.Object{
				{errorClass.Name: &errorClass},
				{errorIdent.String(): err},
			},
		},
	}
	result := unsafe.Evaluate(&state, false, false)
	if result.State != StmtThrow {
		assertionFailed(t, result.State.String(), StmtThrow.String())
	}

	if result.Err == nil {
		t.Error("error is nil")
	}

	errMessage = fmt.Sprintf("%s: %s", builtin.ErrorClass.GetName(), errMessage)
	if result.Err.Error() != errMessage {
		assertionFailed(t, result.Err.Error(), errMessage)
	}

	if result.Value == nil {
		t.Error("value is nil")
	}

	if _, ok := result.Value.(*builtin.ErrorInstance); !ok {
		t.Error("result value is not ErrorInstance")
	}

	if result.Value.(*builtin.ErrorInstance) != err {
		t.Error("result value is not expected error")
	}
}

func TestUnsafe_EvaluateThrownAndCaught(t *testing.T) {
	errorIdent := Ident("Error")
	errorClassName := Ident(builtin.ErrorClass.Name)
	unsafe := &Unsafe{
		Stmts: &BlockStmts{
			Stmts:   []*Stmt{{Throw: makeThrowStmt(&errorIdent)}},
			stmtPos: 0,
		},
		CatchBlocks: []*Catch{
			{
				ErrorVar: "e",
				ErrorType: &AttributeAccess{
					IdentOrCall: &IdentOrCall{
						Ident: &errorClassName,
					},
				},
				Stmts: &BlockStmts{
					Stmts:   []*Stmt{},
					stmtPos: 0,
				},
			},
		},
	}

	errMessage := "This is an error"
	err, _ := builtin.NewErrorInstance(errMessage)
	state := StateImpl{
		interpreter: &testInterpreter{},
		context: &ContextImpl{
			scopes: []map[string]types.Object{
				{builtin.ErrorClass.Name: builtin.ErrorClass},
				{errorIdent.String(): err},
			},
		},
	}
	result := unsafe.Evaluate(&state, false, false)
	if result.State != StmtNone {
		assertionFailed(t, result.State.String(), StmtNone.String())
	}

	if result.Err != nil {
		t.Error("error is not nil")
	}

	if result.Value == nil {
		t.Error("value is nil")
	}

	if _, ok := result.Value.(types.NilInstance); !ok {
		t.Error("result value is not ErrorInstance")
	}
}

func TestUnsafe_EvaluateThrownRethrownAndNotCaught(t *testing.T) {
	errorClass := types.Class{
		Name:  "CustomError",
		Class: types.TypeClass,
		Bases: []*types.Class{builtin.ErrorClass},
	}
	errorIdent := Ident("Error")
	eIdent := Ident("e")
	errorClassName := Ident(errorClass.Name)
	unsafe := &Unsafe{
		Stmts: &BlockStmts{
			Stmts:   []*Stmt{{Throw: makeThrowStmt(&errorIdent)}},
			stmtPos: 0,
		},
		CatchBlocks: []*Catch{
			{
				ErrorVar: eIdent,
				ErrorType: &AttributeAccess{
					IdentOrCall: &IdentOrCall{
						Ident: &errorClassName,
					},
				},
				Stmts: &BlockStmts{
					Stmts:   []*Stmt{{Throw: makeThrowStmt(&eIdent)}},
					stmtPos: 0,
				},
			},
		},
	}

	errMessage := "This is an error"
	err, _ := builtin.NewErrorInstance(errMessage)
	state := StateImpl{
		interpreter: &testInterpreter{},
		context: &ContextImpl{
			scopes: []map[string]types.Object{
				{errorClass.Name: &errorClass},
				{errorIdent.String(): err},
			},
		},
	}
	result := unsafe.Evaluate(&state, false, false)
	if result.State != StmtThrow {
		assertionFailed(t, result.State.String(), StmtThrow.String())
	}

	if result.Err == nil {
		t.Error("error is nil")
	}

	errMessage = fmt.Sprintf("%s: %s", builtin.ErrorClass.GetName(), errMessage)
	if result.Err.Error() != errMessage {
		assertionFailed(t, result.Err.Error(), errMessage)
	}

	if result.Value == nil {
		t.Error("value is nil")
	}

	if _, ok := result.Value.(*builtin.ErrorInstance); !ok {
		t.Error("result value is not ErrorInstance")
	}

	if result.Value.(*builtin.ErrorInstance) != err {
		t.Error("result value is not expected error")
	}
}

func assertionFailed(t *testing.T, actual, expected string) {
	t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", actual, expected)
}
