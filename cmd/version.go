package cmd

import (
	"fmt"

	"readygo/pkg/settings"

	"github.com/spf13/cobra"
)

// go run main.go version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print app version",
	Long:  `Print app version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(settings.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
