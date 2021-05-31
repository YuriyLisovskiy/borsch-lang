package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"strconv"
)

func CastToInt(args... ValueType) (ValueType, error) {
	if len(args) != 1 {
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"ціле() приймає лише один аргумент (отримано %d)", len(args),
		))
	}

	switch vt := args[0].(type) {
	case RealNumberType:
		return IntegerNumberType{Value: int64(vt.Value)}, nil
	case IntegerNumberType:
		return vt, nil
	case StringType:
		intVal, err := strconv.ParseInt(vt.Value, 10, 64)
		if err != nil {
			return NoneType{}, util.RuntimeError(fmt.Sprintf(
				"некоректний літерал для функції ціле() з основою 10: '%s'", vt.Value,
			))
		}

		return IntegerNumberType{Value: intVal}, nil
	case BoolType:
		if vt.Value {
			return IntegerNumberType{Value: 1}, nil
		}

		return IntegerNumberType{Value: 0}, nil
	default:
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"'%s' неможливо інтерпретувати як ціле число", args[0].TypeName(),
		))
	}
}

func CastToReal(args... ValueType) (ValueType, error) {
	if len(args) != 1 {
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"дійсне() приймає лише один аргумент (отримано %d)", len(args),
		))
	}

	switch vt := args[0].(type) {
	case RealNumberType:
		return vt, nil
	case IntegerNumberType:
		return RealNumberType{Value: float64(vt.Value)}, nil
	case StringType:
		realVal, err := strconv.ParseFloat(vt.Value, 64)
		if err != nil {
			return NoneType{}, util.RuntimeError(fmt.Sprintf(
				"не вдалося перетворити рядок у дійсне число: '%s'", vt.Value,
			))
		}

		return RealNumberType{Value: realVal}, nil
	case BoolType:
		if vt.Value {
			return RealNumberType{Value: 1.0}, nil
		}

		return RealNumberType{Value: 0.0}, nil
	default:
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"'%s' неможливо інтерпретувати як дійсне число", args[0].TypeName(),
		))
	}
}

func CastToString(args... ValueType) (ValueType, error) {
	if len(args) != 1 {
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"рядок() приймає лише один аргумент (отримано %d)", len(args),
		))
	}

	switch vt := args[0].(type) {
	case StringType:
		return vt, nil
	case RealNumberType, IntegerNumberType, BoolType, NoneType:
		return StringType{Value: vt.String()}, nil
	default:
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"'%s' неможливо інтерпретувати як рядок", args[0].TypeName(),
		))
	}
}

func CastToBool(args... ValueType) (ValueType, error) {
	if len(args) != 1 {
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"логічне() приймає лише один аргумент (отримано %d)", len(args),
		))
	}

	switch vt := args[0].(type) {
	case RealNumberType:
		return BoolType{Value: vt.Value != 0.0}, nil
	case IntegerNumberType:
		return BoolType{Value: vt.Value != 0}, nil
	case StringType:
		return BoolType{Value: vt.Value != ""}, nil
	case BoolType:
		return vt, nil
	case NoneType:
		return BoolType{Value: false}, nil
	default:
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"'%s' неможливо інтерпретувати як логічне значення", args[0].TypeName(),
		))
	}
}
