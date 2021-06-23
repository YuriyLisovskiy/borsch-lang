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
	Parameters: []types.FunctionParameter{},
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
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "а",
				IsVariadic: true,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			Print(args...)
			return nil, nil
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
	},
	"друкр": types.FunctionType{
		Name: "друкр",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "а",
				IsVariadic: true,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			Print(append(args, types.StringType{Value: "\n"})...)
			return nil, nil
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
	},
	"ввід": types.FunctionType{
		Name: "ввід",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "повідомлення",
				IsVariadic: true,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return Input(args...)
		},
		ReturnType: types.StringTypeHash,
		IsBuiltin:  true,
	},

	// Common
	"паніка": types.FunctionType{
		Name: "паніка",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "повідомлення",
				IsVariadic: false,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			var strArgs []string
			for _, arg := range args {
				strArgs = append(strArgs, arg.String())
			}

			return nil, util.RuntimeError(strings.Join(strArgs, " "))
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
	},
	"середовище": types.FunctionType{
		Name: "середовище",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "ключ",
				IsVariadic: false,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return types.StringType{Value: os.Getenv(args[0].String())}, nil
		},
		ReturnType: types.StringTypeHash,
		IsBuiltin:  true,
	},
	"підтвердити": types.FunctionType{
		Name: "підтвердити",
		Parameters: []types.FunctionParameter{
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
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return nil, Assert(args[0], args[1], args[2].String())
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
	},
	"авторське_право": types.FunctionType{
		Name:       "авторське_право",
		Parameters: []types.FunctionParameter{},
		Code:       nil,
		Callable: func([]types.ValueType, map[string]types.ValueType) (types.ValueType, error) {
			fmt.Printf("Copyright (c) %s %s.\nAll Rights Reserved.\n", build.Years, build.AuthorName)
			return nil, nil
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
	},
	"ліцензія": types.FunctionType{
		Name:       "ліцензія",
		Parameters: []types.FunctionParameter{},
		Code:       nil,
		Callable: func([]types.ValueType, map[string]types.ValueType) (types.ValueType, error) {
			fmt.Println(build.License)
			return nil, nil
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
	},
	"допомога": types.FunctionType{
		Name: "допомога",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.StringTypeHash,
				Name:       "слово",
				IsVariadic: false,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return nil, Help(args[0].String())
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
	},

	// System calls
	"вихід": types.FunctionType{
		Name: "вихід",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.IntegerTypeHash,
				Name:       "код",
				IsVariadic: false,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, kwargs map[string]types.ValueType) (types.ValueType, error) {
			return nil, Exit(args[0].(types.IntegerType).Value)
		},
		ReturnType: types.NoneTypeHash,
		IsBuiltin:  true,
	},

	// Conversion
	"цілий": types.FunctionType{
		Name: "цілий",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToInteger(args...)
		},
		ReturnType: types.IntegerTypeHash,
		IsBuiltin:  true,
	},
	"дійсний": types.FunctionType{
		Name: "дійсний",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToReal(args...)
		},
		ReturnType: types.RealTypeHash,
		IsBuiltin:  true,
	},
	"рядок": types.FunctionType{
		Name: "рядок",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToString(args...)
		},
		ReturnType: types.StringTypeHash,
		IsBuiltin:  true,
	},
	"логічний": types.FunctionType{
		Name: "логічний",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToBool(args...)
		},
		ReturnType: types.BoolTypeHash,
		IsBuiltin:  true,
	},
	"список": types.FunctionType{
		Name: "список",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "елементи",
				IsVariadic: true,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToList(args...)
		},
		ReturnType: types.ListTypeHash,
		IsBuiltin:  true,
	},
	"словник": types.FunctionType{
		Name: "словник",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "значення",
				IsVariadic: true,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return ToDictionary(args...)
		},
		ReturnType: types.DictionaryTypeHash,
		IsBuiltin:  true,
	},

	// Utilities
	"довжина": types.FunctionType{
		Name: "довжина",
		Parameters: []types.FunctionParameter{
			{
				TypeHash:   types.AnyTypeHash,
				Name:       "послідовність",
				IsVariadic: false,
			},
		},
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return Length(args[0])
		},
		ReturnType: types.IntegerTypeHash,
		IsBuiltin:  true,
	},
	"додати": types.FunctionType{
		Name: "додати",
		Parameters: []types.FunctionParameter{
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
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return AppendToList(args[0].(types.ListType), args[1:]...)
		},
		ReturnType: types.ListTypeHash,
		IsBuiltin:  true,
	},
	"вилучити": types.FunctionType{
		Name: "вилучити",
		Parameters: []types.FunctionParameter{
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
		Code: nil,
		Callable: func(args []types.ValueType, _ map[string]types.ValueType) (types.ValueType, error) {
			return RemoveFromDictionary(args[0].(types.DictionaryType), args[1])
		},
		ReturnType: types.DictionaryTypeHash,
		IsBuiltin:  true,
	},
}
