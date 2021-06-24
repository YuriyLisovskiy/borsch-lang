package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/cli/build"
	"github.com/YuriyLisovskiy/borsch/src/util"
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

var GlobalScope = map[string]types.ValueType{

	// I/O
	"друк": types.FunctionType{
		Name: "друк",
		Arguments: []types.FunctionArgument{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "а",
				IsVariadic: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			Print(args...)
			return nil, nil
		},
		ReturnType: types.NoneTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			Print(append(args, types.StringType{Value: "\n"})...)
			return nil, nil
		},
		ReturnType: types.NoneTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return Input(args...)
		},
		ReturnType: types.StringTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			var strArgs []string
			for _, arg := range args {
				strArgs = append(strArgs, arg.String())
			}

			return nil, util.RuntimeError(strings.Join(strArgs, " "))
		},
		ReturnType: types.NoneTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return types.StringType{Value: os.Getenv(args[0].String())}, nil
		},
		ReturnType: types.StringTypeHash,
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
			},
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "фактичне",
				IsVariadic: false,
			},
			{
				TypeHash:   types.StringTypeHash,
				Name:       "повідомлення_про_помилку",
				IsVariadic: false,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return nil, Assert(args[0], args[1], args[2].String())
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"авторське_право": types.FunctionType{
		Name:      "авторське_право",
		Arguments: []types.FunctionArgument{},
		Callable: func([]types.ValueType, map[string]types.ValueType) (types.ValueType, error) {
			fmt.Printf("Copyright (c) %s %s.\nAll Rights Reserved.\n", build.Years, build.AuthorName)
			return nil, nil
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
	"ліцензія": types.FunctionType{
		Name:      "ліцензія",
		Arguments: []types.FunctionArgument{},
		Callable: func([]types.ValueType, map[string]types.ValueType) (types.ValueType, error) {
			fmt.Println(build.License)
			return nil, nil
		},
		ReturnType: types.NoneTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return nil, Help(args[0].String())
		},
		ReturnType: types.NoneTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, kwargs map[string]types.ValueType) (types.ValueType, error) {
			return nil, Exit(args[0].(types.IntegerType).Value)
		},
		ReturnType: types.NoneTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToInteger(args...)
		},
		ReturnType: types.IntegerTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToReal(args...)
		},
		ReturnType: types.RealTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToString(args...)
		},
		ReturnType: types.StringTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToBool(args...)
		},
		ReturnType: types.BoolTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToList(args...)
		},
		ReturnType: types.ListTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToDictionary(args...)
		},
		ReturnType: types.DictionaryTypeHash,
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
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return Length(args[0])
		},
		ReturnType: types.IntegerTypeHash,
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
			},
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "елементи",
				IsVariadic: true,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return AppendToList(args[0].(types.ListType), args[1:]...)
		},
		ReturnType: types.ListTypeHash,
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
			},
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "ключ",
				IsVariadic: false,
			},
		},
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return RemoveFromDictionary(args[0].(types.DictionaryType), args[1])
		},
		ReturnType: types.DictionaryTypeHash,
		IsBuiltin:  true,
		Attributes: map[string]types.ValueType{},
	},
}
