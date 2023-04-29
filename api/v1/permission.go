package api

import (
	"readygo/models"
	"readygo/pkg/settings"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
)

// ListPermissions ListPermissions
func ListPermissions(c *gin.Context) {
	cw := utils.NewContextWrapper(c)
	svc := services.New(&models.Permission{})

	var list []models.PermissionView
	if err := svc.Find(&list, c); err != nil {
		cw.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"list": list,
		"prev": svc.GetPrev(),
		"next": svc.GetNext(),
	}

	cw.Respond(nil, data)
}

// CreatePermission CreatePermission
func CreatePermission(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	binding := models.PermissionCreate{}
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Permission{}
	svc := services.New(&mdl)
	if err := svc.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Create(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         mdl.ID,
		"created_at": mdl.CreatedAt.Time.Format(settings.App.TimeFormat),
	}

	cw.Respond(nil, data)
}

// UpdatePermission UpdatePermission
func UpdatePermission(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	binding := models.PermissionUpdate{}
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Permission{}
	svc := services.New(&mdl)
	if err := svc.LoadByID(c.Param("id")); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Save(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         mdl.ID,
		"updated_at": mdl.UpdatedAt.Time.Format(settings.App.TimeFormat),
	}

	cw.Respond(nil, data)
}

// DeletePermission DeletePermission
func DeletePermission(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	svc := services.New(&models.Permission{})
	if err := svc.LoadByID(c.Param("id")); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Delete(); err != nil {
		cw.Respond(err, nil)
		return
	}

	cw.Respond(nil, nil)
}
