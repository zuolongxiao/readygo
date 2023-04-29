package api

import (
	"fmt"

	"readygo/models"
	"readygo/pkg/errs"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
)

// GetProfile GetProfile
func GetProfile(ctx *gin.Context) {
	cw := utils.NewContextWrapper(ctx)

	adminMdl := models.Admin{}
	adminSvc := services.New(&adminMdl)
	if err := adminSvc.LoadByKey("username", cw.GetUsername()); err != nil {
		cw.Respond(err, nil)
		return
	}

	var profile models.ProfileView
	profile.ID = adminMdl.ID
	profile.Username = adminMdl.Username
	profile.CreatedAt = adminMdl.CreatedAt
	profile.UpdatedAt = adminMdl.UpdatedAt
	profile.Roles = []string{}
	if adminMdl.RoleID > 0 {
		roleMdl := models.Role{}
		roleSvc := services.New(&roleMdl)
		if err := roleSvc.LoadByID(fmt.Sprint(adminMdl.RoleID)); err == nil {
			profile.Roles = append(profile.Roles, roleMdl.Name)
		}
	}

	cw.Respond(nil, profile)
}

// UpdateProfile UpdateProfile
func UpdateProfile(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	binding := models.ProfileUpdate{}
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Admin{}
	svc := services.New(&mdl)
	if err := svc.LoadByKey("username", cw.GetUsername()); err != nil {
		cw.Respond(err, nil)
		return
	}

	if !utils.VerifyPassword(mdl.Password, binding.PasswordOld) {
		cw.Respond(errs.ValidationError("incorrect password"), nil)
		return
	}

	if err := svc.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	if binding.Password != "" {
		hashedPassword, err := utils.HashPassword(binding.Password)
		if err != nil {
			cw.Respond(errs.InternalServerError(err.Error()), nil)
		}
		mdl.Password = hashedPassword
	}

	if err := svc.Save(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         mdl.ID,
		"updated_at": mdl.UpdatedAt,
	}

	cw.Respond(nil, data)
}
