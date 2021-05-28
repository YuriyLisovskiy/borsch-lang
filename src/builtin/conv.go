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
				"недійсний літерал для функції ціле() з основою 10: '%s'", vt.Value,
			))
		}

		return IntegerNumberType{Value: intVal}, nil
	default:
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"'%s' неможливо інтерпретувати як ціле число", args[0].TypeName(),
		))
	}
}
