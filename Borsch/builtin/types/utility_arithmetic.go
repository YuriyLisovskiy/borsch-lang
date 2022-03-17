// Automatically generated - DO NOT EDIT
// Regenerate with: go generate

// Arithmetic operations

package types

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/common"

// Add two objects together returning an Object.
//
// Will raise TypeError if add can't be run on these objects.
func Add(ctx common.Context, a, b common.Value) (common.Value, error) {
	// Try using a to add
	if A, ok := a.(I__add__); ok {
		res, err := A.__add__(ctx, b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_add if different in type to a
	if a.(ObjectInstance).GetClass() != b.(ObjectInstance).GetClass() {
		if B, ok := b.(I__reversed_add__); ok {
			res, err := B.__reversed_add__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(
		TypeError,
		"непідтримувані типи операндів для +: '%s' та '%s'",
		a.(ObjectInstance).GetClass().Name,
		b.(ObjectInstance).GetClass().Name,
	)
}

func InPlaceAdd(ctx common.Context, a, b common.Value) (common.Value, error) {
	if A, ok := a.(I__in_place_add__); ok {
		res, err := A.__in_place_add__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return Add(ctx, a, b)
}
