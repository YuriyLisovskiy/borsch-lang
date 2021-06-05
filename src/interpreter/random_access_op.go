package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"unicode/utf8"
)

func (i *Interpreter) executeRandomAccessOp(
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
		case builtin.StringType:
			runesCount := int64(utf8.RuneCountInString(target.Value))
			if index.Value >= 0 && index.Value < runesCount {
				runes := []rune(target.Value)
				return builtin.StringType{Value: string(runes[index.Value])}, nil
			} else if index.Value < 0 && index.Value >= -runesCount {
				runes := []rune(target.Value)
				return builtin.StringType{Value: string(runes[runesCount + index.Value])}, nil
			} else {
				return nil, util.RuntimeError("індекс рядка за межами послідовності")
			}
		default:
			return nil, util.RuntimeError(fmt.Sprintf(
				"об'єкт з типом '%s' не підтримує індексування", target.TypeName(),
			))
		}
	default:
		return nil, util.RuntimeError("індекси мають бути цілого типу")
	}
}
