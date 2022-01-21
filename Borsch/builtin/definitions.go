package builtin

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/cli/build"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	types "github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

var BuiltinScope map[string]common.Type

func init() {
	BuiltinScope = map[string]common.Type{

		// I/O
		"друк": types.NewFunctionInstance(
			"друк",
			[]types.FunctionArgument{
				{
					Type:       types.Any,
					Name:       "а",
					IsVariadic: true,
					IsNullable: true,
				},
			},
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				Print(ctx, *args...)
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
		),
		"друкр": types.NewFunctionInstance(
			"друкр",
			[]types.FunctionArgument{
				{
					Type:       types.Any,
					Name:       "а",
					IsVariadic: true,
					IsNullable: true,
				},
			},
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				Print(ctx, append(*args, types.StringInstance{Value: "\n"})...)
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
		),
		"ввід": types.NewFunctionInstance(
			"ввід",
			[]types.FunctionArgument{
				{
					Type:       types.String,
					Name:       "повідомлення",
					IsVariadic: true,
					IsNullable: false,
				},
			},
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				return Input(ctx, *args...)
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
		),

		// Common
		"паніка": types.NewFunctionInstance(
			"паніка",
			[]types.FunctionArgument{
				{
					Type:       types.Any,
					Name:       "повідомлення",
					IsVariadic: false,
					IsNullable: true,
				},
			},
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				var strArgs []string
				for _, arg := range *args {
					strArgs = append(strArgs, arg.String(ctx))
				}

				return types.NewNilInstance(), util.RuntimeError(strings.Join(strArgs, " "))
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
		),
		"середовище": types.NewFunctionInstance(
			"середовище",
			[]types.FunctionArgument{
				{
					Type:       types.String,
					Name:       "ключ",
					IsVariadic: false,
					IsNullable: false,
				},
			},
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				return types.StringInstance{Value: os.Getenv((*args)[0].String(ctx))}, nil
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
		),
		"підтвердити": types.NewFunctionInstance(
			"підтвердити",
			[]types.FunctionArgument{
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
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				message := ""
				if len(*args) > 2 {
					messageArgs := (*args)[2:]
					sz := len(messageArgs)
					for c := 0; c < sz; c++ {
						message += messageArgs[c].String(ctx)
						if c < sz-1 {
							message += " "
						}
					}
				}

				return types.NewNilInstance(), Assert(ctx, (*args)[0], (*args)[1], message)
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
		),
		"авторське_право": types.NewFunctionInstance(
			"авторське_право",
			[]types.FunctionArgument{},
			func(common.Context, *[]common.Type, *map[string]common.Type) (common.Type, error) {
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
		),
		"ліцензія": types.NewFunctionInstance(
			"ліцензія",
			[]types.FunctionArgument{},
			func(common.Context, *[]common.Type, *map[string]common.Type) (common.Type, error) {
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
		),
		"допомога": types.NewFunctionInstance(
			"допомога",
			[]types.FunctionArgument{
				{
					Type:       types.String,
					Name:       "слово",
					IsVariadic: false,
					IsNullable: false,
				},
			},
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				return types.NewNilInstance(), Help((*args)[0].String(ctx))
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
		),

		// System calls
		"вихід": types.NewFunctionInstance(
			"вихід",
			[]types.FunctionArgument{
				{
					Type:       types.Integer,
					Name:       "код",
					IsVariadic: false,
					IsNullable: false,
				},
			},
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
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
		),

		// Conversion
		"дійсний":  types.Real,
		"логічний": types.Bool,
		"пакет": types.NewFunctionInstance(
			"пакет",
			[]types.FunctionArgument{
				{
					Type:       types.String,
					Name:       "шлях",
					IsVariadic: false,
					IsNullable: false,
				},
			},
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				packagePath := (*args)[0].(types.StringInstance).Value
				if strings.HasPrefix(packagePath, "!/") {
					packagePath = path.Join(os.Getenv(common.BORSCH_LIB), packagePath[2:])
				}

				pack, err := ImportPackage(BuiltinScope, packagePath, ctx.GetParser())
				if err != nil {
					return nil, err
				}

				return pack, nil
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
		),
		"рядок":   types.String,
		"словник": types.Dictionary,
		"список":  types.List,
		// "функція": types.Function,
		"цілий": types.Integer,

		// Utilities
		"довжина": types.NewFunctionInstance(
			"довжина",
			[]types.FunctionArgument{
				{
					Type:       types.Any,
					Name:       "послідовність",
					IsVariadic: false,
					IsNullable: false,
				},
			},
			func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				sequence := (*args)[0]
				return runOperator(
					ctx,
					ops.LengthOperatorName,
					sequence,
					types.Integer,
				)
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
		),
		"додати": types.NewFunctionInstance(
			"додати",
			[]types.FunctionArgument{
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
			func(_ common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
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
		),
		"копіювати": types.NewFunctionInstance(
			"копіювати",
			[]types.FunctionArgument{
				{
					Type:       types.Any,
					Name:       "значення",
					IsVariadic: false,
					IsNullable: false,
				},
			},
			func(_ common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
				return DeepCopy((*args)[0])
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
		),
		"тип": types.NewFunctionInstance(
			"тип",
			[]types.FunctionArgument{
				{
					Type:       types.Any,
					Name:       "значення",
					IsVariadic: false,
					IsNullable: false,
				},
			},
			func(_ common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
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
		),
	}

	types.BuiltinPackage.Attributes = BuiltinScope
}
