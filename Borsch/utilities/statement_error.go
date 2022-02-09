package utilities

import (
	"errors"

	"github.com/alecthomas/participle/v2/lexer"
)

type ErrorStatement interface {
	String() string
	Position() lexer.Position
}

type RuntimeStatementError struct {
	err       error
	statement ErrorStatement
}

func (e RuntimeStatementError) Error() string {
	return e.err.Error()
}

func (e RuntimeStatementError) Statement() ErrorStatement {
	return e.statement
}

func NewRuntimeStatementError(message string, statement ErrorStatement) RuntimeStatementError {
	return RuntimeStatementError{
		err:       errors.New(message),
		statement: statement,
	}
}
