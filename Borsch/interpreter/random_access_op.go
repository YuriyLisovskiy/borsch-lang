package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeRandomAccessGetOp(ctx *Context, targetNode, indexNode ast.ExpressionNode) (
	types.Type,
	error,
) {
	targetVal, _, err := i.executeNode(ctx, targetNode)
	if err != nil {
		return nil, err
	}

	switch target := targetVal.(type) {
	case types.SequentialType:
		indexVal, _, err := i.executeNode(ctx, indexNode)
		if err != nil {
			return nil, err
		}

		switch index := indexVal.(type) {
		case types.IntegerInstance:
			elem, err := target.GetElement(index.Value)
			if err != nil {
				return nil, util.RuntimeError(err.Error())
			}

			return elem, nil
		default:
			return nil, util.RuntimeError("індекси мають бути цілого типу")
		}
	case types.DictionaryInstance:
		key, _, err := i.executeNode(ctx, indexNode)
		if err != nil {
			return nil, err
		}

		elem, err := target.GetElement(key)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}

		return elem, nil
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор довільного доступу до об'єкта з типом '%s'",
				target.GetTypeName(),
			),
		)
	}
}

func (i *Interpreter) executeRandomAccessSetOp(
	ctx *Context,
	indexNode ast.ExpressionNode,
	variable types.Type,
	value types.Type,
) (types.Type, error) {
	switch container := variable.(type) {
	case types.SequentialType:
		indexVal, _, err := i.executeNode(ctx, indexNode)
		if err != nil {
			return nil, err
		}

		switch index := indexVal.(type) {
		case types.IntegerInstance:
			newIterable, err := container.SetElement(index.Value, value)
			if err != nil {
				return nil, util.RuntimeError(err.Error())
			}

			return newIterable, nil
		default:
			return nil, util.RuntimeError("індекси мають бути цілого типу")
		}
	case types.DictionaryInstance:
		key, _, err := i.executeNode(ctx, indexNode)
		if err != nil {
			return nil, err
		}

		err = container.SetElement(key, value)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}

		return container, nil
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор довільного доступу до об'єкта з типом '%s'",
				container.GetTypeName(),
			),
		)
	}
}
