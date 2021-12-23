package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
)

func (i *Interpreter) executeIfSequence(
	blocks []ast.ConditionBlock, elseBlock []models.Token,
	rootDir string, thisPackage, parentPackage string,
) (types.Type, bool, error) {
	for _, block := range blocks {
		result, _, err := i.executeNode(block.Condition, rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, false, err
		}

		result, err = types.ToBool([]types.Type{result}...)
		if err != nil {
			return nil, false, err
		}

		if result.(types.BoolInstance).Value {
			return i.executeBlock(map[string]types.Type{}, block.Tokens, thisPackage, parentPackage)
		}
	}

	if len(elseBlock) > 0 {
		return i.executeBlock(map[string]types.Type{}, elseBlock, thisPackage, parentPackage)
	}

	return types.NilInstance{}, false, nil
}
