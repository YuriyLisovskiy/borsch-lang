package types

import (
	"fmt"
	"strconv"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func ToInteger(_ common.Context, args ...common.Type) (common.Type, error) {
	if len(args) == 0 {
		return NewIntegerInstance(0), nil
	}

	if len(args) != 1 {
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"'цілий()' приймає лише один аргумент (отримано %d)", len(args),
			),
		)
	}

	switch vt := args[0].(type) {
	case RealInstance:
		return NewIntegerInstance(int64(vt.Value)), nil
	case IntegerInstance:
		return vt, nil
	case StringInstance:
		intVal, err := strconv.ParseInt(vt.Value, 10, 64)
		if err != nil {
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"некоректний літерал для функції 'цілий()' з основою 10: '%s'", vt.Value,
				),
			)
		}

		return NewIntegerInstance(intVal), nil
	case BoolInstance:
		if vt.Value {
			return NewIntegerInstance(1), nil
		}

		return NewIntegerInstance(0), nil
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"'%s' неможливо інтерпретувати як ціле число", args[0].GetTypeName(),
			),
		)
	}
}

func ToReal(_ common.Context, args ...common.Type) (common.Type, error) {
	if len(args) == 0 {
		return NewRealInstance(0.0), nil
	}

	if len(args) != 1 {
		return nil, util.RuntimeError(
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
			return nil, util.RuntimeError(
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
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"'%s' неможливо інтерпретувати як дійсне число", args[0].GetTypeName(),
			),
		)
	}
}

func ToString(ctx common.Context, args ...common.Type) (common.Type, error) {
	if len(args) == 0 {
		return NewStringInstance(""), nil
	}

	if len(args) != 1 {
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"функція 'рядок()' приймає лише один аргумент (отримано %d)", len(args),
			),
		)
	}

	return NewStringInstance(args[0].String(ctx)), nil

	// TODO: remove
	// switch vt := args[0].(type) {
	// case StringInstance:
	// 	return vt, nil
	// case RealInstance, IntegerInstance, BoolInstance, NilInstance:
	// 	return NewStringInstance(vt.String()), nil
	// default:
	// 	return NewStringInstance(vt.String()), nil
	// 	// return nil, util.RuntimeError(fmt.Sprintf(
	// 	// 	"'%s' неможливо інтерпретувати як рядок", args[0].GetTypeName(),
	// 	// ))
	// }
}

func ToBool(_ common.Context, args ...common.Type) (common.Type, error) {
	if len(args) == 0 {
		return NewBoolInstance(false), nil
	}

	if len(args) != 1 {
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"функція 'логічний()' приймає лише один аргумент (отримано %d)", len(args),
			),
		)
	}

	switch vt := args[0].(type) {
	case RealInstance:
		return NewBoolInstance(vt.Value != 0.0), nil
	case IntegerInstance:
		return NewBoolInstance(vt.Value != 0), nil
	case StringInstance:
		return NewBoolInstance(vt.Value != ""), nil
	case BoolInstance:
		return vt, nil
	case NilInstance:
		return NewBoolInstance(false), nil
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"'%s' неможливо інтерпретувати як логічне значення", args[0].GetTypeName(),
			),
		)
	}
}

func ToList(_ common.Context, args ...common.Type) (common.Type, error) {
	list := NewListInstance()
	if len(args) == 0 {
		return list, nil
	}

	for _, arg := range args {
		list.Values = append(list.Values, arg)
	}

	return list, nil
}

func ToDictionary(ctx common.Context, args ...common.Type) (common.Type, error) {
	dict := NewDictionaryInstance()
	if len(args) == 0 {
		return dict, nil
	}

	if len(args) != 2 {
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"функція 'словник()' приймає два аргументи, або жодного (отримано %d)", len(args),
			),
		)
	}

	switch keys := args[0].(type) {
	case ListInstance:
		switch values := args[1].(type) {
		case ListInstance:
			if keys.Length(ctx) != values.Length(ctx) {
				return nil, util.RuntimeError(
					fmt.Sprintf(
						"довжина списку ключів має співпадати з довжиною списку значень",
					),
				)
			}

			length := keys.Length(ctx)
			for i := int64(0); i < length; i++ {
				err := dict.SetElement(keys.Values[i], values.Values[i])
				if err != nil {
					return nil, err
				}
			}

			return dict, nil
		default:
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"функція 'словник()' другим аргументом приймає список значень",
				),
			)
		}
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"функція 'словник()' першим аргументом приймає список ключів",
			),
		)
	}
}
