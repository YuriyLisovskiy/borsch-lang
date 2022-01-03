package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
)

func (i *Interpreter) executeIfSequence(
	ctx *Context,
	blocks []ast.ConditionBlock,
	elseBlock []models.Token,
) (types.Type, bool, error) {
	for _, block := range blocks {
		result, _, err := i.executeNode(ctx, block.Condition)
		if err != nil {
			return nil, false, err
		}

		result, err = types.ToBool([]types.Type{result}...)
		if err != nil {
			return nil, false, err
		}

		if result.(types.BoolInstance).Value {
			return i.executeBlock(ctx, map[string]types.Type{}, block.Tokens)
		}
	}

	if len(elseBlock) > 0 {
		return i.executeBlock(ctx, map[string]types.Type{}, elseBlock)
	}

	return nil, false, nil
}
