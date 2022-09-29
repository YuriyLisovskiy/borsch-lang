package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "borsch",
	Long: `Борщ — це мова програмування, яка дозволяє писати код українською.
Вихідний код доступний на GitHub — https://github.com/YuriyLisovskiy/borsch-lang`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
