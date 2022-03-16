package utilities

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

type RuntimeStatementError struct {
	err       error
	statement types.Statement
}

func (e RuntimeStatementError) Error() string {
	return e.err.Error()
}

func (e RuntimeStatementError) Statement() types.Statement {
	return e.statement
}

func NewRuntimeStatementError(message string, statement types.Statement) RuntimeStatementError {
	return RuntimeStatementError{
		err:       errors.New(message),
		statement: statement,
	}
}
