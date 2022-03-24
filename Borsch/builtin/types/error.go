package types

import (
	"errors"
	"fmt"
)

func ErrorNewf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf("ПомилкаВиконання: %s", fmt.Sprintf(format, args...)))
}
