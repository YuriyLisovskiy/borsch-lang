package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (node *ClassDef) Evaluate(state State) (types.Object, error) {
	// TODO: add doc
	// cls := &types.Class{
	// 	Name:    node.Name.String(),
	// 	IsFinal: node.IsFinal,
	// 	Class:   nil,
	// 	Parent:  state.Package(),
	// }

	ctx := state.Context()
	var bases []*types.Class
	for _, name := range node.Bases {
		base, err := ctx.GetClass(name.String())
		if err != nil {
			return nil, err
		}

		baseClass := base.(*types.Class)
		if baseClass.IsFinal {
			return nil, state.RuntimeError(
				fmt.Sprintf("клас '%s' є закритим для розширення, не наслідуйте цей клас", name),
				node,
			)
		}

		bases = append(bases, baseClass)
	}

	cls := types.ObjectClass.ClassNew(node.Name.String(), map[string]types.Object{}, node.IsFinal, nil, nil)
	// cls.Bases = append(cls.Bases, bases...)
	if len(bases) > 0 {
		cls.Bases = bases
	}

	err := ctx.SetVar(node.Name.String(), cls)
	if err != nil {
		return nil, err
	}

	classContext := ctx.Derive()
	classContext.PushScope(map[string]types.Object{})
	for _, classMember := range node.Members {
		_, err := classMember.Evaluate(state.NewChild().WithContext(classContext), cls)
		if err != nil {
			return nil, err
		}
	}

	cls.Dict = classContext.PopScope()
	return cls, nil
}

func (node *ClassMember) Evaluate(state State, class *types.Class) (types.Object, error) {
	if node.Variable != nil {
		return node.Variable.Evaluate(state)
	}

	if node.Method != nil {
		return node.Method.Evaluate(
			state,
			state.Package().(*types.Package),
			true,
			func(arguments []types.MethodParameter, returnTypes []types.MethodReturnType) error {
				if err := checkMethod(class, arguments, returnTypes); err != nil {
					return state.RuntimeError(err.Error(), node)
				}

				if node.Method.Name == builtin.ConstructorName {
					err := checkConstructor(arguments, returnTypes)
					if err != nil {
						return state.RuntimeError(err.Error(), node)
					}
				}

				return nil
			},
		)
	}

	if node.Operator != nil {
		operator, err := node.Operator.Evaluate(
			state,
			func(parameters []types.MethodParameter, returnTypes []types.MethodReturnType, opName string) error {
				opHash := common.OperatorHashFromString(opName)
				paramsCount := getParamsCountOfOperator(opHash)
				if err := checkOperator(class, parameters, returnTypes, paramsCount, opHash); err != nil {
					return err
				}

				return nil
			},
		)
		if err != nil {
			return nil, err
		}

		class.Operators[common.OperatorHashFromString(operator.Name)] = operator
		return operator, nil
	}

	if node.Class != nil {
		return node.Class.Evaluate(state)
	}

	panic("unreachable")
}

func checkMethod(class *types.Class, args []types.MethodParameter, _ []types.MethodReturnType) error {
	if len(args) == 0 {
		// TODO: ukr error text!
		return errors.New("not enough args, self required")
	}

	if args[0].Class != class {
		return errors.New(
			fmt.Sprintf(
				"перший параметер методу має бути типу '%s' отримано '%s'", class.Name, args[0].Class.Name,
			),
		)
	}

	return nil
}

func checkConstructor(_ []types.MethodParameter, returnTypes []types.MethodReturnType) error {
	switch len(returnTypes) {
	case 0:
		// skip
	case 1:
		if returnTypes[0].Class != types.NilClass {
			return errors.New("конструктор має повертати 'нуль'")
		}
	default:
		return errors.New("ERROR")
	}

	return nil
}

// getParamsCountOfOperator returns a count of parameters minus 1,
// which can be present in the operator depending on whether it is
// the binary, unary or multi-parameter operator.
//
// Attention: -1 marks multi-parameter operator.
func getParamsCountOfOperator(op common.OperatorHash) int {
	if op.IsUnary() {
		return 0
	}

	if op.IsBinary() {
		return 1
	}

	return -1
}
