package builtin

import (
	"bufio"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"os"
	"strings"
)

func Print(args... ValueType) (ValueType, error) {
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, arg.Representation())
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
	return NoneType{}, nil
}

func PrintLn(args... ValueType) (ValueType, error) {
	return Print(append(args, StringType{Value: "\n"})...)
}

func Input(args... ValueType) (ValueType, error) {
	_, err := Print(args...)
	if err != nil {
		return NoneType{}, util.InternalError(err.Error())
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return NoneType{}, util.InternalError(err.Error())
	}

	return StringType{Value: strings.TrimSuffix(input, "\n")}, nil
}
