package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/models"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) executeForEachLoop(
	indexVar, itemVar models.Token, containerValue types.Type,
	body []models.Token, thisPackage, parentPackage string,
) (types.Type, bool, error) {
	switch container := containerValue.(type) {
	case types.SequentialType:
		var err error
		sz := container.Length()
		for idx := int64(0); idx < sz; idx++ {
			scope := map[string]types.Type{}
			if indexVar.Text != "_" {
				scope[indexVar.Text] = types.IntegerType{Value: idx}
			}

			if itemVar.Text != "_" {
				scope[itemVar.Text], err = container.GetElement(idx)
				if err != nil {
					return nil, false, err
				}
			}

			_, forceReturn, err := i.executeBlock(scope, body, thisPackage, parentPackage)
			if err != nil {
				return nil, false, err
			}

			if forceReturn {
				return types.NilType{}, forceReturn, nil
			}
			//if result != nil {
			//	return result, forceReturn, nil
			//}
		}
	default:
		return nil, false, util.RuntimeError(fmt.Sprintf(
			"тип '%s' не є об'єктом, по якому можна ітерувати", container.GetTypeName(),
		))
	}

	return types.NilType{}, false, nil
}
