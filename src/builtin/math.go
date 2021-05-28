package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"math"
)

func Log10(values... ValueType) (ValueType, error) {
	if len(values) != 1 {
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"лог10() приймає лише один аргумент (отримано %d)", len(values),
		))
	}

	switch vt := values[0].(type) {
	case RealNumberType:
		return RealNumberType{Value: math.Log10(vt.Value)}, nil
	default:
		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"лог10() приймає лише дійсне значення, а не '%s'", values[0].TypeName(),
		))
	}
}
