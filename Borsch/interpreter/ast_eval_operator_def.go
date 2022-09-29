package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (node *OperatorDef) Evaluate(
	state State,
	check func(_ []types.MethodParameter, _ []types.MethodReturnType, opName string) error,
) (*types.Method, error) {
	// Assume that node.Op is a valid operator.
	arguments, err := node.ParametersSet.Evaluate(state)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(state, node.ReturnTypes)
	if err != nil {
		return nil, err
	}

	switch node.Op {
	case "+", "-":
		if len(arguments) == 1 {
			node.Op = "_" + node.Op
		}
	}

	if check != nil {
		if err := check(arguments, returnTypes, node.Op); err != nil {
			return nil, err
		}
	}

	methodF := func(ctx types.Context, _ types.Tuple, kwargs types.StringDict) (types.Object, error) {
		return node.Body.Evaluate(state.NewChild().WithContext(ctx))
	}

	return types.MethodNew(node.Op, state.Package().(*types.Package), arguments, returnTypes, methodF), nil
}

func checkOperator(
	class *types.Class,
	params []types.MethodParameter,
	returnTypes []types.MethodReturnType,
	expectedParamsCount int,
	opHash common.OperatorHash,
) error {
	argsLen := len(params)
	if argsLen == 0 {
		// TODO: ukr error text!
		return errors.New("not enough args, self required")
	}

	if params[0].Class != class {
		return errors.New(
			fmt.Sprintf(
				"перший параметр оператора %s має бути типу '%s' отримано '%s'",
				opHash.Sign(),
				class.Name,
				params[0].Class.Name,
			),
		)
	}

	if expectedParamsCount != -1 && argsLen-1 != expectedParamsCount {
		return errors.New(
			fmt.Sprintf(
				"кількість параметрів оператора %s має дорівнювати %d, крім першого, отримано '%d'",
				opHash.Sign(),
				expectedParamsCount,
				argsLen-1,
			),
		)
	}

	// check the return type(s)
	switch opHash {
	case common.LengthOp, common.IntOp:
		return checkSingleReturnType(returnTypes, types.IntClass, opHash)
	case common.BoolOp:
		return checkSingleReturnType(returnTypes, types.BoolClass, opHash)
	case common.StringOp, common.RepresentationOp:
		return checkSingleReturnType(returnTypes, types.StringClass, opHash)
	}

	return nil
}

func checkSingleReturnType(
	retTypes []types.MethodReturnType,
	expectedClass *types.Class,
	opHash common.OperatorHash,
) error {
	if len(retTypes) != 1 {
		return types.NewErrorf("оператор ʼ%sʼ має повертати єдине значення", opHash.Name())
	}

	if retTypes[0].Class != expectedClass {
		return types.NewTypeErrorf(
			"тип, значення якого повертає оператор ʼ%sʼ, має бути ʼ%sʼ, отримано ʼ%sʼ",
			opHash.Name(),
			expectedClass.Name,
			retTypes[0].Class.Name,
		)
	}

	return nil
}
