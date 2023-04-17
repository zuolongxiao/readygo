package cmd

import (
	"github.com/spf13/cobra"
)

// admin
// go run main.go admin
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "admin",
	Long:  `admin`,
}

func init() {
	rootCmd.AddCommand(adminCmd)
}
