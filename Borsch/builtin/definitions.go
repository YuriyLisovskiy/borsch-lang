package builtin

import (
	"fmt"
	"os"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/cli/build"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

/*
types.FunctionInstance{
	Name:       "",
	Arguments: []types.FunctionArgument{},
	Code:       nil,
	Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {

	},
	ReturnType: ,
	IsBuiltin: true,
}
*/

var RuntimeObjects = map[string]types.Type{

	// I/O
	"друк": types.NewFunctionInstance(
		"друк",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "а",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			Print(*args...)
			return types.NilInstance{}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"друкр": types.NewFunctionInstance(
		"друкр",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "а",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			Print(append(*args, types.StringInstance{Value: "\n"})...)
			return types.NilInstance{}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"ввід": types.NewFunctionInstance(
		"ввід",
		[]types.FunctionArgument{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "повідомлення",
				IsVariadic: true,
				IsNullable: false,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			return Input(*args...)
		},
		types.FunctionReturnType{
			TypeHash:   types.StringTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),

	// Common
	"паніка": types.NewFunctionInstance(
		"паніка",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "повідомлення",
				IsVariadic: false,
				IsNullable: true,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			var strArgs []string
			for _, arg := range *args {
				strArgs = append(strArgs, arg.String())
			}

			return types.NilInstance{}, util.RuntimeError(strings.Join(strArgs, " "))
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"середовище": types.NewFunctionInstance(
		"середовище",
		[]types.FunctionArgument{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "ключ",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			return types.StringInstance{Value: os.Getenv((*args)[0].String())}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.StringTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"підтвердити": types.NewFunctionInstance(
		"підтвердити",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "очікуване",
				IsVariadic: false,
				IsNullable: true,
			},
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "фактичне",
				IsVariadic: false,
				IsNullable: true,
			},
			{
				TypeHash:   types.StringTypeHash,
				Name:       "повідомлення_про_помилку",
				IsVariadic: true,
				IsNullable: false,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			message := ""
			if len(*args) > 2 {
				messageArgs := (*args)[2:]
				sz := len(messageArgs)
				for c := 0; c < sz; c++ {
					message += messageArgs[c].String()
					if c < sz - 1 {
						message += " "
					}
				}
			}

			return types.NilInstance{}, Assert((*args)[0], (*args)[1], message)
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"авторське_право": types.NewFunctionInstance(
		"авторське_право",
		[]types.FunctionArgument{},
		func(*[]types.Type, *map[string]types.Type) (types.Type, error) {
			fmt.Printf("Copyright (c) %s %s.\nAll Rights Reserved.\n", build.Years, build.Author)
			return types.NilInstance{}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"ліцензія": types.NewFunctionInstance(
		"ліцензія",
		[]types.FunctionArgument{},
		func(*[]types.Type, *map[string]types.Type) (types.Type, error) {
			fmt.Println(build.License)
			return types.NilInstance{}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"допомога": types.NewFunctionInstance(
		"допомога",
		[]types.FunctionArgument{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "слово",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			return types.NilInstance{}, Help((*args)[0].String())
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),

	// System calls
	"вихід": types.NewFunctionInstance(
		"вихід",
		[]types.FunctionArgument{
			{
				TypeHash:   types.IntegerTypeHash,
				Name:       "код",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			return types.NilInstance{}, Exit((*args)[0].(types.IntegerInstance).Value)
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),

	// Conversion
	"дійсний": types.Real,
	"логічний": types.Bool,
	// "пакет": types.Package,
	"рядок": types.String,
	"словник": types.Dictionary,
	"список": types.List,
	"функція": types.Function,
	"цілий": types.Integer,

	// Utilities
	"довжина": types.NewFunctionInstance(
		"довжина",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "послідовність",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			return Length((*args)[0])
		},
		types.FunctionReturnType{
			TypeHash:   types.IntegerTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"додати": types.NewFunctionInstance(
		"додати",
		[]types.FunctionArgument{
			{
				TypeHash:   types.ListTypeHash,
				Name:       "вхідний_список",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "елементи",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			return AppendToList((*args)[0].(types.ListInstance), (*args)[1:]...)
		},
		types.FunctionReturnType{
			TypeHash:   types.ListTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"вилучити": types.NewFunctionInstance(
		"вилучити",
		[]types.FunctionArgument{
			{
				TypeHash:   types.DictionaryTypeHash,
				Name:       "вхідний_словник",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "ключ",
				IsVariadic: false,
				IsNullable: true,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			return RemoveFromDictionary((*args)[0].(types.DictionaryInstance), (*args)[1])
		},
		types.FunctionReturnType{
			TypeHash:   types.DictionaryTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"копіювати": types.NewFunctionInstance(
		"копіювати",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args *[]types.Type, _ *map[string]types.Type) (types.Type, error) {
			switch value := (*args)[0].(type) {
			case *types.ClassInstance:
				copied := value.Copy()
				return copied, nil
			default:
				return value, nil
			}
		},
		types.FunctionReturnType{
			TypeHash:   types.AnyTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
}

func init() {
	sourCream := types.NewClass("Сметанка", types.BuiltinPackage, map[string]types.Type{
		"ччч": types.NewStringInstance("Песто"),
	}, "", func() (types.Type, error) {
		return nil, nil
	})
	RuntimeObjects["Сметанка"] = sourCream

	RuntimeObjects["сметанка_для_борщу"] = types.NewClassInstance(sourCream, map[string]types.Type{})
}
