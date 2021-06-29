package cmd

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/lang/cli/build"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of the Borsch programming language",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Borsch %s\n", build.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
