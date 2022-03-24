package cmd

import (
	"fmt"
	"os"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "друк інформації про змінні середовища для мови Борщ",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s=\"%s\"\n", builtin.BORSCH_LIB, os.Getenv(builtin.BORSCH_LIB))
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
