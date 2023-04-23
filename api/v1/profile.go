package api

import (
	"fmt"

	"readygo/models"
	"readygo/pkg/errs"
	"readygo/pkg/settings"
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
	profile.CreatedAt = models.LocalTime(adminMdl.CreatedAt.Time)
	profile.UpdatedAt = models.LocalTime(adminMdl.UpdatedAt.Time)
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
	w := utils.NewContextWrapper(c)

	binding := models.ProfileUpdate{}
	if err := w.Bind(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m := models.Admin{}
	s := services.New(&m)
	if err := s.LoadByKey("username", w.GetUsername()); err != nil {
		w.Respond(err, nil)
		return
	}

	if !utils.VerifyPassword(m.Password, binding.PasswordOld) {
		w.Respond(errs.ValidationError("incorrect password"), nil)
		return
	}

	if err := s.Fill(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	if binding.Password != "" {
		hashedPassword, err := utils.HashPassword(binding.Password)
		if err != nil {
			w.Respond(errs.InternalServerError(err.Error()), nil)
		}
		m.Password = hashedPassword
	}

	m.UpdatedBy = w.GetUsername()
	if err := s.Save(); err != nil {
		w.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         m.ID,
		"updated_at": m.UpdatedAt.Time.Format(settings.App.TimeFormat),
	}

	w.Respond(nil, data)
}
