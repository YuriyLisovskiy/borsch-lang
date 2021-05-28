package builtin

import (
	"fmt"
	"strings"
)

func Print(args... ValueType) (ValueType, error) {
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
	return NoneType{}, nil
}

func PrintLn(args... ValueType) (ValueType, error) {
	return Print(append(args, StringType{Value: "\n"})...)
}
