package ast
//
//import (
//	"github.com/YuriyLisovskiy/borsch/src/models"
//)
//
//type RandomAccessSetOperationNode struct {
//	Variable models.Token
//	Index    ExpressionNode
//
//	rowNumber int
//}
//
//func NewRandomAccessSetOperationNode(
//	variable models.Token, index ExpressionNode, rowNumber int,
//) RandomAccessSetOperationNode {
//	return RandomAccessSetOperationNode{
//		Variable:  variable,
//		Index:     index,
//		rowNumber: rowNumber,
//	}
//}
//
//func (n RandomAccessSetOperationNode) String() string {
//	return n.Variable.String() + "[" + n.Index.String() + "]"
//}
//
//func (n RandomAccessSetOperationNode) RowNumber() int {
//	return n.rowNumber
//}
