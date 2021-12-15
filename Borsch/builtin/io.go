package builtin

import (
	"bufio"
	"fmt"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
	"os"
	"strings"
)

func Print(args ...types.Type) {
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
}

func Input(args ...types.Type) (types.Type, error) {
	Print(args...)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, util.InternalError(err.Error())
	}

	return types.StringType{Value: strings.TrimSuffix(input, "\n")}, nil
}
