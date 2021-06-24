package ast

type AttrOpNode struct {
	Base ExpressionNode
	Attr ExpressionNode

	rowNumber int
}

func NewGetAttrOpNode(expression ExpressionNode, attr ExpressionNode, rowNumber int) AttrOpNode {
	return AttrOpNode{
		Base: expression,
		Attr:       attr,
		rowNumber:  rowNumber,
	}
}

func (n AttrOpNode) String() string {
	return n.Base.String() + "." + n.Attr.String()
}

func (n AttrOpNode) RowNumber() int {
	return n.rowNumber
}
