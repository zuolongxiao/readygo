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
		"list":   list,
		"offset": svc.GetOffset(),
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
	s := services.New(&mdl)
	if err := s.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl.CreatedBy = cw.GetUsername()
	if err := s.Create(); err != nil {
		cw.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         mdl.ID,
		"created_at": mdl.CreatedAt.Time,
	}

	cw.Respond(nil, data)
}

// UpdateRole update role
func UpdateRole(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	binding := models.RoleUpdate{}
	if err := w.Bind(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m := models.Role{}
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

// DeleteRole delete role
func DeleteRole(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	s := services.New(&models.Role{})
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

// ListRolePermissions ListRolePermissions
func ListRolePermissions(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	roleID, _ := strconv.Atoi(c.Param("id"))
	if roleID <= 0 {
		w.Respond(errs.ValidationError("invalid role ID"), nil)
		return
	}

	roleModel := models.Role{}
	roleSvc := services.New(&roleModel)
	if err := roleSvc.LoadByID(c.Param("id")); err != nil {
		w.Respond(err, nil)
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
		w.Respond(err, nil)
		return
	}

	permIDsQueryer := models.IDsQueryer{
		List: rolePermList,
		Key:  "PermissionID",
	}
	permSvc := services.New(&models.Permission{})
	var permList []models.RoleView
	if err := permSvc.Find(&permList, &permIDsQueryer); err != nil {
		w.Respond(err, nil)
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
	var lst []List
	for _, row := range rolePermList {
		dst := List{}
		copier.CopyWithOption(&dst, row, copier.Option{IgnoreEmpty: true})
		dst.PermissionName = permNameDict[row.PermissionID]
		dst.RoleName = roleModel.Name
		lst = append(lst, dst)
	}

	data := map[string]interface{}{
		"list":   lst,
		"offset": svc.GetOffset(),
	}

	w.Respond(nil, data)
}

// AddRolePermission AddRolePermission
func AddRolePermission(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	roleID, _ := strconv.Atoi(c.Param("id"))
	if roleID <= 0 {
		w.Respond(errs.ValidationError("invalid role ID"), nil)
		return
	}

	binding := models.AuthorizationBinding{}
	if err := w.Bind(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m := models.Authorization{}
	s := services.New(&m)
	if err := s.Fill(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m.RoleID = uint64(roleID)
	m.CreatedBy = w.GetUsername()
	if err := s.Create(); err != nil {
		w.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         m.ID,
		"updated_at": m.UpdatedAt.Time,
	}

	w.Respond(nil, data)
}

// DeleteRolePermission DeleteRolePermission
func DeleteRolePermission(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	roleID, _ := strconv.Atoi(c.Param("id"))
	if roleID <= 0 {
		w.Respond(errs.ValidationError("invalid role ID"), nil)
		return
	}

	permissionID, _ := strconv.Atoi(c.Param("permissionID"))
	if permissionID <= 0 {
		w.Respond(errs.ValidationError("invalid permission ID"), nil)
		return
	}

	m := models.Authorization{
		RoleID:       uint64(roleID),
		PermissionID: uint64(permissionID),
	}
	s := services.New(&m)
	if err := s.Load(); err != nil {
		w.Respond(err, nil)
		return
	}

	if err := s.Delete(); err != nil {
		w.Respond(err, nil)
		return
	}

	w.Respond(nil, nil)
}
