package v1

import (
	apiv1 "readygo/api/v1"
	"readygo/routing/routes"
)

var district = []routes.Route{
	{
		Method:  "GET",
		Pattern: "/districts",
		Handler: apiv1.ListDistricts,
		Flag:    "Y",
		Desc:    "List districts",
	},
	{
		Method:  "POST",
		Pattern: "/districts",
		Handler: apiv1.CreateDistrict,
		Flag:    "Y",
		Desc:    "Create district",
	},
	{
		Method:  "PUT",
		Pattern: "/districts/:id",
		Handler: apiv1.UpdateDistrict,
		Flag:    "Y",
		Desc:    "Update district",
	},
	{
		Method:  "DELETE",
		Pattern: "/districts/:id",
		Handler: apiv1.DeleteDistrict,
		Flag:    "Y",
		Desc:    "Delete district",
	},
}

func init() {
	Routes = append(Routes, district...)
}
