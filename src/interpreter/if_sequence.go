package interpreter

import (
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/models"
)

func (i *Interpreter) executeIfSequence(
	blocks []ast.ConditionBlock, elseBlock []models.Token, rootDir string, thisPackage, parentPackage string,
) (types.ValueType, error) {
	for _, block := range blocks {
		result, err := i.executeNode(block.Condition, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		result, err = builtin.ToBool([]types.ValueType{result}...)
		if err != nil {
			return nil, err
		}

		if result.(types.BoolType).Value {
			return i.executeBlock(map[string]types.ValueType{}, block.Tokens, thisPackage, parentPackage)
		}
	}

	if len(elseBlock) > 0 {
		return i.executeBlock(map[string]types.ValueType{}, elseBlock, thisPackage, parentPackage)
	}

	return nil, nil
}
