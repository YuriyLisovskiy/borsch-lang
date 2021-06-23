package interpreter

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/ast"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func (i *Interpreter) executeCallOp(
	node *ast.CallOpNode, rootDir, thisPackage, parentPackage string,
) (types.ValueType, error) {
	attr, err := i.getVar(thisPackage, node.CallableName.Text)
	if err != nil {
		return nil, err
	}

	switch function := attr.(type) {
	case types.FunctionType:
		argsLen := len(node.Args)
		parametersLen := len(function.Parameters)
		if parametersLen > 0 && function.Parameters[parametersLen - 1].IsVariadic {
			parametersLen--
			if argsLen > parametersLen {
				argsLen = parametersLen
			}
		}

		if argsLen != parametersLen {
			diffLen := parametersLen - argsLen
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
				for c := argsLen; c < parametersLen; c++ {
					parametersStr += fmt.Sprintf("'%s'", function.Parameters[c].Name)
					if c < parametersLen - 2 {
						parametersStr += ", "
					} else if c < parametersLen - 1 {
						parametersStr += " та "
					}
				}

				return nil, util.RuntimeError(fmt.Sprintf(
					"при виклику '%s()' відсутн%s %d необхідн%s аргумент%s: %s",
					function.Name, end1, diffLen, end2, end3, parametersStr,
				))
			} else {
				end1 := ""
				if parametersLen > 1 && parametersLen < 5 {
					end1 = "и"
				} else if parametersLen != 1 {
					end1 = "ів"
				}

				return nil, util.RuntimeError(fmt.Sprintf(
					"%s() приймає %d необхідн~ аргумент%s, отримано %d",
					function.Name, parametersLen, end1, argsLen,
				))
			}
		}

		var args []types.ValueType
		kwargs := map[string]types.ValueType{}
		var c int
		for c = 0; c < parametersLen; c++ {
			arg, err := i.executeNode(node.Args[c], rootDir, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}

			if function.Parameters[c].TypeHash != types.AnyTypeHash && arg.TypeHash() != function.Parameters[c].TypeHash {
				return nil, util.RuntimeError(fmt.Sprintf(
					"параметр '%s' очікує аргумент з типом '%s', отримано '%s'",
					function.Parameters[c].Name, function.Parameters[c].TypeName(), arg.TypeName(),
				))
			}

			kwargs[function.Parameters[c].Name] = arg
			args = append(args, arg)
		}

		if len(node.Args) - argsLen > 0 {
			argsLen = len(node.Args)
			lastParameter := function.Parameters[parametersLen]
			for k := c; k < argsLen; k++ {
				arg, err := i.executeNode(node.Args[k], rootDir, thisPackage, parentPackage)
				if err != nil {
					return nil, err
				}

				if lastParameter.TypeHash != types.AnyTypeHash && arg.TypeHash() != lastParameter.TypeHash {
					return nil, util.RuntimeError(fmt.Sprintf(
						"параметр '%s' очікує аргумент з типом '%s', отримано '%s'",
						lastParameter.Name, lastParameter.TypeName(), arg.TypeName(),
					))
				}

				if _, ok := kwargs[lastParameter.Name]; !ok {
					kwargs[lastParameter.Name] = arg
				}

				args = append(args, arg)
			}
		}

		var res types.ValueType
		if function.Callable != nil {
			res, err = function.Callable(args, kwargs)
			if err != nil {
				return nil, err
			}
		} else {
			res, err = i.executeBlock(kwargs, function.Code, thisPackage, parentPackage)
			if err != nil {
				return nil, err
			}
		}

		if res != nil && res.TypeHash() != function.ReturnType {
			return nil, util.RuntimeError(fmt.Sprintf(
				"неможливо представити тип '%s', який повернув оператор виклику, типом '%s'",
				res.TypeName(), types.GetTypeName(function.ReturnType),
			))
		}

		return res, nil
	default:
		return nil, util.RuntimeError(fmt.Sprintf(
			"неможливо застосувати оператор виклику до об'єкта '%s' з типом '%s'",
			node.CallableName.Text, attr.TypeName(),
		))
	}
}
