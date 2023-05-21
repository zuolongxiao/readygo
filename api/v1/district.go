package api

import (
	"readygo/models"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// ListDistricts list districts
func ListDistricts(ctx *gin.Context) {
	cw := utils.NewContextWrapper(ctx)
	svc := services.New(&models.District{})

	var list []models.DistrictView
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

// CreateDistrict create district
func CreateDistrict(ctx *gin.Context) {
	cw := utils.NewContextWrapper(ctx)

	binding := models.DistrictCreate{}
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	var mdl models.District
	svc := services.New(&mdl)
	if err := svc.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Create(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	var view models.DistrictView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// UpdateDistrict update district
func UpdateDistrict(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	binding := models.DistrictUpdate{}
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.District{}
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

	var view models.DistrictView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// DeleteDistrict delete district
func DeleteDistrict(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	svc := services.New(&models.District{})
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
