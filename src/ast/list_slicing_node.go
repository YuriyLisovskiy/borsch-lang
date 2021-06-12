package ast

type ListSlicingNode struct {
	Operand    ExpressionNode
	LeftIndex  ExpressionNode
	RightIndex ExpressionNode

	rowNumber int
}

func NewListSlicingNode(
	operand, leftIndex, rightIndex ExpressionNode, rowNumber int,
) ListSlicingNode {
	return ListSlicingNode{
		Operand:    operand,
		LeftIndex:  leftIndex,
		RightIndex: rightIndex,
		rowNumber:  rowNumber,
	}
}

func (n ListSlicingNode) String() string {
	return n.Operand.String() + "[" + n.LeftIndex.String() + ":" + n.RightIndex.String() + "]"
}

func (n ListSlicingNode) RowNumber() int {
	return n.rowNumber
}
