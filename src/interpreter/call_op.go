package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeCallOp(
	node *ast.CallOpNode, callable types.ValueType, rootDir, thisPackage, parentPackage string,
) (types.ValueType, error) {
	switch function := callable.(type) {
	case types.FunctionType:
		parametersLen := len(node.Parameters)
		argsLen := len(function.Arguments)
		if argsLen > 0 && function.Arguments[argsLen - 1].IsVariadic {
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
					if c < argsLen - 2 {
						parametersStr += ", "
					} else if c < argsLen - 1 {
						parametersStr += " та "
					}
				}

				return nil, util.RuntimeError(fmt.Sprintf(
					"при виклику '%s()' відсутн%s %d необхідн%s параметр%s: %s",
					function.Name, end1, diffLen, end2, end3, parametersStr,
				))
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

				return nil, util.RuntimeError(fmt.Sprintf(
					"'%s()' приймає %d необхідн%s параметр%s, отримано %d",
					function.Name, argsLen, end1, end2, parametersLen,
				))
			}
		}

		var args []types.ValueType
		kwargs := map[string]types.ValueType{}
		var c int
		for c = 0; c < argsLen; c++ {
			arg, _, err := i.executeNode(node.Parameters[c], rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			// TODO: remove
			if arg == nil {
				//arg = types.NilType{}
				panic("fatal: argument is nil")
			}

			if arg.TypeHash() == types.NilTypeHash {
				if function.Arguments[c].TypeHash != types.NilTypeHash && !function.Arguments[c].IsNullable {
					return nil, util.RuntimeError(fmt.Sprintf(
						"аргумент '%s' очікує ненульовий параметр, отримано '%s'",
						function.Arguments[c].Name, arg.String(),
					))
				}
			} else if function.Arguments[c].TypeHash != types.AnyTypeHash && arg.TypeHash() != function.Arguments[c].TypeHash {
				return nil, util.RuntimeError(fmt.Sprintf(
					"аргумент '%s' очікує параметр з типом '%s', отримано '%s'",
					function.Arguments[c].Name, function.Arguments[c].TypeName(), arg.TypeName(),
				))
			}

			kwargs[function.Arguments[c].Name] = arg
			args = append(args, arg)
		}

		if len(function.Arguments) > 0 {
			if lastArgument := function.Arguments[len(function.Arguments)-1]; lastArgument.IsVariadic {
				lastParameter := types.NewListType()
				if len(node.Parameters) - parametersLen > 0 {
					parametersLen = len(node.Parameters)
					for k := c; k < parametersLen; k++ {
						arg, _, err := i.executeNode(node.Parameters[k], rootDir, thisPackage, parentPackage)
						if err != nil {
							return nil, err
						}

						if arg.TypeHash() == types.NilTypeHash {
							if lastArgument.TypeHash != types.NilTypeHash && !lastArgument.IsNullable {
								return nil, util.RuntimeError(fmt.Sprintf(
									"аргумент '%s' очікує ненульовий параметр, отримано '%s'",
									lastArgument.Name, arg.String(),
								))
							}
						} else if lastArgument.TypeHash != types.AnyTypeHash && arg.TypeHash() != lastArgument.TypeHash {
							return nil, util.RuntimeError(fmt.Sprintf(
								"аргумент '%s' очікує список параметрів з типом '%s', отримано '%s'",
								lastArgument.Name, lastArgument.TypeName(), arg.TypeName(),
							))
						}

						lastParameter.Values = append(lastParameter.Values, arg)
						args = append(args, arg)
					}
				}

				kwargs[lastArgument.Name] = lastParameter
			}
		}

		res, err := function.Callable(args, kwargs)
		if err != nil {
			return nil, err
		}

		// TODO: remove
		if res == nil {
			res = types.NilType{}
			panic("fatal: returned value is nil")
		}

		if res.TypeHash() == types.NilTypeHash {
			if function.ReturnType.TypeHash != types.NilTypeHash && !function.ReturnType.IsNullable {
				return nil, util.RuntimeError(fmt.Sprintf(
					"'%s()' повертає ненульове значення, отримано '%s'",
					function.Name, res.String(),
				))
			}
		} else if function.ReturnType.TypeHash != types.AnyTypeHash && res.TypeHash() != function.ReturnType.TypeHash {
			return nil, util.RuntimeError(fmt.Sprintf(
				"'%s()' повертає значення типу '%s', отримано значення з типом '%s'",
				function.Name, function.ReturnType.String(), res.TypeName(),
			))
		}

		return res, nil
	default:
		return nil, util.RuntimeError(fmt.Sprintf(
			"неможливо застосувати оператор виклику до об'єкта '%s' з типом '%s'",
			node.CallableName.Text, callable.TypeName(),
		))
	}
}
