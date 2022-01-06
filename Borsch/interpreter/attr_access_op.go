package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

func (i *Interpreter) executeAttrAccessOp(ctx *Context, node *ast.AttrAccessOpNode) (types.Type, error) {
	val, _, err := i.executeNode(ctx, node.Base)
	if err != nil {
		return nil, err
	}

	switch attr := node.Attr.(type) {
	case ast.VariableNode:
		val, err = val.GetAttribute(attr.Variable.Text)
		if err != nil {
			return nil, err
		}

		return val, nil
	case ast.CallOpNode:
		res, err := val.GetAttribute(attr.CallableName.Text)
		if err != nil {
			return nil, err
		}
		
		res, err = i.executeCallOp(ctx.WithParent(val), &attr, res)
		if err != nil {
			return nil, err
		}

		return res, nil
	default:
		panic("fatal: invalid node")
	}
}
