package builtin

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

func evalPrint(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	parts := *args
	var strArgs []string
	for _, arg := range parts {
		argStr, err := arg.String(state)
		if err != nil {
			return nil, err
		}

		strArgs = append(strArgs, argStr)
	}

	fmt.Print(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Join(strArgs, ""), `\n`, "\n", -1,
				), `\r`, "\r", -1,
			), `\t`, "\t", -1,
		),
	)

	return types.NewNilInstance(), nil
}

func evalPrintLine(state common.State, args *[]common.Value, kwargs *map[string]common.Value) (common.Value, error) {
	*args = append(*args, types.NewStringInstance("\n"))
	return evalPrint(state, args, kwargs)
}

func evalInput(state common.State, args *[]common.Value, kwargs *map[string]common.Value) (common.Value, error) {
	_, err := evalPrint(state, args, kwargs)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, utilities.InternalError(err.Error())
	}

	return types.StringInstance{Value: strings.TrimSuffix(input, "\n")}, nil
}

func makePrintFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"друк",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "х",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		evalPrint,
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}

func makePrintLnFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"друкр",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "а",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		evalPrintLine,
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}

func makeInputFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"ввід",
		[]types.FunctionParameter{
			{
				Type:       types.String,
				Name:       "повідомлення",
				IsVariadic: true,
				IsNullable: false,
			},
		},
		evalInput,
		[]types.FunctionReturnType{
			{
				Type:       types.String,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
