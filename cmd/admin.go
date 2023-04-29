package cmd

import (
	"github.com/spf13/cobra"
)

type EmptyContextWrapper struct{}

func (cw EmptyContextWrapper) GetUsername() string {
	return ""
}

// go run main.go admin
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin",
	Long:  `Admin`,
}

func init() {
	rootCmd.AddCommand(adminCmd)
}
