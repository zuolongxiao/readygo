package cmd

import (
	"github.com/spf13/cobra"
)

// admin
// go run main.go admin
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin",
	Long:  `Admin`,
}

func init() {
	rootCmd.AddCommand(adminCmd)
}
