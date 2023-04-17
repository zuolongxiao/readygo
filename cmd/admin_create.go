package cmd

import (
	"fmt"
	"readygo/models"
	"readygo/pkg/db"
	"readygo/pkg/errs"
	"readygo/services"
	"readygo/utils"

	"github.com/spf13/cobra"
)

var username string
var password string

// 创建管理员
// go run main.go admin create -u zuolongxiao -p zuolongxiao
var adminCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create admin",
	Long:  `Create admin`,
	Run: func(cmd *cobra.Command, args []string) {
		if password == "" {
			password = utils.RandomString(12)
		}
		if err := createAdmin(username, password); err != nil {
			fmt.Printf("Admin created failed with error: %v\n", err)
			return
		}

		fmt.Printf("Admin created successfully\n")
		fmt.Printf("Username: %s\n", username)
		fmt.Printf("Password: %s\n", password)
	},
}

func init() {
	adminCmd.AddCommand(adminCreateCmd)
	adminCreateCmd.Flags().StringVarP(&username, "username", "u", "", "用户名")
	adminCreateCmd.Flags().StringVarP(&password, "password", "p", "", "密码")

	adminCreateCmd.MarkFlagRequired("username")
}

func createAdmin(username, password string) error {
	if err := db.Setup(); err != nil {
		fmt.Println(err)
		return err
	}

	if len(username) < 2 || len(username) > 40 {
		return errs.ValidationError("username length must be 2-40")
	}
	if len(password) < 2 || len(password) > 40 {
		return errs.ValidationError("password length must be 2-40")
	}

	admin := models.Admin{
		Username: username,
		Password: password,
	}
	s := services.New(&admin)

	return s.Create()
}
