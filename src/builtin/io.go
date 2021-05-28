package builtin

import (
	"fmt"
	"strings"
)

func Print(args... string) (string, error) {
	fmt.Print(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Join(args, " "), `\n`, "\n", -1,
				), `\r`, "\r", -1,
			), `\t`, "\t", -1,
		),
	)
	return "", nil
}

func PrintLn(args... string) (string, error) {
	return Print(append(args, "\n")...)
}
