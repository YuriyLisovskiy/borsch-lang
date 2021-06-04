package ast

type AST struct {
	CodeNodes []ExpressionNode
}

func NewAST() *AST {
	return &AST{
		[]ExpressionNode{},
	}
}

func (a *AST) AddNode(node ExpressionNode) {
	a.CodeNodes = append(a.CodeNodes, node)
}

func (a AST) String() string {
	return "AST"
}
