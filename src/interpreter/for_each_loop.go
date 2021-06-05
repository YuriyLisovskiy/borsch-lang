package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin"
	"github.com/YuriyLisovskiy/borsch/src/models"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeForEachLoop(
	indexVar, itemVar models.Token, containerValue builtin.ValueType,
	body []models.Token, currentFile string,
) (builtin.ValueType, error) {
	switch container := containerValue.(type) {
	case builtin.StringType:
		runes := []rune(container.Value)
		for idx, obj := range runes {
			scope := map[string]builtin.ValueType{}
			if indexVar.Text != "_" {
				scope[indexVar.Text] = builtin.IntegerNumberType{Value: int64(idx)}
			}

			if itemVar.Text != "_" {
				scope[itemVar.Text] = builtin.StringType{Value: string(obj)}
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
