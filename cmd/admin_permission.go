package cmd

import (
	"fmt"
	"readygo/models"
	"readygo/pkg/db"
	"readygo/pkg/errs"
	v1 "readygo/routing/routes/v1"
	"readygo/services"
	"reflect"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Load permissions into database
// go run main.go admin permission
var adminPermissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "Load permissions to database",
	Long:  `Load permissions to database`,
	Run: func(cmd *cobra.Command, args []string) {
		loadPermissions()
	},
}

func init() {
	adminCmd.AddCommand(adminPermissionCmd)
}

func loadPermissions() {
	if err := db.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	cw := EmptyContextWrapper{}

	for _, v := range v1.Routes {
		handler := runtime.FuncForPC(reflect.ValueOf(v.Handler).Pointer()).Name()
		tmp := strings.Split(handler, ".")
		name := tmp[len(tmp)-1:][0]

		if v.Flag == "" || v.Flag == "-" {
			fmt.Printf("name: %s, ignore\n", name)
			continue
		}

		isEnabled := "N"
		if v.Flag == "Y" {
			isEnabled = "Y"
		}

		mdl := models.Permission{
			Name: name,
		}

		svc := services.New(&mdl)
		if err := svc.Load(); err != nil {
			if _, ok := err.(errs.NotFoundError); ok {
				mdl.Title = v.Desc
				mdl.IsEnabled = isEnabled
				if err := svc.Create(cw); err != nil {
					fmt.Printf("name: %s, create error: %s\n", name, err.Error())
				} else {
					fmt.Printf("name: %s, added\n", name)
				}
			} else {
				fmt.Printf("name: %s, load error: %s\n", name, err.Error())
				break
			}
		} else {
			fmt.Printf("name: %s, exists, skip\n", name)
		}
	}

	fmt.Println("done")
}
