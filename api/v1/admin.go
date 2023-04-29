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
	cw := utils.NewContextWrapper(c)
	svc := services.New(&models.Admin{})

	var aminList []models.AdminView
	if err := svc.Find(&aminList, c); err != nil {
		cw.Respond(err, nil)
		return
	}

	roleIDsQueryer := models.IDsQueryer{
		List: aminList,
		Key:  "RoleID",
	}
	roleSvc := services.New(&models.Role{})
	var roleList []models.RoleView
	if err := roleSvc.Find(&roleList, &roleIDsQueryer); err != nil {
		cw.Respond(err, nil)
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
		"prev": svc.GetPrev(),
		"next": svc.GetNext(),
	}

	cw.Respond(nil, resp)
}

// CreateAdmin CreateAdmin
func CreateAdmin(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	var binding models.AdminCreate
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Admin{}
	svc := services.New(&mdl)
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

	mdl.IPAddr = c.ClientIP()
	if err := svc.Create(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	var view models.AdminView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// UpdateAdmin UpdateAdmin
func UpdateAdmin(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	var binding models.AdminUpdate
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	if id == settings.App.SuperAdminID {
		cw.Respond(errs.ForbiddenError("super admin not allowed to change"), nil)
		return
	}

	mdl := models.Admin{}
	svc := services.New(&mdl)
	if err := svc.LoadByID(c.Param("id")); err != nil {
		cw.Respond(err, nil)
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

	var view models.AdminView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// DeleteAdmin DeleteAdmin
func DeleteAdmin(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	id, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	if id == settings.App.SuperAdminID {
		cw.Respond(errs.ForbiddenError("super admin role not allowed to delete"), nil)
		return
	}

	svc := services.New(&models.Admin{})
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
