package interpreter

import (
	"github.com/YuriyLisovskiy/borsch/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch/Borsch/models"
)

func (i *Interpreter) executeIfSequence(
	blocks []ast.ConditionBlock, elseBlock []models.Token,
	rootDir string, thisPackage, parentPackage string,
) (types.ValueType, bool, error) {
	for _, block := range blocks {
		result, _, err := i.executeNode(block.Condition, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, false, err
		}

		result, err = builtin.ToBool([]types.ValueType{result}...)
		if err != nil {
			return nil, false, err
		}

		if result.(types.BoolType).Value {
			return i.executeBlock(map[string]types.ValueType{}, block.Tokens, thisPackage, parentPackage)
		}
	}

	if len(elseBlock) > 0 {
		return i.executeBlock(map[string]types.ValueType{}, elseBlock, thisPackage, parentPackage)
	}

	return types.NilType{}, false, nil
}
