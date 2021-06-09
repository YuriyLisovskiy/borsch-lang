package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"strconv"
)

func ToInteger(args ...types.ValueType) (types.ValueType, error) {
	if len(args) == 0 {
		return types.IntegerType{Value: 0}, nil
	}

	if len(args) != 1 {
		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"функція 'цілий()' приймає лише один аргумент (отримано %d)", len(args),
		))
	}

	switch vt := args[0].(type) {
	case types.RealType:
		return types.IntegerType{Value: int64(vt.Value)}, nil
	case types.IntegerType:
		return vt, nil
	case types.StringType:
		intVal, err := strconv.ParseInt(vt.Value, 10, 64)
		if err != nil {
			return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
				"некоректний літерал для функції 'цілий()' з основою 10: '%s'", vt.Value,
			))
		}

		return types.IntegerType{Value: intVal}, nil
	case types.BoolType:
		if vt.Value {
			return types.IntegerType{Value: 1}, nil
		}

		return types.IntegerType{Value: 0}, nil
	default:
		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"'%s' неможливо інтерпретувати як ціле число", args[0].TypeName(),
		))
	}
}

func ToReal(args ...types.ValueType) (types.ValueType, error) {
	if len(args) == 0 {
		return types.RealType{Value: 0.0}, nil
	}

	if len(args) != 1 {
		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"функція 'дійсний()' приймає лише один аргумент (отримано %d)", len(args),
		))
	}

	switch vt := args[0].(type) {
	case types.RealType:
		return vt, nil
	case types.IntegerType:
		return types.RealType{Value: float64(vt.Value)}, nil
	case types.StringType:
		realVal, err := strconv.ParseFloat(vt.Value, 64)
		if err != nil {
			return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
				"не вдалося перетворити рядок у дійсне число: '%s'", vt.Value,
			))
		}

		return types.RealType{Value: realVal}, nil
	case types.BoolType:
		if vt.Value {
			return types.RealType{Value: 1.0}, nil
		}

		return types.RealType{Value: 0.0}, nil
	default:
		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"'%s' неможливо інтерпретувати як дійсне число", args[0].TypeName(),
		))
	}
}

func ToString(args ...types.ValueType) (types.ValueType, error) {
	if len(args) == 0 {
		return types.StringType{Value: ""}, nil
	}

	if len(args) != 1 {
		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"функція 'рядок()' приймає лише один аргумент (отримано %d)", len(args),
		))
	}

	switch vt := args[0].(type) {
	case types.StringType:
		return vt, nil
	case types.RealType, types.IntegerType, types.BoolType, types.NoneType:
		return types.StringType{Value: vt.String()}, nil
	default:
		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"'%s' неможливо інтерпретувати як рядок", args[0].TypeName(),
		))
	}
}

func ToBool(args ...types.ValueType) (types.ValueType, error) {
	if len(args) == 0 {
		return types.BoolType{Value: false}, nil
	}

	if len(args) != 1 {
		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"функція 'логічний()' приймає лише один аргумент (отримано %d)", len(args),
		))
	}

	switch vt := args[0].(type) {
	case types.RealType:
		return types.BoolType{Value: vt.Value != 0.0}, nil
	case types.IntegerType:
		return types.BoolType{Value: vt.Value != 0}, nil
	case types.StringType:
		return types.BoolType{Value: vt.Value != ""}, nil
	case types.BoolType:
		return vt, nil
	case types.NoneType:
		return types.BoolType{Value: false}, nil
	default:
		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"'%s' неможливо інтерпретувати як логічне значення", args[0].TypeName(),
		))
	}
}
