package utilities

import (
	"errors"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type RuntimeStatementError struct {
	err       error
	statement common.Statement
}

func (e RuntimeStatementError) Error() string {
	return e.err.Error()
}

func (e RuntimeStatementError) Statement() common.Statement {
	return e.statement
}

func NewRuntimeStatementError(message string, statement common.Statement) RuntimeStatementError {
	return RuntimeStatementError{
		err:       errors.New(message),
		statement: statement,
	}
}
