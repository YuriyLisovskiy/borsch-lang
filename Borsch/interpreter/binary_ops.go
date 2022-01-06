package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeBinaryOp(ctx *Context, node *ast.BinOperationNode) (types.Type, error) {
	if node.Operator.Type.Name == models.Assign {
		rightNode, _, err := i.executeNode(ctx, node.RightNode)
		if err != nil {
			return nil, err
		}

		switch leftNode := node.LeftNode.(type) {
		case ast.VariableNode:
			return nil, i.setVar(ctx.GetPackageFromParent(), leftNode.Variable.Text, rightNode)
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
					err = i.setVar(ctx.GetPackageFromParent(), external.Variable.Text, variable)
				}

				break
			}

			return nil, nil
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

			return nil, nil
		default:
			panic("fatal: invalid node")
		}
	}

	left, _, err := i.executeNode(ctx, node.LeftNode)
	if err != nil {
		return nil, err
	}

	right, _, err := i.executeNode(ctx, node.RightNode)
	if err != nil {
		return nil, err
	}

	switch node.Operator.Type.Name {
	case models.ExponentOp:
		return i.executeArithmeticOp(left, right, ops.PowOp)
	case models.ModuloOp:
		return i.executeArithmeticOp(left, right, ops.ModuloOp)
	case models.Add:
		return i.executeArithmeticOp(left, right, ops.AddOp)
	case models.Sub:
		return i.executeArithmeticOp(left, right, ops.SubOp)
	case models.Mul:
		return i.executeArithmeticOp(left, right, ops.MulOp)
	case models.Div:
		return i.executeArithmeticOp(left, right, ops.DivOp)
	case models.AndOp:
		return i.executeLogicalOp(left, right, ops.AndOp)
	case models.OrOp:
		return i.executeLogicalOp(left, right, ops.OrOp)
	case models.BitwiseLeftShiftOp:
		return i.executeBitwiseOp(left, right, ops.BitwiseLeftShiftOp)
	case models.BitwiseRightShiftOp:
		return i.executeBitwiseOp(left, right, ops.BitwiseRightShiftOp)
	case models.BitwiseAndOp:
		return i.executeBitwiseOp(left, right, ops.BitwiseAndOp)
	case models.BitwiseXorOp:
		return i.executeBitwiseOp(left, right, ops.BitwiseXorOp)
	case models.BitwiseOrOp:
		return i.executeBitwiseOp(left, right, ops.BitwiseOrOp)
	case models.EqualsOp:
		return i.executeComparisonOp(left, right, ops.EqualsOp)
	case models.NotEqualsOp:
		return i.executeComparisonOp(left, right, ops.NotEqualsOp)
	case models.GreaterOp:
		return i.executeComparisonOp(left, right, ops.GreaterOp)
	case models.GreaterOrEqualsOp:
		return i.executeComparisonOp(left, right, ops.GreaterOrEqualsOp)
	case models.LessOp:
		return i.executeComparisonOp(left, right, ops.LessOp)
	case models.LessOrEqualsOp:
		return i.executeComparisonOp(left, right, ops.LessOrEqualsOp)
	default:
		panic("fatal: invalid binary operator")
	}
}
