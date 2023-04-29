package v1

import (
	apiv1 "readygo/api/v1"
	"readygo/routing/routes"
)

var sys = []routes.Route{
	{
		Method:  "POST",
		Pattern: "/auth",
		Handler: apiv1.Auth,
		Flag:    "-",
		Desc:    "Authentication",
	},
	{
		Method:  "GET",
		Pattern: "/captcha",
		Handler: apiv1.GenerateCaptcha,
		Flag:    "-",
		Desc:    "Generate Captcha",
	},
	{
		Method:  "GET",
		Pattern: "/admins",
		Handler: apiv1.ListAdmins,
		Flag:    "Y",
		Desc:    "List admins",
	},
	{
		Method:  "POST",
		Pattern: "/admins",
		Handler: apiv1.CreateAdmin,
		Flag:    "Y",
		Desc:    "Create admin",
	},
	{
		Method:  "PUT",
		Pattern: "/admins/:id",
		Handler: apiv1.UpdateAdmin,
		Flag:    "Y",
		Desc:    "Update admin",
	},
	{
		Method:  "DELETE",
		Pattern: "/admins/:id",
		Handler: apiv1.DeleteAdmin,
		Flag:    "Y",
		Desc:    "Delete admin",
	},

	{
		Method:  "GET",
		Pattern: "/roles",
		Handler: apiv1.ListRoles,
		Flag:    "Y",
		Desc:    "List roles",
	},
	{
		Method:  "POST",
		Pattern: "/roles",
		Handler: apiv1.CreateRole,
		Flag:    "Y",
		Desc:    "Create role",
	},
	{
		Method:  "PUT",
		Pattern: "/roles/:id",
		Handler: apiv1.UpdateRole,
		Flag:    "Y",
		Desc:    "Update role",
	},
	{
		Method:  "DELETE",
		Pattern: "/roles/:id",
		Handler: apiv1.DeleteRole,
		Flag:    "Y",
		Desc:    "Delete role",
	},
	{
		Method:  "POST",
		Pattern: "/roles/:id/permissions",
		Handler: apiv1.AddRolePermission,
		Flag:    "Y",
		Desc:    "Add permission to role",
	},
	{
		Method:  "PUT",
		Pattern: "/roles/:id/permissions",
		Handler: apiv1.UpdateRolePermission,
		Flag:    "Y",
		Desc:    "Update permissions of a role",
	},
	{
		Method:  "GET",
		Pattern: "/roles/:id/permissions",
		Handler: apiv1.ListRolePermissions,
		Flag:    "Y",
		Desc:    "List permissions of role",
	},
	{
		Method:  "DELETE",
		Pattern: "/roles/:id/permissions/:permissionID",
		Handler: apiv1.DeleteRolePermission,
		Flag:    "Y",
		Desc:    "Delete permission of role",
	},

	{
		Method:  "GET",
		Pattern: "/permissions",
		Handler: apiv1.ListPermissions,
		Flag:    "Y",
		Desc:    "List permissions",
	},
	{
		Method:  "POST",
		Pattern: "/permissions",
		Handler: apiv1.CreatePermission,
		Flag:    "Y",
		Desc:    "Create permission",
	},
	{
		Method:  "PUT",
		Pattern: "/permissions/:id",
		Handler: apiv1.UpdatePermission,
		Flag:    "Y",
		Desc:    "Update permission",
	},
	{
		Method:  "DELETE",
		Pattern: "/permissions/:id",
		Handler: apiv1.DeletePermission,
		Flag:    "Y",
		Desc:    "Delete permission",
	},

	{
		Method:  "GET",
		Pattern: "/profile",
		Handler: apiv1.GetProfile,
		Flag:    "-",
		Desc:    "Get profile",
	},
	{
		Method:  "PUT",
		Pattern: "/profile",
		Handler: apiv1.UpdateProfile,
		Flag:    "-",
		Desc:    "Update profile",
	},
}

func init() {
	Routes = append(Routes, sys...)
}
