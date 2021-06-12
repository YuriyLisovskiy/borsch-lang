package cmd

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/interpreter"
	"github.com/spf13/cobra"
	"os"
)

var stdRoot string

var rootCmd = &cobra.Command{
	Use: "borsch",
	Long: `Borsch is a programming language that lets you write code in Ukrainian.
The source code is available at https://github.com/YuriyLisovskiy/borsch-lang`,
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

		interpret := interpreter.NewInterpreter(stdRoot)
		if len(args) > 0 {
			filePath := args[0]
			err := interpret.ExecuteFile(filePath)
			if err != nil {
				fmt.Println(fmt.Sprintf("Відстеження (стек викликів):\n%s", err.Error()))
			}
		} else {
			runInteractiveConsole(interpret)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(
		&stdRoot, "stdlib", "l", "", "path to root directory of Borsch standard library",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
