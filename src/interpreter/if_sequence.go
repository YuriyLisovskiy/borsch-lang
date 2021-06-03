package interpreter

import (
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/models"
)

func (e *Interpreter) executeIfSequence(
	blocks []ast.ConditionBlock, elseBlock []models.Token, rootDir string, currentFile string,
) (builtin.ValueType, error) {
	for _, block := range blocks {
		result, err := e.executeNode(block.Condition, rootDir, currentFile)
		if err != nil {
			return nil, err
		}

		result, err = builtin.CastToBool([]builtin.ValueType{result}...)
		if err != nil {
			return nil, err
		}

		if result.(builtin.BoolType).Value {
			return e.executeBlock(block.Tokens, currentFile)
		}
	}

	if len(elseBlock) > 0 {
		return e.executeBlock(elseBlock, currentFile)
	}

	return builtin.NoneType{}, nil
}
