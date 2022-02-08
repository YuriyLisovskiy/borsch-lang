package types

import (
	"fmt"
	"strconv"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

func ToReal(_ common.State, args ...common.Value) (common.Value, error) {
	if len(args) == 0 {
		return NewRealInstance(0.0), nil
	}

	if len(args) != 1 {
		return nil, utilities.RuntimeError(
			fmt.Sprintf(
				"функція 'дійсний()' приймає лише один аргумент (отримано %d)", len(args),
			),
		)
	}

	switch vt := args[0].(type) {
	case RealInstance:
		return vt, nil
	case IntegerInstance:
		return NewRealInstance(float64(vt.Value)), nil
	case StringInstance:
		realVal, err := strconv.ParseFloat(vt.Value, 64)
		if err != nil {
			return nil, utilities.RuntimeError(
				fmt.Sprintf(
					"не вдалося перетворити рядок у дійсне число: '%s'", vt.Value,
				),
			)
		}

		return NewRealInstance(realVal), nil
	case BoolInstance:
		if vt.Value {
			return NewRealInstance(1.0), nil
		}

		return NewRealInstance(0.0), nil
	default:
		return nil, utilities.RuntimeError(
			fmt.Sprintf(
				"'%s' неможливо інтерпретувати як дійсне число", args[0].GetTypeName(),
			),
		)
	}
}

func ToString(state common.State, args ...common.Value) (common.Value, error) {
	if len(args) == 0 {
		return NewStringInstance(""), nil
	}

	if len(args) != 1 {
		return nil, utilities.RuntimeError(
			fmt.Sprintf(
				"функція 'рядок()' приймає лише один аргумент (отримано %d)", len(args),
			),
		)
	}

	argStr, err := args[0].String(state)
	if err != nil {
		return nil, err
	}

	return NewStringInstance(argStr), nil
}
