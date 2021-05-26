package ast

type AST struct {
	CodeRows []ExpressionNode
}

func NewStatementsNode() AST {
	return AST{
		[]ExpressionNode{},
	}
}

func (sn *AST) AddNode(node ExpressionNode) {
	sn.CodeRows = append(sn.CodeRows, node)
}

func (sn AST) String() string {
	return "AST"
}
