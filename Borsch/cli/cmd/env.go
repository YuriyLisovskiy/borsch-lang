package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "друк інформації про змінні середовища для мови Борщ",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("BORSCH_STD=\"%s\"\n", os.Getenv("BORSCH_STD"))
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
