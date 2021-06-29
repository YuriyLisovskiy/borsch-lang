package ast

type RandomAccessOperationNode struct {
	Operand ExpressionNode
	Index   ExpressionNode
	IsSet   bool

	rowNumber int
}

func NewRandomAccessOperationNode(operand, index ExpressionNode, rowNumber int) RandomAccessOperationNode {
	return RandomAccessOperationNode{
		Operand:   operand,
		Index:     index,
		rowNumber: rowNumber,
	}
}

func (n RandomAccessOperationNode) String() string {
	return n.Operand.String() + "[" + n.Index.String() + "]"
}

func (n RandomAccessOperationNode) RowNumber() int {
	return n.rowNumber
}
