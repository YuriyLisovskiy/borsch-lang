package interpreter

import (
	"fmt"
	"testing"

	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/internal/common"
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

func TestThrow_EvaluateSuccess(t *testing.T) {
	errorIdent := Ident("err")
	throwNode := makeThrowStmt(&errorIdent)
	errMessage := "This is an error"
	exc := types2.NewError(errMessage)
	state := StateImpl{
		context: &ContextImpl{
			scopes: []map[string]types2.Object{
				{errorIdent.String(): exc},
			},
		},
		stacktrace: &common.StackTrace{},
	}

	result := throwNode.Evaluate(&state)
	if result.State != StmtThrow {
		t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", result.State.String(), StmtThrow.String())
	}

	errMessage = fmt.Sprintf("%s: %s", types2.ErrorClass.Name, errMessage)
	if result.Err.Error() != errMessage {
		t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", result.Err.Error(), errMessage)
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
		context: &ContextImpl{
			scopes: []map[string]types2.Object{
				{errorIdent.String(): types2.String(errMessage)},
			},
		},
		stacktrace: &common.StackTrace{},
	}

	result := throwNode.Evaluate(&state)
	if result.State != StmtNone {
		t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", result.State.String(), StmtNone.String())
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
	unsafe := &Block{
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
	err := types2.NewError(errMessage)
	state := StateImpl{
		context: &ContextImpl{
			scopes: []map[string]types2.Object{
				{errorIdent.String(): err},
			},
		},
		stacktrace: &common.StackTrace{},
	}
	result := unsafe.Evaluate(&state, false, false)
	if result.State == StmtThrow {
		t.Error("state is StmtThrow")
	}

	if result.Err != nil {
		t.Error("error is not nil")
	}

	if _, ok := result.Value.(types2.NilType); !ok {
		t.Error("result value is not NilInstance")
	}
}

func TestUnsafe_EvaluateThrownAndNotCaught(t *testing.T) {
	errorClass := types2.ErrorClass.ClassNew("CustomError", map[string]types2.Object{}, false, nil, nil)
	errorClassName := Ident(errorClass.Name)
	errorIdent := Ident("Error")
	unsafe := &Block{
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
	err := types2.NewError(errMessage)
	state := StateImpl{
		context: &ContextImpl{
			scopes: []map[string]types2.Object{
				{errorClass.Name: errorClass},
				{errorIdent.String(): err},
			},
		},
		stacktrace: &common.StackTrace{},
	}
	result := unsafe.Evaluate(&state, false, false)
	if result.State != StmtThrow {
		t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", result.State.String(), StmtThrow.String())
	}

	if result.Err == nil {
		t.Error("error is nil")
	}

	errMessage = fmt.Sprintf("%s: %s", types2.ErrorClass.Name, errMessage)
	if result.Err.Error() != errMessage {
		t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", result.Err.Error(), errMessage)
	}

	if result.Value == nil {
		t.Error("value is nil")
	}

	if _, ok := result.Value.(*types2.Error); !ok {
		t.Error("result value is not ErrorInstance")
	}

	if result.Value.(*types2.Error) != err {
		t.Error("result value is not expected error")
	}
}

func TestUnsafe_EvaluateThrownAndCaught(t *testing.T) {
	errorIdent := Ident("Error")
	errorClassName := Ident(types2.ErrorClass.Name)
	unsafe := &Block{
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
	err := types2.NewError(errMessage)
	state := StateImpl{
		context: &ContextImpl{
			scopes: []map[string]types2.Object{
				{types2.ErrorClass.Name: types2.ErrorClass},
				{errorIdent.String(): err},
			},
		},
		stacktrace: &common.StackTrace{},
	}
	result := unsafe.Evaluate(&state, false, false)
	if result.State != StmtNone {
		t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", result.State.String(), StmtNone.String())
	}

	if result.Err != nil {
		t.Error("error is not nil")
	}

	if result.Value == nil {
		t.Error("value is nil")
	}

	if _, ok := result.Value.(types2.NilType); !ok {
		t.Error("result value is not ErrorInstance")
	}
}

func TestUnsafe_EvaluateThrownRethrownAndNotCaught(t *testing.T) {
	errorClass := types2.ErrorClass.ClassNew("CustomError", map[string]types2.Object{}, false, nil, nil)
	errorIdent := Ident("Error")
	eIdent := Ident("e")
	errorClassName := Ident(errorClass.Name)
	unsafe := &Block{
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
	err := types2.NewError(errMessage)
	state := StateImpl{
		context: &ContextImpl{
			scopes: []map[string]types2.Object{
				{errorClass.Name: errorClass},
				{errorIdent.String(): err},
			},
		},
		stacktrace: &common.StackTrace{},
	}
	result := unsafe.Evaluate(&state, false, false)
	if result.State != StmtThrow {
		t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", result.State.String(), StmtThrow.String())
	}

	if result.Err == nil {
		t.Error("error is nil")
	}

	errMessage = fmt.Sprintf("%s: %s", types2.ErrorClass.Name, errMessage)
	if result.Err.Error() != errMessage {
		t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", result.Err.Error(), errMessage)
	}

	if result.Value == nil {
		t.Error("value is nil")
	}

	if _, ok := result.Value.(*types2.Error); !ok {
		t.Error("result value is not ErrorInstance")
	}

	if result.Value.(*types2.Error) != err {
		t.Error("result value is not expected error")
	}
}
