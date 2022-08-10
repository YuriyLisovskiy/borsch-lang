package types

import (
	"errors"
	"fmt"
)

func ErrorNewf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf("ПомилкаВиконання: %s", fmt.Sprintf(format, args...)))
}

func AssertionErrorNewf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf("ПомилкаПрипущення: %s", fmt.Sprintf(format, args...)))
}

func ZeroDivisionErrorNewf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf("ПомилкаДіленняНаНуль: %s", fmt.Sprintf(format, args...)))
}

func OperatorNotSupportedErrorNew(operator, lType, rType string) error {
	return errors.New(
		fmt.Sprintf(
			"ПомилкаТипу: екземпляри типів '%s' і '%s' не підтримують оператор '%s'",
			lType,
			rType,
			operator,
		),
	)
}
