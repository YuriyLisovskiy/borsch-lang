package interpreter

import (
	"reflect"

	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

func evalBinaryOperator(
	state State,
	valueToSet types2.Object,
	operator func(ctx types2.Context, a, b types2.Object) (types2.Object, error),
	current OperatorEvaluatable,
	next OperatorEvaluatable,
) (types2.Object, error) {
	left, err := current.Evaluate(state, valueToSet)
	if err != nil {
		return nil, err
	}

	if !reflect.ValueOf(next).IsNil() {
		right, err := next.Evaluate(state, valueToSet)
		if err != nil {
			return nil, err
		}

		return operator(state.Context(), left, right)
	}

	return left, nil
}

func evalUnaryOperator(
	state State,
	operator func(ctx types2.Context, a types2.Object) (types2.Object, error),
	next OperatorEvaluatable,
) (types2.Object, error) {
	if next != nil {
		value, err := next.Evaluate(state, nil)
		if err != nil {
			return nil, err
		}

		return operator(state.Context(), value)
	}

	panic("unreachable")
}

// evalSlicingOperation: "ranges_" len should be greater than 0
func evalSlicingOperation(
	state State,
	variable types2.Object,
	ranges_ []*Range,
	valueToSet types2.Object,
) (types2.Object, error) {
	switch iterable := variable.(type) {
	case types2.ISequence:
		errMsg := ""
		if ranges_[0].IsSlicing {
			errMsg = "ліва межа має бути цілого типу"
		} else {
			errMsg = "індекс має бути цілого типу"
		}

		leftIdx, err := mustInt(
			state, ranges_[0].LeftBound, func(t types2.Object) error {
				return types2.NewTypeErrorf("%s, отримано %s", errMsg, t.Class().Name)
			},
		)
		if err != nil {
			return nil, err
		}

		ctx := state.Context()
		length, err := iterable.Length(ctx)
		if err != nil {
			return nil, err
		}

		if leftIdx < 0 {
			leftIdx = length + leftIdx
		}

		var element types2.Object
		if ranges_[0].RightBound != nil {
			rightIdx, err := mustInt(
				state, ranges_[0].RightBound, func(t types2.Object) error {
					return types2.NewTypeErrorf("права межа має бути цілого типу, отримано %s", t.Class().Name)
				},
			)
			if err != nil {
				return nil, err
			}

			if rightIdx < 0 {
				rightIdx = length + rightIdx
			}

			element, err = iterable.Slice(ctx, leftIdx, rightIdx)
			if err != nil {
				return nil, err
			}

			if len(ranges_) == 1 {
				// valueToSet is ignored, return error maybe.
				return element, nil
			}
		} else if ranges_[0].IsSlicing {
			element, err = iterable.Slice(ctx, leftIdx, length)
			if err != nil {
				return nil, err
			}

			if len(ranges_) == 1 {
				// valueToSet is ignored, return error maybe.
				return element, nil
			}
		} else {
			if len(ranges_) == 1 {
				if valueToSet != nil {
					return iterable.SetElement(ctx, leftIdx, valueToSet)
				}

				return iterable.GetElement(ctx, leftIdx)
			}

			element, err = iterable.GetElement(ctx, leftIdx)
			if err != nil {
				return nil, err
			}
		}

		element, err = evalSlicingOperation(state, element, ranges_[1:], valueToSet)
		if err != nil {
			return nil, err
		}

		return iterable.SetElement(ctx, leftIdx, element)
	default:
		operatorDescription := ""
		if ranges_[0].IsSlicing {
			operatorDescription = "зрізу"
		} else {
			operatorDescription = "довільного доступу"
		}

		return nil, types2.NewTypeErrorf(
			"неможливо застосувати оператор %s до об'єкта з типом '%s'",
			operatorDescription, variable.Class().Name,
		)
	}
}

func mustInt(state State, expression *Expression, errFunc func(types2.Object) error) (types2.Int, error) {
	value, err := expression.Evaluate(state, nil)
	if err != nil {
		return 0, err
	}

	switch integer := value.(type) {
	case types2.Int:
		return integer, nil
	default:
		return 0, errFunc(value)
	}
}

func checkValuesCountToUnpack(lhsLen, rhsLen int64) error {
	if lhsLen > rhsLen {
		return types2.NewValueErrorf(
			"недостатньо значень для розпакування (потрібно %d, отримано %d)",
			lhsLen,
			rhsLen,
		)
	}

	if lhsLen < rhsLen {
		return types2.NewValueErrorf("занадто багато значень для розпакування (потрібно %d)", lhsLen)
	}

	return nil
}

func evalRhs(state State, _ /*lhs*/ []*Expression, rhs []*Expression) (*types2.Tuple, error) {
	rhsSeq := types2.Tuple{}
	for _, expr := range rhs {
		obj, err := expr.Evaluate(state, nil)
		if err != nil {
			return nil, err
		}

		// TODO: possible optimization: validate types of lhs[i] and
		//  evaluated i-th expression.
		//  Note: lhs can contain only one expression when unpacking a
		//  statement like: к = 1, 2.0, "Привіт"
		rhsSeq = append(rhsSeq, obj)
	}

	return &rhsSeq, nil
}

func unpack(state State, lhs []*Expression, rhs []*Expression) (types2.Object, error) {
	lhsLen := len(lhs)
	rhsLen := len(rhs)
	if lhsLen == 1 && lhsLen < rhsLen {
		// unpack right-hand values into left-hand single variable:
		//  к = 1, 2.0, "Привіт";
		rhsSeq, err := evalRhs(state, lhs, rhs)
		if err != nil {
			return nil, err
		}

		return lhs[0].Evaluate(state, rhsSeq)
	}

	if rhsLen == 1 && lhsLen > rhsLen {
		// unpack right-hand tuple into left-hand variables:
		//  а, б, в = [1, 2.0, "Привіт"];
		obj, err := rhs[0].Evaluate(state, nil)
		if err != nil {
			return nil, err
		}

		return unpackFromSequence(state, lhs, obj, false)
	}

	if err := checkValuesCountToUnpack(int64(lhsLen), int64(rhsLen)); err != nil {
		return nil, err
	}

	rhsSeq, err := evalRhs(state, lhs, rhs)
	if err != nil {
		return nil, err
	}

	return unpackFromSequence(state, lhs, rhsSeq, true)
}

// unpackFromSequence unpacks right-hand tuple into left-hand variables:
//
//	а, б, в = 1, 2.0, "Привіт"
//	а, б, в = [1, 2.0, "Привіт"]
//	а, б, в = (1, 2.0, "Привіт")
//	а, б, в = <any sequence>
func unpackFromSequence(state State, dest []*Expression, src types2.Object, skipSrcCheck bool) (types2.Object, error) {
	var sequence types2.ISequence
	if skipSrcCheck {
		sequence = src.(types2.ISequence)
	} else {
		var ok bool
		sequence, ok = src.(types2.ISequence)
		if !ok {
			return nil, types2.NewErrorf(
				"неможливо розпакувати об'єкт з типом '%s', оскільки він не є послідовністю",
				src.Class().Name,
			)
		}

		seqLen, err := sequence.Length(state.Context())
		if err != nil {
			return nil, err
		}

		if err := checkValuesCountToUnpack(int64(len(dest)), int64(seqLen)); err != nil {
			return nil, err
		}
	}

	result := types2.NewList()
	ctx := state.Context()
	for i, expr := range dest {
		item, _ := sequence.GetElement(ctx, types2.Int(i))
		obj, err := expr.Evaluate(state, item)
		if err != nil {
			return nil, err
		}

		result.Values = append(result.Values, obj)
	}

	return result, nil
}

func evalReturnTypes(state State, returnTypes []*ReturnType) ([]types2.MethodReturnType, error) {
	var result []types2.MethodReturnType
	if len(returnTypes) == 0 {
		result = append(
			result, types2.MethodReturnType{
				Class:      types2.NilClass,
				IsNullable: false,
			},
		)
	} else {
		for _, returnType := range returnTypes {
			r, err := returnType.Evaluate(state.Context())
			if err != nil {
				return nil, err
			}

			// TODO: check if return type is class!
			result = append(result, *r)
		}
	}

	return result, nil
}

func getCurrentValue(ctx types2.Context, prevValue types2.Object, ident string) (types2.Object, error) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return types2.GetAttribute(ctx, prevValue, ident)
	}

	return ctx.GetVar(ident)
}

func setCurrentValue(ctx types2.Context, prevValue types2.Object, ident string, valueToSet types2.Object) (
	types2.Object,
	error,
) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return prevValue, types2.SetAttribute(ctx, prevValue, ident, valueToSet)
	}

	return valueToSet, ctx.SetVar(ident, valueToSet)
}

func checkForNilAttribute(ident string) error {
	switch ident {
	case "нуль", "нульове":
		return types2.NewAttributeErrorf("'%s' не є атрибутом", ident)
	}

	return nil
}

func updateArgs(state State, arguments []*Expression, args *types2.Tuple) error {
	for _, expressionArgument := range arguments {
		arg, err := expressionArgument.Evaluate(state, nil)
		if err != nil {
			return err
		}

		*args = append(*args, arg)
	}

	return nil
}
