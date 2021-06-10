package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of the Borsch programming language",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Borsch %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
