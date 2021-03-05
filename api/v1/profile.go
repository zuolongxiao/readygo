package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zuolongxiao/readygo/models"
	"github.com/zuolongxiao/readygo/pkg/errs"
	"github.com/zuolongxiao/readygo/pkg/utils"
	"github.com/zuolongxiao/readygo/services"
)

// GetProfile GetProfile
func GetProfile(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	am := models.Admin{}
	as := services.New(&am)
	if err := as.LoadByKey("username", w.GetUsername()); err != nil {
		w.Respond(err, nil)
		return
	}

	var profile models.ProfileView
	profile.ID = am.ID
	profile.Username = am.Username
	profile.CreatedAt = am.CreatedAt.Time
	profile.UpdatedAt = am.UpdatedAt.Time
	if am.RoleID > 0 {
		rm := models.Role{}
		rs := services.New(&rm)
		if err := rs.LoadByID(fmt.Sprint(am.RoleID)); err == nil {
			profile.Role = rm.Name
		}
	}

	w.Respond(nil, profile)
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

	m.UpdatedBy = w.GetUsername()
	if err := s.Update(); err != nil {
		w.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         m.ID,
		"updated_at": m.UpdatedAt.Time,
	}

	w.Respond(nil, data)
}
