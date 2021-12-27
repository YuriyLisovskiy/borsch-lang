package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeListSlicing(ctx *Context, node *ast.ListSlicingNode) (types.Type, error) {
	container, _, err := i.executeNode(ctx, node.Operand)
	if err != nil {
		return nil, err
	}

	if container.GetTypeHash() == types.ListTypeHash {
		fromIdx, _, err := i.executeNode(ctx, node.LeftIndex)
		if err != nil {
			return nil, err
		}

		if fromIdx.GetTypeHash() == types.IntegerTypeHash {
			toIdx, _, err := i.executeNode(ctx, node.RightIndex)
			if err != nil {
				return nil, err
			}

			if toIdx.GetTypeHash() == types.IntegerTypeHash {
				res, err := container.(types.ListInstance).Slice(
					fromIdx.(types.IntegerInstance).Value, toIdx.(types.IntegerInstance).Value,
				)
				return res, err
			}

			return nil, util.RuntimeError("правий індекс має бути цілого типу")
		}

		return nil, util.RuntimeError("лівий індекс має бути цілого типу")
	}

	return nil, util.RuntimeError(fmt.Sprintf(
		"неможливо застосувати оператор відсікання списку до об'єкта з типом '%s'",
		container.GetTypeName(),
	))
}
