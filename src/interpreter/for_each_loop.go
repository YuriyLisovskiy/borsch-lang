package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/models"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeForEachLoop(
	indexVar, itemVar models.Token, containerValue types.ValueType,
	body []models.Token, currentFile string,
) (types.ValueType, error) {
	switch container := containerValue.(type) {
	case types.SequentialType:
		var err error
		for idx := int64(0); idx < container.Length(); idx++ {
			scope := map[string]types.ValueType{}
			if indexVar.Text != "_" {
				scope[indexVar.Text] = types.IntegerType{Value: idx}
			}

			if itemVar.Text != "_" {
				scope[itemVar.Text], err = container.GetElement(idx)
				if err != nil {
					return nil, err
				}
			}

			result, err := i.executeBlock(scope, body, currentFile)
			if err != nil {
				return nil, err
			}

			if result != nil {
				return result, nil
			}
		}
	default:
		return nil, util.RuntimeError(fmt.Sprintf(
			"тип '%s' не є об'єктом, по якому можна ітерувати", container.TypeName(),
		))
	}

	return nil, nil
}
