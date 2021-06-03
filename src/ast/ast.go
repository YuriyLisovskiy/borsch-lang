package ast

type AST struct {
	CodeRows []ExpressionNode
}

func NewAST() *AST {
	return &AST{
		[]ExpressionNode{},
	}
}

func (a *AST) AddNode(node ExpressionNode) {
	a.CodeRows = append(a.CodeRows, node)
}

func (a AST) String() string {
	return "AST"
}
