package ast

type ExpressionNode interface {
	String() string
	RowNumber() int
}
