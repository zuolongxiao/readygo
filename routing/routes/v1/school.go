package v1

import (
	apiv1 "readygo/api/v1"
	"readygo/routing/routes"
)

var school = []routes.Route{
	{
		Method:  "GET",
		Pattern: "/schools",
		Handler: apiv1.ListSchools,
		Flag:    "Y",
		Desc:    "List schools",
	},
	{
		Method:  "POST",
		Pattern: "/schools",
		Handler: apiv1.CreateSchool,
		Flag:    "Y",
		Desc:    "Create school",
	},
	{
		Method:  "PUT",
		Pattern: "/schools/:id",
		Handler: apiv1.UpdateSchool,
		Flag:    "Y",
		Desc:    "Update school",
	},
	{
		Method:  "DELETE",
		Pattern: "/schools/:id",
		Handler: apiv1.DeleteSchool,
		Flag:    "Y",
		Desc:    "Delete school",
	},
}

func init() {
	Routes = append(Routes, school...)
}
