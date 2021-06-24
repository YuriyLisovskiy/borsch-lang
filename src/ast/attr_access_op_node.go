package ast

type AttrAccessOpNode struct {
	Base ExpressionNode
	Attr ExpressionNode

	rowNumber int
}

func NewAttrAccessOpNode(expression ExpressionNode, attr ExpressionNode, rowNumber int) AttrAccessOpNode {
	return AttrAccessOpNode{
		Base: expression,
		Attr:       attr,
		rowNumber:  rowNumber,
	}
}

func (n AttrAccessOpNode) String() string {
	return n.Base.String() + "." + n.Attr.String()
}

func (n AttrAccessOpNode) RowNumber() int {
	return n.rowNumber
}
