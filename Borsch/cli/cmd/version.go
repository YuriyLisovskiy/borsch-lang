package cmd

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/cli/build"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "версія",
	Short: "друк номеру збірки мови",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", build.LanguageName, build.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
