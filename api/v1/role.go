package api

import (
	"strconv"

	"readygo/models"
	"readygo/pkg/errs"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// ListRoles list roles
func ListRoles(ctx *gin.Context) {
	cw := utils.NewContextWrapper(ctx)
	svc := services.New(&models.Role{})

	var list []models.RoleView
	if err := svc.Find(&list, ctx); err != nil {
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

// CreateRole create role
func CreateRole(ctx *gin.Context) {
	cw := utils.NewContextWrapper(ctx)

	binding := models.RoleCreate{}
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	var mdl models.Role
	svc := services.New(&mdl)
	if err := svc.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Create(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	var view models.RoleView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// UpdateRole update role
func UpdateRole(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	binding := models.RoleUpdate{}
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Role{}
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

	var view models.RoleView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// DeleteRole delete role
func DeleteRole(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	svc := services.New(&models.Role{})
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

// ListRolePermissions ListRolePermissions
func ListRolePermissions(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	if roleID <= 0 {
		cw.Respond(errs.ValidationError("invalid role ID"), nil)
		return
	}

	roleModel := models.Role{}
	roleSvc := services.New(&roleModel)
	if err := roleSvc.LoadByID(c.Param("id")); err != nil {
		cw.Respond(err, nil)
		return
	}

	kvQueryer := models.KeyValueQueryer{
		Entries: map[string]string{
			"role_id": c.Param("id"),
		},
	}
	svc := services.New(&models.Authorization{})
	var rolePermList []models.AuthorizationView
	if err := svc.Find(&rolePermList, &kvQueryer); err != nil {
		cw.Respond(err, nil)
		return
	}

	permIDsQueryer := models.IDsQueryer{
		List: rolePermList,
		Key:  "PermissionID",
	}
	permSvc := services.New(&models.Permission{})
	var permList []models.RoleView
	if err := permSvc.Find(&permList, &permIDsQueryer); err != nil {
		cw.Respond(err, nil)
		return
	}
	permNameDict := make(map[uint64]string)
	for _, perm := range permList {
		permNameDict[perm.ID] = perm.Name
	}
	type List struct {
		models.AuthorizationView
		RoleName       string `json:"role_name"`
		PermissionName string `json:"permission_name"`
	}
	lst := []List{}
	for _, row := range rolePermList {
		dst := List{}
		copier.CopyWithOption(&dst, row, copier.Option{IgnoreEmpty: true})
		dst.PermissionName = permNameDict[row.PermissionID]
		dst.RoleName = roleModel.Name
		lst = append(lst, dst)
	}

	data := map[string]interface{}{
		"list": lst,
		"prev": svc.GetPrev(),
		"next": svc.GetNext(),
	}

	cw.Respond(nil, data)
}

// AddRolePermission AddRolePermission
func AddRolePermission(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	if roleID <= 0 {
		cw.Respond(errs.ValidationError("invalid role ID"), nil)
		return
	}

	binding := models.AuthorizationBinding{}
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Authorization{}
	svc := services.New(&mdl)
	if err := svc.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl.RoleID = roleID
	if err := svc.Create(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         mdl.ID,
		"updated_at": mdl.UpdatedAt,
	}

	cw.Respond(nil, data)
}

// DeleteRolePermission DeleteRolePermission
func DeleteRolePermission(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	if roleID <= 0 {
		cw.Respond(errs.ValidationError("invalid role ID"), nil)
		return
	}

	permissionID, _ := strconv.ParseUint(c.Param("permissionID"), 10, 0)
	if permissionID <= 0 {
		cw.Respond(errs.ValidationError("invalid permission ID"), nil)
		return
	}

	mdl := models.Authorization{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	svc := services.New(&mdl)
	if err := svc.Load(); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Delete(); err != nil {
		cw.Respond(err, nil)
		return
	}

	cw.Respond(nil, nil)
}

// UpdateRolePermission UpdateRolePermission
func UpdateRolePermission(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	if roleID <= 0 {
		cw.Respond(errs.ValidationError("invalid role ID"), nil)
		return
	}

	binding := make(map[string][]uint64)
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}
	bindingPermissionIDs, ok := binding["permission_ids"]
	if !ok {
		cw.Respond(errs.ValidationError("missing permission_ids"), nil)
		return
	}
	bindingPermissionIDs = utils.RemoveDuplicate(bindingPermissionIDs)

	roleModel := models.Role{}
	roleSvc := services.New(&roleModel)
	if err := roleSvc.LoadByID(c.Param("id")); err != nil {
		cw.Respond(err, nil)
		return
	}

	kvQueryer := models.KeyValueQueryer{
		Entries: map[string]string{
			"role_id": c.Param("id"),
		},
	}
	svc := services.New(&models.Authorization{})
	var rolePermList []models.AuthorizationView
	if err := svc.Find(&rolePermList, &kvQueryer); err != nil {
		cw.Respond(err, nil)
		return
	}

	permissionIDs := []uint64{}
	for _, auth := range rolePermList {
		permissionIDs = append(permissionIDs, auth.PermissionID)
	}

	addPermissionIDs := utils.Difference(bindingPermissionIDs, permissionIDs)
	subPermissionIDs := utils.Difference(permissionIDs, bindingPermissionIDs)

	// db.DB.Delete(&models.Authorization{}, "role_id = ? AND permission_id IN ?", roleID, subPermissionIds)
	for _, permissionID := range subPermissionIDs {
		mdl := models.Authorization{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
		svc := services.New(&mdl)
		if err := svc.Load(); err != nil {
			continue
		}
		if err := svc.Delete(); err != nil {
			continue
		}
	}

	for _, permissionID := range addPermissionIDs {
		mdl := models.Authorization{}
		svc := services.New(&mdl)
		binding := models.AuthorizationBinding{
			RoleID:       roleID,
			PermissionID: permissionID,
		}

		if err := svc.Fill(&binding); err != nil {
			continue
		}

		if err := svc.Create(cw); err != nil {
			continue
		}
	}

	cw.Respond(nil, nil)
}
