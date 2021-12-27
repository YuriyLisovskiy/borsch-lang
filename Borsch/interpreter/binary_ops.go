package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeBinaryOp(ctx *Context, node *ast.BinOperationNode) (types.Type, error) {
	switch node.Operator.Type.Name {
	case models.ExponentOp:
		return i.executeArithmeticOp(ctx, node.LeftNode, node.RightNode, ops.PowOp)
	case models.ModuloOp:
		return i.executeArithmeticOp(ctx, node.LeftNode, node.RightNode, ops.ModuloOp)
	case models.Add:
		return i.executeArithmeticOp(ctx, node.LeftNode, node.RightNode, ops.AddOp)
	case models.Sub:
		return i.executeArithmeticOp(ctx, node.LeftNode, node.RightNode, ops.SubOp)
	case models.Mul:
		return i.executeArithmeticOp(ctx, node.LeftNode, node.RightNode, ops.MulOp)
	case models.Div:
		return i.executeArithmeticOp(ctx, node.LeftNode, node.RightNode, ops.DivOp)
	case models.AndOp:
		return i.executeLogicalOp(ctx, node.LeftNode, node.RightNode, ops.AndOp)
	case models.OrOp:
		return i.executeLogicalOp(ctx, node.LeftNode, node.RightNode, ops.OrOp)
	case models.BitwiseLeftShiftOp:
		return i.executeBitwiseOp(ctx, node.LeftNode, node.RightNode, ops.BitwiseLeftShiftOp)
	case models.BitwiseRightShiftOp:
		return i.executeBitwiseOp(ctx, node.LeftNode, node.RightNode, ops.BitwiseRightShiftOp)
	case models.BitwiseAndOp:
		return i.executeBitwiseOp(ctx, node.LeftNode, node.RightNode, ops.BitwiseAndOp)
	case models.BitwiseXorOp:
		return i.executeBitwiseOp(ctx, node.LeftNode, node.RightNode, ops.BitwiseXorOp)
	case models.BitwiseOrOp:
		return i.executeBitwiseOp(ctx, node.LeftNode, node.RightNode, ops.BitwiseOrOp)
	case models.EqualsOp:
		return i.executeComparisonOp(ctx, node.LeftNode, node.RightNode, ops.EqualsOp)
	case models.NotEqualsOp:
		return i.executeComparisonOp(ctx, node.LeftNode, node.RightNode, ops.NotEqualsOp)
	case models.GreaterOp:
		return i.executeComparisonOp(ctx, node.LeftNode, node.RightNode, ops.GreaterOp)
	case models.GreaterOrEqualsOp:
		return i.executeComparisonOp(ctx, node.LeftNode, node.RightNode, ops.GreaterOrEqualsOp)
	case models.LessOp:
		return i.executeComparisonOp(ctx, node.LeftNode, node.RightNode, ops.LessOp)
	case models.LessOrEqualsOp:
		return i.executeComparisonOp(ctx, node.LeftNode, node.RightNode, ops.LessOrEqualsOp)
	case models.Assign:
		rightNode, _, err := i.executeNode(ctx, node.RightNode)
		if err != nil {
			return nil, err
		}

		switch leftNode := node.LeftNode.(type) {
		case ast.VariableNode:
			return rightNode, i.setVar(ctx.package_.Name, leftNode.Variable.Text, rightNode)
		case ast.CallOpNode:
			return nil, util.RuntimeError("неможливо присвоїти значення виклику функції")
		case ast.RandomAccessOperationNode:
			variable, _, err := i.executeNode(ctx, leftNode.Operand)
			if err != nil {
				return nil, err
			}

			variable, err = i.executeRandomAccessSetOp(ctx, leftNode.Index, variable, rightNode)
			if err != nil {
				return nil, err
			}

			operand := leftNode.Operand
			for {
				switch external := operand.(type) {
				case ast.RandomAccessOperationNode:
					opVar, _, err := i.executeNode(ctx, external.Operand)
					if err != nil {
						return nil, err
					}

					variable, err = i.executeRandomAccessSetOp(ctx, external.Index, opVar, variable)
					if err != nil {
						return nil, err
					}

					operand = external.Operand
					continue
				case ast.VariableNode:
					err = i.setVar(ctx.package_.Name, external.Variable.Text, variable)
				}

				break
			}

			return variable, nil
		case ast.AttrAccessOpNode:
			base, _, err := i.executeNode(ctx, leftNode.Base)
			if err != nil {
				return nil, err
			}

			switch attrNode := leftNode.Attr.(type) {
			case ast.VariableNode:
				base, err = base.SetAttribute(attrNode.Variable.Text, rightNode)
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
