package v1

import (
	apiv1 "readygo/api/v1"
	"readygo/routing/routes"
)

var tag = []routes.Route{
	{
		Method:  "GET",
		Pattern: "/tags/:id",
		Handler: apiv1.GetTag,
		Flag:    "Y",
		Desc:    "Get tag",
	},
	{
		Method:  "GET",
		Pattern: "/tags",
		Handler: apiv1.ListTags,
		Flag:    "Y",
		Desc:    "List tags",
	},
	{
		Method:  "POST",
		Pattern: "/tags",
		Handler: apiv1.CreateTag,
		Flag:    "Y",
		Desc:    "Create tag",
	},
	{
		Method:  "PUT",
		Pattern: "/tags/:id",
		Handler: apiv1.UpdateTag,
		Flag:    "Y",
		Desc:    "Update tag",
	},
	{
		Method:  "DELETE",
		Pattern: "/tags/:id",
		Handler: apiv1.DeleteTag,
		Flag:    "Y",
		Desc:    "Delete tag",
	},
}

func init() {
	Routes = append(Routes, tag...)
}
