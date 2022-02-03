package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func CheckFunctionArguments(
	state common.State,
	function *FunctionInstance,
	args *[]common.Type,
	_ *map[string]common.Type,
) error {
	parametersLen := len(*args)
	argsLen := len(function.Parameters)
	if argsLen > 0 && function.Parameters[argsLen-1].IsVariadic {
		argsLen--
		if parametersLen > argsLen {
			parametersLen = argsLen
		}
	}

	if parametersLen != argsLen {
		return makeArgumentError(argsLen, parametersLen, function.Parameters, function.Name)
	}

	var c int
	for c = 0; c < argsLen; c++ {
		parameter := function.Parameters[c]
		if parameter.Type == Any {
			continue
		}

		arg := (*args)[c]
		argPrototype := arg.(ObjectInstance).GetPrototype()
		if argPrototype == Nil && parameter.IsNullable {
			continue
		}

		if parameter.Type == argPrototype {
			continue
		}

		return util.RuntimeError(
			fmt.Sprintf(
				"аргумент '%s' очікує параметр з типом '%s', отримано '%s'",
				parameter.Name, parameter.GetTypeName(), arg.GetTypeName(),
			),
		)

		// if argPrototype == Nil {
		// 	if function.Parameters[c].Type != Nil && !function.Parameters[c].IsNullable {
		// 		argStr, err := arg.String(state)
		// 		if err != nil {
		// 			return err
		// 		}
		//
		// 		return util.RuntimeError(
		// 			fmt.Sprintf(
		// 				"аргумент '%s' очікує ненульовий параметр, отримано '%s'",
		// 				function.Parameters[c].Name,
		// 				argStr,
		// 			),
		// 		)
		// 	}
		// } else if function.Parameters[c].Type != Any && argPrototype != function.Parameters[c].Type {
		// 	return util.RuntimeError(
		// 		fmt.Sprintf(
		// 			"аргумент '%s' очікує параметр з типом '%s', отримано '%s'",
		// 			function.Parameters[c].Name, function.Parameters[c].GetTypeName(), arg.GetTypeName(),
		// 		),
		// 	)
		// }
	}

	if len(function.Parameters) > 0 {
		if lastArgument := function.Parameters[len(function.Parameters)-1]; lastArgument.IsVariadic {
			if len(*args)-parametersLen > 0 {
				parametersLen = len(*args)
				for k := c; k < parametersLen; k++ {
					arg := (*args)[k]
					argPrototype := arg.(ObjectInstance).GetPrototype()
					if argPrototype == Nil {
						if lastArgument.Type != Nil && !lastArgument.IsNullable {
							argStr, err := arg.String(state)
							if err != nil {
								return err
							}

							return util.RuntimeError(
								fmt.Sprintf(
									"аргумент '%s' очікує ненульовий параметр, отримано '%s'",
									lastArgument.Name,
									argStr,
								),
							)
						}
					} else if lastArgument.Type != nil && argPrototype != lastArgument.Type {
						return util.RuntimeError(
							fmt.Sprintf(
								"аргумент '%s' очікує список параметрів з типом '%s', отримано '%s'",
								lastArgument.Name,
								lastArgument.GetTypeName(),
								argPrototype.GetTypeName(),
							),
						)
					}
				}
			}
		}
	}

	return nil
}

func CheckResult(state common.State, result common.Type, function *FunctionInstance) error {
	if len(function.ReturnTypes) == 1 {
		err := checkSingleResult(state, result, function.ReturnTypes[0], function.Name)
		if err != nil {
			return errors.New(fmt.Sprintf(err.Error(), ""))
		}

		return nil
	}

	switch value := result.(type) {
	case ListInstance:
		if int64(len(function.ReturnTypes)) != value.Length(state) {
			var expectedTypes []string
			for _, retType := range function.ReturnTypes {
				expectedTypes = append(expectedTypes, retType.String())
			}

			var typesGot []string
			for _, retType := range value.Values {
				typesGot = append(typesGot, retType.GetTypeName())
			}

			return util.RuntimeError(
				fmt.Sprintf(
					"'%s' повертає значення з типами (%s), отримано (%s)",
					function.Name,
					strings.Join(expectedTypes, ", "),
					strings.Join(typesGot, ", "),
				),
			)
		}

		// TODO: check values in list
		for i, returnType := range function.ReturnTypes {
			if err := checkSingleResult(state, value.Values[i], returnType, function.Name); err != nil {
				return errors.New(fmt.Sprintf(err.Error(), fmt.Sprintf(" на позиції %d", i+1)))
			}
		}
	default:
		var expectedTypes []string
		for _, retType := range function.ReturnTypes {
			expectedTypes = append(expectedTypes, retType.String())
		}

		return util.RuntimeError(
			fmt.Sprintf(
				"'%s()' повертає значення з типами '(%s)', отримано '%s'",
				function.Name,
				strings.Join(expectedTypes, ", "),
				value.GetTypeName(),
			),
		)
	}

	return nil
}

func checkSingleResult(
	state common.State,
	result common.Type,
	returnType FunctionReturnType,
	funcName string,
) error {
	if result.(ObjectInstance).GetPrototype() == Nil {
		if returnType.Type != Nil && !returnType.IsNullable {
			resultStr, err := result.String(state)
			if err != nil {
				return err
			}

			return util.RuntimeError(
				fmt.Sprintf(
					"%s() повертає ненульове значення%s, отримано '%s'",
					funcName, "%s", resultStr,
				),
			)
		}
	} else if returnType.Type != Any && result.(ObjectInstance).GetPrototype() != returnType.Type {
		return util.RuntimeError(
			fmt.Sprintf(
				"'%s()' повертає значення типу '%s'%s, отримано значення з типом '%s'",
				funcName, returnType.String(), "%s", result.GetTypeName(),
			),
		)
	}

	return nil
}

func makeArgumentError(argsLen, parametersLen int, params []FunctionParameter, funcName string) error {
	diffLen := argsLen - parametersLen
	if diffLen > 0 {
		end1, end2, end3 := getEndings(diffLen)
		parametersStr := ""
		for c := parametersLen; c < argsLen; c++ {
			parametersStr += fmt.Sprintf("'%s'", params[c].Name)
			if c < argsLen-2 {
				parametersStr += ", "
			} else if c < argsLen-1 {
				parametersStr += " та "
			}
		}

		return util.RuntimeError(
			fmt.Sprintf(
				"при виклику '%s()' відсутн%s %d необхідн%s параметр%s: %s",
				funcName, end1, diffLen, end2, end3, parametersStr,
			),
		)
	}

	_, end1, end2 := getEndings(argsLen)
	return util.RuntimeError(
		fmt.Sprintf(
			"'%s()' приймає %d необхідн%s параметр%s, отримано %d",
			funcName, argsLen, end1, end2, parametersLen,
		),
	)
}

func getEndings(value int) (string, string, string) {
	end1 := "ій"
	end2 := "ий"
	end3 := ""
	if value > 1 && value < 5 {
		end1 = "і"
		end2 = "і"
		end3 = "и"
	} else if value != 1 {
		end1 = "і"
		end2 = "их"
		end3 = "ів"
	}

	return end1, end2, end3
}
