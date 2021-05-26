package std

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func Sum(values... string) (string, error) {
	res := 0.0
	for _, value := range values {
		number, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return "", errors.New(
				fmt.Sprintf("стд::сума: неможливо виконати додавання з нечисловим значенням '%s'", value),
			)
		}

		res += number
	}

	return fmt.Sprintf("%f", res), nil
}

func Log10(values... string) (string, error) {
	if len(values) != 1 {
		return "", errors.New(fmt.Sprintf("лог10: функція приймає лише один аргумент"))
	}

	number, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("стд::лог10: неможливо обчислити логарифм від нечислового значення '%s'", values[0]),
		)
	}

	return fmt.Sprintf("%f", math.Log10(number)), nil
}
