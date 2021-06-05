package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeRandomAccessGetOp(
	targetNode, indexNode ast.ExpressionNode, rootDir string, currentFile string,
) (builtin.ValueType, error) {
	indexVal, err := i.executeNode(indexNode, rootDir, currentFile)
	if err != nil {
		return nil, err
	}

	switch index := indexVal.(type) {
	case builtin.IntegerNumberType:
		targetVal, err := i.executeNode(targetNode, rootDir, currentFile)
		if err != nil {
			return nil, err
		}

		switch target := targetVal.(type) {
		case builtin.IterableType:
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
	indexNode ast.ExpressionNode, variable, value builtin.ValueType,
	rootDir string, currentFile string,
) (builtin.ValueType, error) {
	switch iterable := variable.(type) {
	case builtin.IterableType:
		indexVal, err := i.executeNode(indexNode, rootDir, currentFile)
		if err != nil {
			return nil, err
		}

		switch index := indexVal.(type) {
		case builtin.IntegerNumberType:
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
