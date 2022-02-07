package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func checkArgument(parameter *FunctionParameter, arg common.Value, isVariadic bool) error {
	if parameter == nil {
		return errors.New("checkArgument: parameter is nil")
	}

	if arg == nil {
		return errors.New("checkArgument: arg is nil")
	}

	if parameter.Type == Any {
		return nil
	}

	argPrototype := arg.(ObjectInstance).GetClass()
	if argPrototype == Nil && parameter.IsNullable {
		return nil
	}

	if parameter.Type == argPrototype || argPrototype.HasBase(parameter.Type) {
		return nil
	}

	if isVariadic {
		return util.RuntimeError(
			fmt.Sprintf(
				"аргумент '%s' очікує набір параметрів з типом '%s' або його похідними, отримано '%s'",
				parameter.Name, parameter.GetTypeName(), arg.GetTypeName(),
			),
		)
	}

	return util.RuntimeError(
		fmt.Sprintf(
			"аргумент '%s' очікує параметр з типом '%s' або його похідними, отримано '%s'",
			parameter.Name, parameter.GetTypeName(), arg.GetTypeName(),
		),
	)
}

func CheckFunctionArguments(function *FunctionInstance, args *[]common.Value, _ *map[string]common.Value) error {
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
		if err := checkArgument(&function.Parameters[c], (*args)[c], false); err != nil {
			return err
		}
	}

	if len(function.Parameters) > 0 {
		if lastArgument := function.Parameters[len(function.Parameters)-1]; lastArgument.IsVariadic {
			if len(*args)-parametersLen > 0 {
				parametersLen = len(*args)
				for k := c; k < parametersLen; k++ {
					if err := checkArgument(&lastArgument, (*args)[k], true); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func CheckResult(state common.State, result common.Value, function *FunctionInstance) error {
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
	result common.Value,
	returnType FunctionReturnType,
	funcName string,
) error {
	if result.(ObjectInstance).GetClass() == Nil {
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
	} else if returnType.Type != Any && result.(ObjectInstance).GetClass() != returnType.Type {
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
