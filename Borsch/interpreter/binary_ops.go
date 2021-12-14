package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeBinaryOp(
	node *ast.BinOperationNode, rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	switch node.Operator.Type.Name {
	case models.ExponentOp:
		return i.executeArithmeticOp(node.LeftNode, node.RightNode, ops.PowOp, rootDir, thisPackage, parentPackage)
	case models.ModuloOp:
		return i.executeArithmeticOp(node.LeftNode, node.RightNode, ops.ModuloOp, rootDir, thisPackage, parentPackage)
	case models.Add:
		return i.executeArithmeticOp(node.LeftNode, node.RightNode, ops.AddOp, rootDir, thisPackage, parentPackage)
	case models.Sub:
		return i.executeArithmeticOp(node.LeftNode, node.RightNode, ops.SubOp, rootDir, thisPackage, parentPackage)
	case models.Mul:
		return i.executeArithmeticOp(node.LeftNode, node.RightNode, ops.MulOp, rootDir, thisPackage, parentPackage)
	case models.Div:
		return i.executeArithmeticOp(node.LeftNode, node.RightNode, ops.DivOp, rootDir, thisPackage, parentPackage)
	case models.AndOp:
		return i.executeLogicalOp(node.LeftNode, node.RightNode, ops.AndOp, rootDir, thisPackage, parentPackage)
	case models.OrOp:
		return i.executeLogicalOp(node.LeftNode, node.RightNode, ops.OrOp, rootDir, thisPackage, parentPackage)
	case models.BitwiseLeftShiftOp:
		return i.executeBitwiseOp(node.LeftNode, node.RightNode, ops.BitwiseLeftShiftOp, rootDir, thisPackage, parentPackage)
	case models.BitwiseRightShiftOp:
		return i.executeBitwiseOp(node.LeftNode, node.RightNode, ops.BitwiseRightShiftOp, rootDir, thisPackage, parentPackage)
	case models.BitwiseAndOp:
		return i.executeBitwiseOp(node.LeftNode, node.RightNode, ops.BitwiseAndOp, rootDir, thisPackage, parentPackage)
	case models.BitwiseXorOp:
		return i.executeBitwiseOp(node.LeftNode, node.RightNode, ops.BitwiseXorOp, rootDir, thisPackage, parentPackage)
	case models.BitwiseOrOp:
		return i.executeBitwiseOp(node.LeftNode, node.RightNode, ops.BitwiseOrOp, rootDir, thisPackage, parentPackage)
	case models.EqualsOp:
		return i.executeComparisonOp(node.LeftNode, node.RightNode, ops.EqualsOp, rootDir, thisPackage, parentPackage)
	case models.NotEqualsOp:
		return i.executeComparisonOp(node.LeftNode, node.RightNode, ops.NotEqualsOp, rootDir, thisPackage, parentPackage)
	case models.GreaterOp:
		return i.executeComparisonOp(node.LeftNode, node.RightNode, ops.GreaterOp, rootDir, thisPackage, parentPackage)
	case models.GreaterOrEqualsOp:
		return i.executeComparisonOp(node.LeftNode, node.RightNode, ops.GreaterOrEqualsOp, rootDir, thisPackage, parentPackage)
	case models.LessOp:
		return i.executeComparisonOp(node.LeftNode, node.RightNode, ops.LessOp, rootDir, thisPackage, parentPackage)
	case models.LessOrEqualsOp:
		return i.executeComparisonOp(node.LeftNode, node.RightNode, ops.LessOrEqualsOp, rootDir, thisPackage, parentPackage)
	case models.Assign:
		rightNode, _, err := i.executeNode(node.RightNode, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		switch leftNode := node.LeftNode.(type) {
		case ast.VariableNode:
			return rightNode, i.setVar(thisPackage, leftNode.Variable.Text, rightNode)
		case ast.CallOpNode:
			return nil, util.RuntimeError("неможливо присвоїти значення виклику функції")
		case ast.RandomAccessOperationNode:
			variable, _, err := i.executeNode(leftNode.Operand, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			variable, err = i.executeRandomAccessSetOp(
				leftNode.Index, variable, rightNode, rootDir, thisPackage, parentPackage,
			)
			if err != nil {
				return nil, err
			}

			operand := leftNode.Operand
			for {
				switch external := operand.(type) {
				case ast.RandomAccessOperationNode:
					opVar, _, err := i.executeNode(external.Operand, rootDir, thisPackage, parentPackage)
					if err != nil {
						return nil, err
					}

					variable, err = i.executeRandomAccessSetOp(
						external.Index, opVar, variable, rootDir, thisPackage, parentPackage,
					)
					if err != nil {
						return nil, err
					}

					operand = external.Operand
					continue
				case ast.VariableNode:
					err = i.setVar(thisPackage, external.Variable.Text, variable)
				}

				break
			}

			return variable, nil
		case ast.AttrAccessOpNode:
			base, _, err := i.executeNode(leftNode.Base, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			switch attrNode := leftNode.Attr.(type) {
			case ast.VariableNode:
				base, err = base.SetAttr(attrNode.Variable.Text, rightNode)
				if err != nil {
					return nil, err
				}
			case ast.CallOpNode:
				return nil, util.RuntimeError("неможливо присвоїти значення виклику функції")
			default:
				panic("fatal: invalid node")
			}

			return base, nil
		default:
			panic("fatal: invalid node")
		}
	default:
		panic("fatal: invalid binary operator")
	}
}
