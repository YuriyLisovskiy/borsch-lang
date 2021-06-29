package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch/Borsch/cli/build"
	"github.com/YuriyLisovskiy/borsch/Borsch/util"
	"os"
	"strings"
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
	"друк": types.FunctionType{
		Name: "друк",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "а",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			Print(args...)
			return types.NilType{}, nil
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"друкр": types.FunctionType{
		Name: "друкр",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "а",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			Print(append(args, types.StringType{Value: "\n"})...)
			return types.NilType{}, nil
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"ввід": types.FunctionType{
		Name: "ввід",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "повідомлення",
				IsVariadic: true,
				IsNullable: false,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return Input(args...)
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.StringTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},

	// Common
	"паніка": types.FunctionType{
		Name: "паніка",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "повідомлення",
				IsVariadic: false,
				IsNullable: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			var strArgs []string
			for _, arg := range args {
				strArgs = append(strArgs, arg.String())
			}

			return types.NilType{}, util.RuntimeError(strings.Join(strArgs, " "))
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"середовище": types.FunctionType{
		Name: "середовище",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "ключ",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return types.StringType{Value: os.Getenv(args[0].String())}, nil
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.StringTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"підтвердити": types.FunctionType{
		Name: "підтвердити",
		Arguments: []types.FunctionArgument{
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
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
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
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"авторське_право": types.FunctionType{
		Name:      "авторське_право",
		Arguments: []types.FunctionArgument{},
		Callable: func([]types.ValueType, map[string]types.ValueType) (types.ValueType, error) {
			fmt.Printf("Copyright (c) %s %s.\nAll Rights Reserved.\n", build.Years, build.AuthorName)
			return types.NilType{}, nil
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"ліцензія": types.FunctionType{
		Name:      "ліцензія",
		Arguments: []types.FunctionArgument{},
		Callable: func([]types.ValueType, map[string]types.ValueType) (types.ValueType, error) {
			fmt.Println(build.License)
			return types.NilType{}, nil
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"допомога": types.FunctionType{
		Name: "допомога",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "слово",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return types.NilType{}, Help(args[0].String())
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},

	// System calls
	"вихід": types.FunctionType{
		Name: "вихід",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.IntegerTypeHash,
				Name:       "код",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		Callable: func(args []types.ValueType, kwargs map[string]types.ValueType) (types.ValueType, error) {
			return types.NilType{}, Exit(args[0].(types.IntegerType).Value)
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.NilTypeHash,
			IsNullable: true,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},

	// Conversion
	"цілий": types.FunctionType{
		Name: "цілий",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToInteger(args...)
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.IntegerTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"дійсний": types.FunctionType{
		Name: "дійсний",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToReal(args...)
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.RealTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"рядок": types.FunctionType{
		Name: "рядок",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToString(args...)
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.StringTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"логічний": types.FunctionType{
		Name: "логічний",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToBool(args...)
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.BoolTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"список": types.FunctionType{
		Name: "список",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "елементи",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToList(args...)
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.ListTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"словник": types.FunctionType{
		Name: "словник",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToDictionary(args...)
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.DictionaryTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},

	// Utilities
	"довжина": types.FunctionType{
		Name: "довжина",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "послідовність",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return Length(args[0])
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.IntegerTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"додати": types.FunctionType{
		Name: "додати",
		Arguments: []types.FunctionArgument{
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
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return AppendToList(args[0].(types.ListType), args[1:]...)
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.ListTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"вилучити": types.FunctionType{
		Name: "вилучити",
		Arguments: []types.FunctionArgument{
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
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return RemoveFromDictionary(args[0].(types.DictionaryType), args[1])
		},
		ReturnType: types.FunctionReturnType{
			TypeHash:   types.DictionaryTypeHash,
			IsNullable: false,
		},
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
}
