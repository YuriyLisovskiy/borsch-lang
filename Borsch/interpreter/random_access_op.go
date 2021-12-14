package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeRandomAccessGetOp(
	targetNode, indexNode ast.ExpressionNode, rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	targetVal, _, err := i.executeNode(targetNode, rootDir, thisPackage, parentPackage)
	if err != nil {
		return nil, err
	}

	switch target := targetVal.(type) {
	case types.SequentialType:
		indexVal, _, err := i.executeNode(indexNode, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		switch index := indexVal.(type) {
		case types.IntegerType:
			elem, err := target.GetElement(index.Value)
			if err != nil {
				return nil, util.RuntimeError(err.Error())
			}

			return elem, nil
		default:
			return nil, util.RuntimeError("індекси мають бути цілого типу")
		}
	case types.DictionaryType:
		key, _, err := i.executeNode(indexNode, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		elem, err := target.GetElement(key)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}

		return elem, nil
	default:
		return nil, util.RuntimeError(fmt.Sprintf(
			"неможливо застосувати оператор довільного доступу до об'єкта з типом '%s'",
			target.TypeName(),
		))
	}
}

func (i *Interpreter) executeRandomAccessSetOp(
	indexNode ast.ExpressionNode, variable, value types.ValueType,
	rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	switch container := variable.(type) {
	case types.SequentialType:
		indexVal, _, err := i.executeNode(indexNode, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		switch index := indexVal.(type) {
		case types.IntegerType:
			newIterable, err := container.SetElement(index.Value, value)
			if err != nil {
				return nil, util.RuntimeError(err.Error())
			}

			return newIterable, nil
		default:
			return nil, util.RuntimeError("індекси мають бути цілого типу")
		}
	case types.DictionaryType:
		key, _, err := i.executeNode(indexNode, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		err = container.SetElement(key, value)
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}

		return container, nil
	default:
		return nil, util.RuntimeError(fmt.Sprintf(
			"неможливо застосувати оператор довільного доступу до об'єкта з типом '%s'",
			container.TypeName(),
		))
	}
}
