package api

import (
	"strconv"

	"readygo/models"
	"readygo/pkg/errs"
	"readygo/pkg/settings"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// ListAdmins ListAdmins
func ListAdmins(c *gin.Context) {
	w := utils.NewContextWrapper(c)
	s := services.New(&models.Admin{})

	var aminList []models.AdminView
	if err := s.Find(&aminList, c); err != nil {
		w.Respond(err, nil)
		return
	}

	roleIDsQueryer := models.IDsQueryer{
		List: aminList,
		Key:  "RoleID",
	}
	roleSvc := services.New(&models.Role{})
	var roleList []models.RoleView
	if err := roleSvc.Find(&roleList, &roleIDsQueryer); err != nil {
		w.Respond(err, nil)
		return
	}

	roleNameDict := make(map[uint64]string)
	for _, role := range roleList {
		roleNameDict[role.ID] = role.Name
	}

	type List struct {
		models.AdminView
		RoleName string `json:"role_name"`
	}
	var lst []List
	for _, row := range aminList {
		dst := List{}
		copier.CopyWithOption(&dst, row, copier.Option{IgnoreEmpty: true})
		dst.RoleName = roleNameDict[row.RoleID]
		lst = append(lst, dst)
	}

	resp := map[string]interface{}{
		"list": lst,
		"prev": s.GetPrev(),
		"next": s.GetNext(),
	}

	w.Respond(nil, resp)
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
	if id == settings.App.SuperAdminID {
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
	if id == settings.App.SuperAdminID {
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
