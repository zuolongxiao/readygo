package commands

import (
	"github.com/zuolongxiao/readygo/models"
	"github.com/zuolongxiao/readygo/pkg/errs"
	"github.com/zuolongxiao/readygo/services"
)

// CreateAdmin create an admin
func CreateAdmin(username, password string) error {
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
