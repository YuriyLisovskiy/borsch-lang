package builtin

import (
	"errors"
	"fmt"
	"os"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/std"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/cli/build"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

var (
	PrintFunction     *types.FunctionInstance
	PrintLineFunction *types.FunctionInstance
	InputFunction     *types.FunctionInstance
	PanicFunction     *types.FunctionInstance
	EnvFunction       *types.FunctionInstance
	AssertFunction    *types.FunctionInstance
	CopyrightFunction *types.FunctionInstance
	LicenceFunction   *types.FunctionInstance
	HelpFunction      *types.FunctionInstance
	ExitFunction      *types.FunctionInstance
	ImportFunction    *types.FunctionInstance
	LengthFunction    *types.FunctionInstance
	AddToListFunction *types.FunctionInstance
	DeepCopyFunction  *types.FunctionInstance
	TypeFunction      *types.FunctionInstance
)

func initRuntime() {
	types.Init()
	std.Init()
	PrintFunction = types.NewFunctionInstance(
		"друк",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "х",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			return types.NewNilInstance(), Print(state, *args...)
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	PrintLineFunction = types.NewFunctionInstance(
		"друкр",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "а",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			return types.NewNilInstance(), Print(state, append(*args, types.StringInstance{Value: "\n"})...)
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	InputFunction = types.NewFunctionInstance(
		"ввід",
		[]types.FunctionParameter{
			{
				Type:       types.String,
				Name:       "повідомлення",
				IsVariadic: true,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			return Input(state, *args...)
		},
		[]types.FunctionReturnType{
			{
				Type:       types.String,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	PanicFunction = types.NewFunctionInstance(
		"панікувати",
		[]types.FunctionParameter{
			{
				Type:       std.ErrorClass,
				Name:       "помилка",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			self := (*args)[0]
			msg, err := self.String(state)
			if err != nil {
				return nil, err
			}

			return types.NewNilInstance(), errors.New(fmt.Sprintf("%s: %s", self.GetTypeName(), msg))
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	EnvFunction = types.NewFunctionInstance(
		"середовище",
		[]types.FunctionParameter{
			{
				Type:       types.String,
				Name:       "ключ",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			argStr, err := (*args)[0].String(state)
			if err != nil {
				return nil, err
			}

			return types.StringInstance{Value: os.Getenv(argStr)}, nil
		},
		[]types.FunctionReturnType{
			{
				Type:       types.String,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	AssertFunction = types.NewFunctionInstance(
		"підтвердити",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "очікуване",
				IsVariadic: false,
				IsNullable: true,
			},
			{
				Type:       types.Any,
				Name:       "фактичне",
				IsVariadic: false,
				IsNullable: true,
			},
			{
				Type:       types.String,
				Name:       "повідомлення_про_помилку",
				IsVariadic: true,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			message := ""
			if len(*args) > 2 {
				messageArgs := (*args)[2:]
				sz := len(messageArgs)
				for c := 0; c < sz; c++ {
					argStr, err := messageArgs[c].String(state)
					if err != nil {
						return nil, err
					}

					message += argStr
					if c < sz-1 {
						message += " "
					}
				}
			}

			return types.NewNilInstance(), Assert(state, (*args)[0], (*args)[1], message)
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	CopyrightFunction = types.NewFunctionInstance(
		"авторське_право",
		[]types.FunctionParameter{},
		func(common.State, *[]common.Value, *map[string]common.Value) (common.Value, error) {
			fmt.Printf("Copyright (c) %s %s.\nAll Rights Reserved.\n", build.Years, build.Author)
			return types.NewNilInstance(), nil
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	LicenceFunction = types.NewFunctionInstance(
		"ліцензія",
		[]types.FunctionParameter{},
		func(common.State, *[]common.Value, *map[string]common.Value) (common.Value, error) {
			fmt.Println(build.License)
			return types.NewNilInstance(), nil
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	HelpFunction = types.NewFunctionInstance(
		"допомога",
		[]types.FunctionParameter{
			{
				Type:       types.String,
				Name:       "слово",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			argStr, err := (*args)[0].String(state)
			if err != nil {
				return nil, err
			}

			return types.NewNilInstance(), Help(argStr)
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	ExitFunction = types.NewFunctionInstance(
		"вихід",
		[]types.FunctionParameter{
			{
				Type:       types.Integer,
				Name:       "код",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(_ common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			os.Exit(int((*args)[0].(types.IntegerInstance).Value))
			return types.NewNilInstance(), nil
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	ImportFunction = types.NewFunctionInstance(
		"імпорт",
		[]types.FunctionParameter{
			{
				Type:       types.String,
				Name:       "шлях",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			return state.GetInterpreter().Import(
				state,
				(*args)[0].(types.StringInstance).Value,
				// state.GetCurrentPackage().(*types.PackageInstance),
			)
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Package,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	LengthFunction = types.NewFunctionInstance(
		"довжина",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "послідовність",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			sequence := (*args)[0]
			if !sequence.HasAttribute(common.LengthOperatorName) {
				return nil, errors.New(fmt.Sprintf("об'єкт типу '%s' не має довжини", sequence.GetTypeName()))
			}

			return runUnaryOperator(state, common.LengthOperatorName, sequence, types.Integer)
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Integer,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	AddToListFunction = types.NewFunctionInstance(
		"додати",
		[]types.FunctionParameter{
			{
				Type:       types.List,
				Name:       "вхідний_список",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       types.Any,
				Name:       "елементи",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(_ common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			list := (*args)[0].(types.ListInstance)
			values := (*args)[1:]
			for _, value := range values {
				list.Values = append(list.Values, value)
			}

			return list, nil
		},
		[]types.FunctionReturnType{
			{
				Type:       types.List,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	DeepCopyFunction = types.NewFunctionInstance(
		"копіювати",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "значення",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(_ common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			return deepCopy((*args)[0])
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Any,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)

	TypeFunction = types.NewFunctionInstance(
		"тип",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "значення",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(_ common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			return types.GetTypeOfInstance((*args)[0])
		},
		[]types.FunctionReturnType{
			{
				Type:       types.Any,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
