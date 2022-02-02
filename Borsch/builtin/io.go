package builtin

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func Print(state common.State, args ...common.Type) error {
	var strArgs []string
	for _, arg := range args {
		argStr, err := arg.String(state)
		if err != nil {
			return err
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

	return nil
}

func Input(state common.State, args ...common.Type) (common.Type, error) {
	err := Print(state, args...)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, util.InternalError(err.Error())
	}

	return types.StringInstance{Value: strings.TrimSuffix(input, "\n")}, nil
}
