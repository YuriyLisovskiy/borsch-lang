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

func Print(ctx common.Context, args ...common.Type) {
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, arg.String(ctx))
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

func Input(ctx common.Context, args ...common.Type) (common.Type, error) {
	Print(ctx, args...)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, util.InternalError(err.Error())
	}

	return types.StringInstance{Value: strings.TrimSuffix(input, "\n")}, nil
}
