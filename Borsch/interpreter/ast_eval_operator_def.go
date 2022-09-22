package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (node *OperatorDef) Evaluate(
	state State,
	check func([]types.MethodParameter, []types.MethodReturnType) error,
) (*types.Method, error) {
	arguments, err := node.ParametersSet.Evaluate(state)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(state, node.ReturnTypes)
	if err != nil {
		return nil, err
	}

	if check != nil {
		if err := check(arguments, returnTypes); err != nil {
			return nil, err
		}
	}

	methodF := func(ctx types.Context, _ types.Tuple, kwargs types.StringDict) (types.Object, error) {
		return node.Body.Evaluate(state.NewChild().WithContext(ctx))
	}

	return types.MethodNew(node.Op, state.Package().(*types.Package), arguments, returnTypes, methodF), nil
}

func checkOperator(class *types.Class, args []types.MethodParameter, otherArgsCount int, op common.OperatorHash) error {
	argsLen := len(args)
	if argsLen == 0 {
		// TODO: ukr error text!
		return errors.New("not enough args, self required")
	}

	if args[0].Class != class {
		return errors.New(
			fmt.Sprintf(
				"перший параметер оператора %s має бути типу '%s' отримано '%s'",
				op.Sign(),
				class.Name,
				args[0].Class.Name,
			),
		)
	}

	if argsLen-1 != otherArgsCount {
		return errors.New(
			fmt.Sprintf(
				"кількість параметрів оператора %s має дорівнювати %d, крім першого, отримано '%d'",
				op.Sign(),
				otherArgsCount,
				argsLen-1,
			),
		)
	}

	return nil
}
