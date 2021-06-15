package builtin

import (
	"bufio"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"os"
	"strings"
)

func Print(args ...types.ValueType) (types.ValueType, error) {
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, arg.String())
	}

	fmt.Print(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Join(strArgs, " "), `\n`, "\n", -1,
				), `\r`, "\r", -1,
			), `\t`, "\t", -1,
		),
	)
	return nil, nil
}

func PrintLn(args ...types.ValueType) (types.ValueType, error) {
	return Print(append(args, types.StringType{Value: "\n"})...)
}

func Input(args ...types.ValueType) (types.ValueType, error) {
	_, err := Print(args...)
	if err != nil {
		return nil, util.InternalError(err.Error())
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, util.InternalError(err.Error())
	}

	return types.StringType{Value: strings.TrimSuffix(input, "\n")}, nil
}
