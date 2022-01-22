package builtin

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func DeepCopy(object common.Type) (common.Type, error) {
	switch value := object.(type) {
	case *types.ClassInstance:
		copied := value.Copy()
		return copied, nil
	default:
		return value, nil
	}
}

func runOperator(
	ctx common.Context,
	name string,
	object common.Type,
	expectedType *types.Class,
) (common.Type, error) {
	attribute, err := object.GetAttribute(name)
	if err != nil {
		return nil, util.RuntimeError(fmt.Sprintf("об'єкт типу '%s' не має довжини", object.GetTypeName()))
	}

	switch operator := attribute.(type) {
	case *types.FunctionInstance:
		args := []common.Type{object}
		kwargs := map[string]common.Type{operator.Arguments[0].Name: object}
		if err := types.CheckFunctionArguments(ctx, operator, &args, &kwargs); err != nil {
			return nil, err
		}

		result, err := operator.Call(ctx, &args, &kwargs)
		if err != nil {
			return nil, err
		}

		if err := types.CheckResult(ctx, result, operator); err != nil {
			return nil, err
		}

		if result.(types.ObjectInstance).GetPrototype() != expectedType {
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"'%s' має повертати значення з типом '%s', отримано '%s'",
					name, expectedType.GetTypeName(), result.GetTypeName(),
				),
			)
		}

		return result, nil
	default:
		return nil, util.ObjectIsNotCallable(name, attribute.GetTypeName())
	}
}

func ImportPackage(
	baseScope map[string]common.Type,
	newPackagePath string,
	parser common.Parser,
	parentPackage *types.PackageInstance,
) (common.Type, error) {
	if strings.HasPrefix(newPackagePath, "!/") {
		newPackagePath = path.Join(os.Getenv(common.BORSCH_LIB), newPackagePath[2:])
	} else if !path.IsAbs(newPackagePath) {
		var err error
		if parentPackage != nil {
			baseDir := path.Dir(parentPackage.Name)
			newPackagePath = path.Join(baseDir, newPackagePath)
		} else {
			newPackagePath, err = filepath.Abs(newPackagePath)
		}

		if err != nil {
			return nil, err
		}
	}

	currPackage := parentPackage
	for currPackage != nil {
		if currPackage.Name == newPackagePath {
			return nil, util.RuntimeError("циклічний імпорт заборонений")
		}

		currPackage = currPackage.Parent
	}

	packageCode, err := util.ReadFile(newPackagePath)
	if err != nil {
		return nil, err
	}

	ast, err := parser.Parse(newPackagePath, string(packageCode))
	if err != nil {
		return nil, err
	}

	context := parser.NewContext(newPackagePath, parentPackage)
	context.PushScope(baseScope)
	package_, err := ast.Evaluate(context)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Відстеження (стек викликів):\n%s", err.Error()))
	}

	return package_, nil
}
