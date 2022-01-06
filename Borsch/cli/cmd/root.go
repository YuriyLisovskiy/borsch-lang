package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/cli/build"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/interpreter"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
	"github.com/spf13/cobra"
)

var (
	stdRoot string
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
			stdRoot = os.Getenv("BORSCH_STD")
		}

		if len(stdRoot) == 0 {
			fmt.Print("Увага: змінна середовища BORSCH_STD необхідна для використання стандартної бібліотеки\n\n")
		}

		if len(args) > 0 {
			filePath := args[0]
			interpret := interpreter.NewInterpreter(stdRoot, filePath, "")
			content, err := util.ReadFile(filePath)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				context := interpreter.NewContext(filePath, "", stdRoot)
				err = interpret.ExecuteFile(context, content)
				if err != nil {
					fmt.Println(fmt.Sprintf("Відстеження (стек викликів):\n%s", err.Error()))
				}
			}
		} else {
			interpret := interpreter.NewInterpreter(stdRoot, builtin.RootPackageName, "")
			fmt.Printf("%s %s (%s, %s)\n", build.LanguageName, build.Version, build.Time, strings.Title(runtime.GOOS))
			fmt.Println(
				"Надрукуйте \"допомога();\", \"авторське_право();\" або \"ліцензія();\" для детальнішої інформації.\n" +
					"Натисніть CONTROL+D або CONTROL+C для виходу.",
			)
			runInteractiveConsole(interpret)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(
		&stdRoot, "lib", "l", "", "шлях до директорії зі стандартною бібліотекою мови",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
