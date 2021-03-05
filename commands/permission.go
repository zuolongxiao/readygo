package commands

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/zuolongxiao/readygo/models"
	"github.com/zuolongxiao/readygo/pkg/errs"
	v1 "github.com/zuolongxiao/readygo/routing/routes/v1"
	"github.com/zuolongxiao/readygo/services"
)

// LoadPermissions load permission into database
func LoadPermissions() {
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

		m := models.Permission{
			Name: name,
		}

		s := services.New(&m)
		if err := s.Load(); err != nil {
			if _, ok := err.(errs.NotFoundError); ok {
				m.Title = v.Desc
				m.IsEnabled = isEnabled
				if err := s.Create(); err != nil {
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
}
