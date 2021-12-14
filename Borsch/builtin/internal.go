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
types.FunctionType{
	Name:       "",
	Arguments: []types.FunctionArgument{},
	Code:       nil,
	Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {

	},
	ReturnType: ,
	IsBuiltin: true,
}
*/

var RuntimeFunctions = map[string]types.ValueType{

	// I/O
	"друк": types.NewFunctionType(
		"друк",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "а",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			Print(args...)
			return types.NilType{}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"друкр": types.NewFunctionType(
		"друкр",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "а",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			Print(append(args, types.StringType{Value: "\n"})...)
			return types.NilType{}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"ввід": types.NewFunctionType(
		"ввід",
		[]types.FunctionArgument{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "повідомлення",
				IsVariadic: true,
				IsNullable: false,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return Input(args...)
		},
		types.FunctionReturnType{
			TypeHash:   types.StringTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),

	// Common
	"паніка": types.NewFunctionType(
		"паніка",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "повідомлення",
				IsVariadic: false,
				IsNullable: true,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			var strArgs []string
			for _, arg := range args {
				strArgs = append(strArgs, arg.String())
			}

			return types.NilType{}, util.RuntimeError(strings.Join(strArgs, " "))
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"середовище": types.NewFunctionType(
		"середовище",
		[]types.FunctionArgument{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "ключ",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return types.StringType{Value: os.Getenv(args[0].String())}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.StringTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"підтвердити": types.NewFunctionType(
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
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			message := ""
			if len(args) > 2 {
				messageArgs := args[2:]
				sz := len(messageArgs)
				for c := 0; c < sz; c++ {
					message += messageArgs[c].String()
					if c < sz - 1 {
						message += " "
					}
				}
			}

			return types.NilType{}, Assert(args[0], args[1], message)
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"авторське_право": types.NewFunctionType(
		"авторське_право",
		[]types.FunctionArgument{},
		func([]types.ValueType, map[string]types.ValueType) (types.ValueType, error) {
			fmt.Printf("Copyright (c) %s %s.\nAll Rights Reserved.\n", build.Years, build.Author)
			return types.NilType{}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"ліцензія": types.NewFunctionType(
		"ліцензія",
		[]types.FunctionArgument{},
		func([]types.ValueType, map[string]types.ValueType) (types.ValueType, error) {
			fmt.Println(build.License)
			return types.NilType{}, nil
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"допомога": types.NewFunctionType(
		"допомога",
		[]types.FunctionArgument{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "слово",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return types.NilType{}, Help(args[0].String())
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),

	// System calls
	"вихід": types.NewFunctionType(
		"вихід",
		[]types.FunctionArgument{
			{
				TypeHash:   types.IntegerTypeHash,
				Name:       "код",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args []types.ValueType, kwargs map[string]types.ValueType) (types.ValueType, error) {
			return types.NilType{}, Exit(args[0].(types.IntegerType).Value)
		},
		types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),

	// Conversion
	"цілий": types.NewFunctionType(
		"цілий",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToInteger(args...)
		},
		types.FunctionReturnType{
			TypeHash:   types.IntegerTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"дійсний": types.NewFunctionType(
		"дійсний",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToReal(args...)
		},
		types.FunctionReturnType{
			TypeHash:   types.RealTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"рядок": types.NewFunctionType(
		"рядок",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToString(args...)
		},
		types.FunctionReturnType{
			TypeHash:   types.StringTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"логічний": types.NewFunctionType(
		"логічний",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToBool(args...)
		},
		types.FunctionReturnType{
			TypeHash:   types.BoolTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"список": types.NewFunctionType(
		"список",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "елементи",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToList(args...)
		},
		types.FunctionReturnType{
			TypeHash:   types.ListTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"словник": types.NewFunctionType(
		"словник",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToDictionary(args...)
		},
		types.FunctionReturnType{
			TypeHash:   types.DictionaryTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),

	// Utilities
	"довжина": types.NewFunctionType(
		"довжина",
		[]types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "послідовність",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return Length(args[0])
		},
		types.FunctionReturnType{
			TypeHash:   types.IntegerTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"додати": types.NewFunctionType(
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
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return AppendToList(args[0].(types.ListType), args[1:]...)
		},
		types.FunctionReturnType{
			TypeHash:   types.ListTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
	"вилучити": types.NewFunctionType(
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
		func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return RemoveFromDictionary(args[0].(types.DictionaryType), args[1])
		},
		types.FunctionReturnType{
			TypeHash:   types.DictionaryTypeHash,
			IsNullable: false,
		},
		types.BuiltinPackage,
		"", // TODO: add doc
	),
}
