package interpreter

import (
	"fmt"
	"reflect"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

func evalBinaryOperator(
	state State,
	valueToSet types.Object,
	operator func(ctx types.Context, a, b types.Object) (types.Object, error),
	current OperatorEvaluatable,
	next OperatorEvaluatable,
) (types.Object, error) {
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
	operator func(ctx types.Context, a types.Object) (types.Object, error),
	next OperatorEvaluatable,
) (types.Object, error) {
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
	variable types.Object,
	ranges_ []*Range,
	valueToSet types.Object,
) (types.Object, error) {
	switch iterable := variable.(type) {
	case types.ISequence:
		errMsg := ""
		if ranges_[0].IsSlicing {
			errMsg = "ліва межа має бути цілого типу"
		} else {
			errMsg = "індекс має бути цілого типу"
		}

		leftIdx, err := mustInt(
			state, ranges_[0].LeftBound, func(t types.Object) error {
				return types.NewTypeErrorf("%s, отримано %s", errMsg, t.Class().Name)
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

		var element types.Object
		if ranges_[0].RightBound != nil {
			rightIdx, err := mustInt(
				state, ranges_[0].RightBound, func(t types.Object) error {
					return types.NewTypeErrorf("права межа має бути цілого типу, отримано %s", t.Class().Name)
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

		return nil, types.NewTypeErrorf(
			"неможливо застосувати оператор %s до об'єкта з типом '%s'",
			operatorDescription, variable.Class().Name,
		)
	}
}

func mustInt(state State, expression *Expression, errFunc func(types.Object) error) (types.Int, error) {
	value, err := expression.Evaluate(state, nil)
	if err != nil {
		return 0, err
	}

	switch integer := value.(type) {
	case types.Int:
		return integer, nil
	default:
		return 0, errFunc(value)
	}
}

func checkValuesCountToUnpack(lhsLen, rhsLen int64) error {
	if lhsLen > rhsLen {
		return types.NewValueErrorf(
			"недостатньо значень для розпакування (потрібно %d, отримано %d)",
			lhsLen,
			rhsLen,
		)
	}

	if lhsLen < rhsLen {
		return types.NewValueErrorf("занадто багато значень для розпакування (потрібно %d)", lhsLen)
	}

	return nil
}

func evalRhs(state State, _ /*lhs*/ []*Expression, rhs []*Expression) (*types.Tuple, error) {
	rhsSeq := types.Tuple{}
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

func unpack(state State, lhs []*Expression, rhs []*Expression) (types.Object, error) {
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

	// lhsLen := len(lhs)
	// rhsLen := len(rhs)
	// if lhsLen > rhsLen {
	// 	return unpackList(state, lhs, rhs[0])
	// }
	//
	// sequence, result, err := getSequenceOrResult(state, lhs, rhs)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if result != nil {
	// 	return result, err
	// }
	//
	// if lhsLen > len(sequence) {
	// 	// TODO: return unable to unpack
	// 	panic(fmt.Sprintf("unable to unpack %d elements to %d vars", len(sequence), lhsLen))
	// } else if lhsLen < len(sequence) {
	// 	// TODO: return unable to unpack
	// 	panic(fmt.Sprintf("unable to unpack %d elements to %d vars", len(sequence), lhsLen))
	// }
	//
	// var i int
	// list := types.NewList()
	// for i = 0; i < lhsLen-1; i++ {
	// 	element, err := lhs[i].Evaluate(state, sequence[i])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	list.Values = append(list.Values, element)
	// }
	//
	// if i < len(sequence)-1 {
	// 	rest := types.NewList()
	// 	rest.Values = sequence[i:]
	// 	list.Values = append(list.Values, rest)
	// } else {
	// 	element, err := lhs[i].Evaluate(state, sequence[i])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	list.Values = append(list.Values, element)
	// }
	//
	// return list, nil
}

// unpackFromSequence unpacks right-hand tuple into left-hand variables:
//  а, б, в = 1, 2.0, "Привіт"
//  а, б, в = [1, 2.0, "Привіт"]
//  а, б, в = (1, 2.0, "Привіт")
//  а, б, в = <any sequence>
func unpackFromSequence(state State, dest []*Expression, src types.Object, skipSrcCheck bool) (types.Object, error) {
	var sequence types.ISequence
	if skipSrcCheck {
		sequence = src.(types.ISequence)
	} else {
		var ok bool
		sequence, ok = src.(types.ISequence)
		if !ok {
			return nil, types.NewErrorf(
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

	result := types.NewList()
	ctx := state.Context()
	for i, expr := range dest {
		item, _ := sequence.GetElement(ctx, types.Int(i))
		obj, err := expr.Evaluate(state, item)
		if err != nil {
			return nil, err
		}

		result.Values = append(result.Values, obj)
	}

	return result, nil
}

func getSequenceOrResult(state State, lhs []*Expression, rhs []*Expression) (
	[]types.Object,
	types.Object,
	error,
) {
	rhsLen := len(rhs)
	var sequence []types.Object
	if rhsLen == 1 {
		element, err := rhs[0].Evaluate(state, nil)
		if err != nil {
			return nil, nil, err
		}

		switch list := element.(type) {
		case *types.List:
			if len(lhs) == 1 {
				result, err := lhs[0].Evaluate(state, list)
				if err != nil {
					return nil, nil, err
				}

				return nil, result, nil
			}

			sequence = list.Values
		default:
			sequence = append(sequence, element)
		}
	} else {
		for _, expr := range rhs {
			element, err := expr.Evaluate(state, nil)
			if err != nil {
				return nil, nil, err
			}

			sequence = append(sequence, element)
		}
	}

	return sequence, nil, nil
}

func unpackList(state State, lhs []*Expression, rhs *Expression) (types.Object, error) {
	element, err := rhs.Evaluate(state, nil)
	if err != nil {
		return nil, err
	}

	switch list := element.(type) {
	case types.ISequence:
		lhsLen := types.Int(len(lhs))
		rhsLen, err := list.Length(state.Context())
		if err != nil {
			return nil, err
		}

		if lhsLen > rhsLen {
			// TODO: return error
			panic(fmt.Sprintf("unable to unpack %d elements of %s to %d vars", rhsLen, element.Class().Name, lhsLen))
		}

		ctx := state.Context()
		var i types.Int
		resultList := types.NewList()
		for i = 0; i < lhsLen-1; i++ {
			value, err := list.GetElement(ctx, i)
			if err != nil {
				return nil, err
			}

			item, err := lhs[i].Evaluate(state, value)
			if err != nil {
				return nil, err
			}

			resultList.Values = append(resultList.Values, item)
		}

		listLen, err := list.Length(state.Context())
		if i < listLen-1 {
			rest, err := list.Slice(ctx, i, listLen)
			if err != nil {
				return nil, err
			}

			resultList.Values = append(resultList.Values, rest)
		} else {
			value, err := list.GetElement(ctx, i)
			if err != nil {
				return nil, err
			}

			element, err := lhs[i].Evaluate(state, value)
			if err != nil {
				return nil, err
			}

			resultList.Values = append(resultList.Values, element)
		}

		return resultList, nil
	}

	// TODO: return error
	return nil, types.NewErrorf("unable to unpack %s", element.Class().Name)
}

func evalReturnTypes(state State, returnTypes []*ReturnType) ([]types.MethodReturnType, error) {
	var result []types.MethodReturnType
	if len(returnTypes) == 0 {
		result = append(
			result, types.MethodReturnType{
				Class:      types.NilClass,
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

func getCurrentValue(ctx types.Context, prevValue types.Object, ident string) (types.Object, error) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return types.GetAttribute(ctx, prevValue, ident)
	}

	return ctx.GetVar(ident)
}

func setCurrentValue(ctx types.Context, prevValue types.Object, ident string, valueToSet types.Object) (
	types.Object,
	error,
) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return prevValue, types.SetAttribute(ctx, prevValue, ident, valueToSet)
	}

	return valueToSet, ctx.SetVar(ident, valueToSet)
}

func checkForNilAttribute(ident string) error {
	switch ident {
	case "нуль", "нульове":
		return types.NewAttributeErrorf("'%s' не є атрибутом", ident)
	}

	return nil
}

func updateArgs(state State, arguments []*Expression, args *types.Tuple) error {
	for _, expressionArgument := range arguments {
		arg, err := expressionArgument.Evaluate(state, nil)
		if err != nil {
			return err
		}

		*args = append(*args, arg)
	}

	return nil
}
