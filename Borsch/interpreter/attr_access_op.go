package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

func (i *Interpreter) executeAttrAccessOp(
	node *ast.AttrAccessOpNode, rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	val, _, err := i.executeNode(node.Base, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	switch attr := node.Attr.(type) {
	case ast.VariableNode:
		val, err = val.GetAttr(attr.Variable.Text)
		if err != nil {
			return nil, err
		}

		return val, nil
	case ast.CallOpNode:
		val, err = val.GetAttr(attr.CallableName.Text)
		if err != nil {
			return nil, err
		}

		val, err = i.executeCallOp(&attr, val, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		return val, nil
	default:
		panic("fatal: invalid node")
	}
}
