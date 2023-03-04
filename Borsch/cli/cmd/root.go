package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/interpreter"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
	"github.com/alecthomas/participle/v2"
	"github.com/spf13/cobra"
)

var (
	stdRoot string
	codeArg string
)

var rootCmd = &cobra.Command{
	Use: "borsch",
	Long: `Борщ — це мова програмування, яка дозволяє писати код українською.
Вихідний код доступний на GitHub — https://github.com/YuriyLisovskiy/borsch-lang`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			fileInfo, err := os.Stat(args[0])
			if err != nil || fileInfo.IsDir() {
				return fmt.Errorf("'%s' is not a file", args[0])
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(stdRoot) == 0 {
			stdRoot = os.Getenv(builtin.BORSCH_LIB)
		}

		if len(stdRoot) == 0 {
			fmt.Printf(
				"Увага: змінна середовища '%s' необхідна для використання стандартної бібліотеки\n\n",
				builtin.BORSCH_LIB,
			)
		}

		if len(codeArg) > 0 {
			runCode(codeArg)
		} else if len(args) > 0 {
			filePath, err := filepath.Abs(args[0])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			runFile(filePath)
		} else {
			// interpret := interpreter.NewInterpreter(stdRoot, builtin.RootPackageName, "")
			// fmt.Printf("%s %s (%s, %s)\n", build.LanguageName, build.Version, build.Time, strings.Title(runtime.GOOS))
			// fmt.Println(
			// 	"Надрукуйте \"допомога();\", \"авторське_право();\" або \"ліцензія();\" для детальнішої інформації.\n" +
			// 		"Натисніть CONTROL+D або CONTROL+C для виходу.",
			// )
			// runInteractiveConsole(interpret)
		}
	},
}

func runCode(code string) {
	run(
		func(i interpreter.Interpreter) (types.Object, error) {
			return i.Evaluate("__вхід__", code, nil)
		},
	)
}

func runFile(filename string) {
	run(
		func(i interpreter.Interpreter) (types.Object, error) {
			return i.Import(filename)
		},
	)
}

func run(fn func(i interpreter.Interpreter) (types.Object, error)) {
	parser, err := interpreter.NewParser()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	stacktrace := &common.StackTrace{}
	state := interpreter.NewInitialState(nil, nil, stacktrace)
	i := interpreter.NewInterpreter(parser, state)
	_, err = fn(i)
	if err != nil {
		if pErr, ok := err.(participle.UnexpectedTokenError); ok {
			text := processParseError(pErr.Message())
			err = utilities.ParseError(pErr.Position(), pErr.Unexpected.Value, text)
		}

		fmt.Println(fmt.Sprintf("Відстеження (стек викликів):\n%s", stacktrace.String(err)))
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(
		&stdRoot, "lib", "l", "", "шлях до каталогу зі стандартною бібліотекою мови",
	)
	rootCmd.Flags().StringVarP(
		&codeArg, "code", "c", "", "вихідний код програми",
	)
}

func processParseError(text string) string {
	return strings.Replace(
		strings.Replace(
			text,
			"unexpected token",
			"неочікуваний токен",
			1,
		),
		"expected",
		"очікуваний",
		1,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
