package borsch

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/internal/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "версія",
	Short: "друк номеру збірки мови",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", config.LanguageName, config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
