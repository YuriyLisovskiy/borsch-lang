package interpreter

import (
	"testing"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func makeThrowStmt(name *string) *Throw {
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
																	SlicingOrSubscription: &SlicingOrSubscription{
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
	errorIdent := "err"
	throwNode := makeThrowStmt(&errorIdent)
	errMessage := "This is an error"
	exc, _ := builtin.NewErrorInstance(errMessage)
	state := StateImpl{
		context: &ContextImpl{
			scopes: []map[string]common.Value{
				{errorIdent: exc},
			},
		},
	}

	result := throwNode.Evaluate(&state)
	if result.State != StmtThrow {
		assertionFailed(t, result.State.String(), StmtThrow.String())
	}

	if result.Err.Error() != errMessage {
		assertionFailed(t, result.Err.Error(), errMessage)
	}

	if result.Value != exc {
		t.Error("value is not exception")
	}
}

func TestThrow_EvaluateFail_NotAnErrorInstance(t *testing.T) {
	errorIdent := "err"
	throwNode := makeThrowStmt(&errorIdent)
	errMessage := "This is an error"
	state := StateImpl{
		context: &ContextImpl{
			scopes: []map[string]common.Value{
				{errorIdent: types.NewStringInstance(errMessage)},
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

func assertionFailed(t *testing.T, actual, expected string) {
	t.Errorf("Assertion failed:\nActual:\n%s\n\nExpected:\n%s", actual, expected)
}
