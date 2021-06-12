package api

import (
	"strconv"

	"readygo/models"
	"readygo/pkg/errs"
	"readygo/pkg/settings"
	"readygo/pkg/utils"
	"readygo/services"

	"github.com/gin-gonic/gin"
)

// ListAdmins ListAdmins
func ListAdmins(c *gin.Context) {
	w := utils.NewContextWrapper(c)
	s := services.New(&models.Admin{})

	var list []models.AdminView
	if err := s.Find(&list, c); err != nil {
		w.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"list":   list,
		"offset": s.GetOffset(),
	}

	w.Respond(nil, data)
}

// CreateAdmin CreateAdmin
func CreateAdmin(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	var binding models.AdminCreate
	if err := w.Bind(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m := models.Admin{}
	s := services.New(&m)
	if err := s.Fill(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m.IPAddr = c.ClientIP()
	m.CreatedBy = w.GetUsername()
	if err := s.Create(); err != nil {
		w.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         m.ID,
		"created_at": m.CreatedAt.Time,
	}

	w.Respond(nil, data)
}

// UpdateAdmin UpdateAdmin
func UpdateAdmin(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	var binding models.AdminUpdate
	if err := w.Bind(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	if id == settings.AppSetting.SuperAdminID {
		w.Respond(errs.ForbiddenError("super admin cannot be modified"), nil)
		return
	}

	m := models.Admin{}
	s := services.New(&m)
	if err := s.LoadByID(c.Param("id")); err != nil {
		w.Respond(err, nil)
		return
	}

	if err := s.Fill(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m.UpdatedBy = w.GetUsername()
	if err := s.Save(); err != nil {
		w.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         m.ID,
		"updated_at": m.UpdatedAt.Time,
	}

	w.Respond(nil, data)
}

// DeleteAdmin DeleteAdmin
func DeleteAdmin(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	id, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	if id == settings.AppSetting.SuperAdminID {
		w.Respond(errs.ForbiddenError("super admin cannot be deleted"), nil)
		return
	}

	s := services.New(&models.Admin{})
	if err := s.LoadByID(c.Param("id")); err != nil {
		w.Respond(err, nil)
		return
	}

	if err := s.Delete(); err != nil {
		w.Respond(err, nil)
		return
	}

	w.Respond(nil, nil)
}
