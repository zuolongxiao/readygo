package v1

import (
	apiv1 "readygo/api/v1"
	"readygo/routing/routes"
)

var zone = []routes.Route{
	{
		Method:  "GET",
		Pattern: "/zones",
		Handler: apiv1.ListZones,
		Flag:    "Y",
		Desc:    "List zones",
	},
	{
		Method:  "POST",
		Pattern: "/zones",
		Handler: apiv1.CreateZone,
		Flag:    "Y",
		Desc:    "Create zone",
	},
	{
		Method:  "PUT",
		Pattern: "/zones/:id",
		Handler: apiv1.UpdateZone,
		Flag:    "Y",
		Desc:    "Update zone",
	},
	{
		Method:  "DELETE",
		Pattern: "/zones/:id",
		Handler: apiv1.DeleteZone,
		Flag:    "Y",
		Desc:    "Delete zone",
	},
}

func init() {
	Routes = append(Routes, zone...)
}
