package cmd

import (
	"fmt"
	"readygo/models"
	"readygo/pkg/db"

	"github.com/spf13/cobra"
)

// 数据库迁移
// go run main.go admin migrate
var adminMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate tables",
	Long:  `Migrate tables`,
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

func init() {
	adminCmd.AddCommand(adminMigrateCmd)
}

func migrate() {
	if err := db.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	db.DB.AutoMigrate(models.Migrations...)

	fmt.Println("done")
}
