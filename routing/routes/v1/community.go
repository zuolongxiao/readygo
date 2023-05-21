package v1

import (
	apiv1 "readygo/api/v1"
	"readygo/routing/routes"
)

var community = []routes.Route{
	{
		Method:  "GET",
		Pattern: "/communities",
		Handler: apiv1.ListCommunities,
		Flag:    "Y",
		Desc:    "List communities",
	},
	{
		Method:  "POST",
		Pattern: "/communities",
		Handler: apiv1.CreateCommunity,
		Flag:    "Y",
		Desc:    "Create community",
	},
	{
		Method:  "PUT",
		Pattern: "/communities/:id",
		Handler: apiv1.UpdateCommunity,
		Flag:    "Y",
		Desc:    "Update community",
	},
	{
		Method:  "DELETE",
		Pattern: "/communities/:id",
		Handler: apiv1.DeleteCommunity,
		Flag:    "Y",
		Desc:    "Delete community",
	},
}

func init() {
	Routes = append(Routes, community...)
}
