package api

import (
	"readygo/models"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// ListZones ListZones
func ListZones(c *gin.Context) {
	cw := utils.NewContextWrapper(c)
	svc := services.New(&models.Zone{})

	var zoneList []models.ZoneView
	if err := svc.Find(&zoneList, c); err != nil {
		cw.Respond(err, nil)
		return
	}

	districtIDsQueryer := models.IDsQueryer{
		List: zoneList,
		Key:  "DistrictID",
	}
	districtSvc := services.New(&models.District{})
	var districtList []models.DistrictView
	if err := districtSvc.Find(&districtList, &districtIDsQueryer); err != nil {
		cw.Respond(err, nil)
		return
	}

	districtNameDict := make(map[uint64]string)
	for _, district := range districtList {
		districtNameDict[district.ID] = district.Name
	}

	type List struct {
		models.ZoneView
		DistrictName string `json:"district_name"`
	}
	var lst []List
	for _, row := range zoneList {
		dst := List{}
		copier.CopyWithOption(&dst, row, copier.Option{IgnoreEmpty: true})
		dst.DistrictName = districtNameDict[row.DistrictID]
		lst = append(lst, dst)
	}

	resp := map[string]interface{}{
		"list": lst,
		"prev": svc.GetPrev(),
		"next": svc.GetNext(),
	}

	cw.Respond(nil, resp)
}

// CreateZone CreateZone
func CreateZone(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	var binding models.ZoneCreate
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Zone{}
	svc := services.New(&mdl)
	if err := svc.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Create(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	var view models.ZoneView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// UpdateZone UpdateZone
func UpdateZone(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	var binding models.ZoneUpdate
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Zone{}
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

	var view models.ZoneView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// DeleteZone DeleteZone
func DeleteZone(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	svc := services.New(&models.Zone{})
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
