package types

import "testing"

func Test_checkNumberOfArgs_TooFewArgs_SameMinMax(t *testing.T) {
	result := checkNumberOfArgs("", 3, 4, 4, 4)
	if result == nil {
		t.Error("should return error")
	}
}

func Test_checkNumberOfArgs_TooMuchArgs_SameMinMax(t *testing.T) {
	result := checkNumberOfArgs("", 4, 3, 3, 3)
	if result == nil {
		t.Error("should return error")
	}
}

func Test_checkNumberOfArgs_TooFewArgs_SmallerThanMin(t *testing.T) {
	result := checkNumberOfArgs("", 1, 4, 2, 4)
	if result == nil {
		t.Error("should return error")
	}
}

func Test_checkNumberOfArgs_TooMuchArgs_GreaterThanMax(t *testing.T) {
	result := checkNumberOfArgs("", 4, 3, 2, 3)
	if result == nil {
		t.Error("should return error")
	}
}
