package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeRandomAccessGetOp(
	targetNode, indexNode ast.ExpressionNode, rootDir string, currentFile string,
) (types.ValueType, error) {
	indexVal, err := i.executeNode(indexNode, rootDir, currentFile)
	if err != nil {
		return nil, err
	}

	switch index := indexVal.(type) {
	case types.IntegerType:
		targetVal, err := i.executeNode(targetNode, rootDir, currentFile)
		if err != nil {
			return nil, err
		}

		switch target := targetVal.(type) {
		case types.SequentialType:
			elem, err := target.GetElement(index.Value)
			if err != nil {
				return nil, util.RuntimeError(err.Error())
			}

			return elem, nil
		default:
			return nil, util.RuntimeError(fmt.Sprintf(
				"об'єкт з типом '%s' не підтримує індексування", target.TypeName(),
			))
		}
	default:
		return nil, util.RuntimeError("індекси мають бути цілого типу")
	}
}

func (i *Interpreter) executeRandomAccessSetOp(
	indexNode ast.ExpressionNode, variable, value types.ValueType,
	rootDir string, currentFile string,
) (types.ValueType, error) {
	switch iterable := variable.(type) {
	case types.SequentialType:
		indexVal, err := i.executeNode(indexNode, rootDir, currentFile)
		if err != nil {
			return nil, err
		}

		switch index := indexVal.(type) {
		case types.IntegerType:
			newIterable, err := iterable.SetElement(index.Value, value)
			if err != nil {
				return nil, util.RuntimeError(err.Error())
			}

			return newIterable, nil
		default:
			return nil, util.RuntimeError("індекси мають бути цілого типу")
		}
	default:
		return nil, util.RuntimeError(fmt.Sprintf(
			"об'єкт з типом '%s' не підтримує індексування", iterable.TypeName(),
		))
	}
}
