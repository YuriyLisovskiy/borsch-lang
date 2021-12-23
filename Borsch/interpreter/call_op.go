package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ast"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (i *Interpreter) callFunction(
	node *ast.CallOpNode,
	function *types.FunctionInstance,
	args *[]types.Type,
	kwargs *map[string]types.Type,
	skipArgsCount int,
	rootDir string,
	thisPackage string,
	parentPackage string,
) (types.Type, error) {
	parametersLen := len(node.Parameters) + skipArgsCount
	argsLen := len(function.Arguments)
	if argsLen > 0 && function.Arguments[argsLen-1].IsVariadic {
		argsLen--
		if parametersLen > argsLen {
			parametersLen = argsLen
		}
	}

	if parametersLen != argsLen {
		diffLen := argsLen - parametersLen
		if diffLen > 0 {
			end1 := "ій"
			end2 := "ий"
			end3 := ""
			if diffLen > 1 && diffLen < 5 {
				end1 = "і"
				end2 = "і"
				end3 = "и"
			} else if diffLen != 1 {
				end1 = "і"
				end2 = "их"
				end3 = "ів"
			}

			parametersStr := ""
			for c := parametersLen; c < argsLen; c++ {
				parametersStr += fmt.Sprintf("'%s'", function.Arguments[c].Name)
				if c < argsLen-2 {
					parametersStr += ", "
				} else if c < argsLen-1 {
					parametersStr += " та "
				}
			}

			return nil, util.RuntimeError(
				fmt.Sprintf(
					"при виклику '%s()' відсутн%s %d необхідн%s параметр%s: %s",
					function.Name, end1, diffLen, end2, end3, parametersStr,
				),
			)
		} else {
			end1 := "ий"
			end2 := ""
			if argsLen > 1 && argsLen < 5 {
				end1 = "і"
				end2 = "и"
			} else if argsLen != 1 {
				end1 = "их"
				end2 = "ів"
			}

			return nil, util.RuntimeError(
				fmt.Sprintf(
					"'%s()' приймає %d необхідн%s параметр%s, отримано %d",
					function.Name, argsLen, end1, end2, parametersLen,
				),
			)
		}
	}

	var c int
	for c = 0; c < argsLen - skipArgsCount; c++ {
		arg, _, err := i.executeNode(node.Parameters[c], rootDir, thisPackage, parentPackage)
		if err != nil {
			return nil, err
		}

		// TODO: remove
		if arg == nil {
			// arg = types.NewNilInstance()
			panic("fatal: argument is nil")
		}

		if arg.GetTypeHash() == types.NilTypeHash {
			if function.Arguments[c + skipArgsCount].TypeHash != types.NilTypeHash && !function.Arguments[c + skipArgsCount].IsNullable {
				return nil, util.RuntimeError(
					fmt.Sprintf(
						"аргумент '%s' очікує ненульовий параметр, отримано '%s'",
						function.Arguments[c + skipArgsCount].Name, arg.String(),
					),
				)
			}
		} else if function.Arguments[c + skipArgsCount].TypeHash != types.AnyTypeHash && arg.GetTypeHash() != function.Arguments[c + skipArgsCount].TypeHash {
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"аргумент '%s' очікує параметр з типом '%s', отримано '%s'",
					function.Arguments[c + skipArgsCount].Name, function.Arguments[c + skipArgsCount].TypeName(), arg.GetTypeName(),
				),
			)
		}

		(*kwargs)[function.Arguments[c + skipArgsCount].Name] = arg
		*args = append(*args, arg)
	}

	if len(function.Arguments) - skipArgsCount > 0 {
		if lastArgument := function.Arguments[len(function.Arguments)-1]; lastArgument.IsVariadic {
			lastParameter := types.NewListInstance()
			if len(node.Parameters) + skipArgsCount - parametersLen > 0 {
				parametersLen = len(node.Parameters)
				for k := c; k < parametersLen; k++ {
					arg, _, err := i.executeNode(node.Parameters[k], rootDir, thisPackage, parentPackage)
					if err != nil {
						return nil, err
					}

					if arg.GetTypeHash() == types.NilTypeHash {
						if lastArgument.TypeHash != types.NilTypeHash && !lastArgument.IsNullable {
							return nil, util.RuntimeError(
								fmt.Sprintf(
									"аргумент '%s' очікує ненульовий параметр, отримано '%s'",
									lastArgument.Name, arg.String(),
								),
							)
						}
					} else if lastArgument.TypeHash != types.AnyTypeHash && arg.GetTypeHash() != lastArgument.TypeHash {
						return nil, util.RuntimeError(
							fmt.Sprintf(
								"аргумент '%s' очікує список параметрів з типом '%s', отримано '%s'",
								lastArgument.Name, lastArgument.TypeName(), arg.GetTypeName(),
							),
						)
					}

					lastParameter.Values = append(lastParameter.Values, arg)
					*args = append(*args, arg)
				}
			}

			(*kwargs)[lastArgument.Name] = lastParameter
		}
	}

	res, err := function.Call(args, kwargs)
	if err != nil {
		return nil, err
	}

	if res == nil {
		res = types.NewNilInstance()
		// TODO: remove
		// panic("fatal: returned value is nil")
	}

	if res.GetTypeHash() == types.NilTypeHash {
		if function.ReturnType.TypeHash != types.NilTypeHash && !function.ReturnType.IsNullable {
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"'%s()' повертає ненульове значення, отримано '%s'",
					function.Name, res.String(),
				),
			)
		}
	} else if function.ReturnType.TypeHash != types.AnyTypeHash && res.GetTypeHash() != function.ReturnType.TypeHash {
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"'%s()' повертає значення типу '%s', отримано значення з типом '%s'",
				function.Name, function.ReturnType.String(), res.GetTypeName(),
			),
		)
	}

	return res, nil
}

func (i *Interpreter) executeCallOp(
	node *ast.CallOpNode, object types.Type, rootDir, thisPackage, parentPackage string,
) (types.Type, error) {
	switch callable := object.(type) {
	case *types.FunctionInstance:
		return i.callFunction(
			node,
			callable,
			&[]types.Type{},
			&map[string]types.Type{},
			0,
			rootDir,
			thisPackage,
			parentPackage,
		)
	case *types.Class:
		constructorAttribute, err := callable.GetAttribute(ops.ConstructorName)
		if err == nil {
			switch constructor := constructorAttribute.(type) {
			case *types.FunctionInstance:
				instance, err := callable.GetEmptyInstance()
				if err != nil {
					return nil, err
				}

				args := []types.Type{instance}
				kwargs := map[string]types.Type{"я": instance}
				_, err = i.callFunction(
					node,
					constructor,
					&args,
					&kwargs,
					1,
					rootDir,
					thisPackage,
					parentPackage,
				)
				if err != nil {
					return nil, err
				}

				return args[0], nil
			}
		}
	case types.Instance:
		callOperatorAttribute, err := callable.GetClass().GetAttribute(ops.CallOperatorName)
		if err == nil {
			switch callOperator := callOperatorAttribute.(type) {
			case *types.FunctionInstance:
				args := []types.Type{object}
				kwargs := map[string]types.Type{"я": object}
				return i.callFunction(
					node,
					callOperator,
					&args,
					&kwargs,
					1,
					rootDir,
					thisPackage,
					parentPackage,
				)
			}
		}
	}

	return nil, util.ObjectIsNotCallable(node.CallableName.Text, object.GetTypeName())
}
