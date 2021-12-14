package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeListSlicing(
	node *ast.ListSlicingNode, rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	container, _, err := i.executeNode(node.Operand, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	if container.TypeHash() == types.ListTypeHash {
		fromIdx, _, err := i.executeNode(node.LeftIndex, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		if fromIdx.TypeHash() == types.IntegerTypeHash {
			toIdx, _, err := i.executeNode(node.RightIndex, rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			if toIdx.TypeHash() == types.IntegerTypeHash {
				res, err := container.(types.ListType).Slice(
					fromIdx.(types.IntegerType).Value, toIdx.(types.IntegerType).Value,
				)
				return res, err
			}

			return nil, util.RuntimeError("правий індекс має бути цілого типу")
		}

		return nil, util.RuntimeError("лівий індекс має бути цілого типу")
	}

	return nil, util.RuntimeError(fmt.Sprintf(
		"неможливо застосувати оператор відсікання списку до об'єкта з типом '%s'",
		container.TypeName(),
	))
}
