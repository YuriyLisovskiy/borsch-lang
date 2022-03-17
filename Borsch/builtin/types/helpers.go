package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

// Represent calls __представлення__ on the object or returns a sensible default.
func Represent(state common.State, self common.Value) (common.Value, error) {
	if I, ok := self.(I__represent__); ok {
		return I.__represent__()
	} else if res, ok, err := TypeCall0(state, self, common.RepresentOperatorName); ok {
		return res, err
	}

	return NewStringInstance(fmt.Sprintf("<%s instance at %p>", self.(ObjectInstance).GetClass().Name, self)), nil
}

// Str calls common.StringOperator on the object and if not found
// calls common.RepresentOperator.
func Str(state common.State, self common.Value) (common.Value, error) {
	if I, ok := self.(I__str__); ok {
		return I.__str__()
	} else if res, ok, err := TypeCall0(state, self, common.StringOperatorName); ok {
		return res, err
	}

	return Represent(state, self)
}

// StrAsString returns object as a string.
//
// Calls Str then makes sure the output is a string.
func StrAsString(state common.State, self common.Value) (string, error) {
	res, err := Str(state, self)
	if err != nil {
		return "", err
	}

	str, ok := res.(String)
	if !ok {
		return "", ErrorNewf(
			TypeError,
			"result of __str__ must be string, not '%s'",
			res.(ObjectInstance).GetClass().Name,
		)
	}

	return str.Value, nil
}
