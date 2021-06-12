package api

import (
	"readygo/models"
	"readygo/pkg/utils"
	"readygo/services"

	"github.com/gin-gonic/gin"
)

// ListPermissions ListPermissions
func ListPermissions(c *gin.Context) {
	w := utils.NewContextWrapper(c)
	s := services.New(&models.Permission{})

	var list []models.PermissionView
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

// CreatePermission CreatePermission
func CreatePermission(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	binding := models.PermissionCreate{}
	if err := w.Bind(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m := models.Permission{}
	s := services.New(&m)
	if err := s.Fill(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

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

// UpdatePermission UpdatePermission
func UpdatePermission(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	binding := models.PermissionUpdate{}
	if err := w.Bind(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m := models.Permission{}
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

// DeletePermission DeletePermission
func DeletePermission(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	s := services.New(&models.Permission{})
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
