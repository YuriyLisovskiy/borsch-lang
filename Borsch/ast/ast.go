package ast

type AST struct {
	Nodes []ExpressionNode
}

func NewAST() *AST {
	return &AST{
		[]ExpressionNode{},
	}
}

func (a *AST) AddNode(node ExpressionNode) {
	a.Nodes = append(a.Nodes, node)
}

func (a AST) String() string {
	return "AST"
}
