package api

import (
	"readygo/models"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// ListSchools ListSchools
func ListSchools(c *gin.Context) {
	cw := utils.NewContextWrapper(c)
	svc := services.New(&models.School{})

	var schoolList []models.SchoolView
	if err := svc.Find(&schoolList, c); err != nil {
		cw.Respond(err, nil)
		return
	}

	districtIDsQueryer := models.IDsQueryer{
		List: schoolList,
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

	zoneIDsQueryer := models.IDsQueryer{
		List: schoolList,
		Key:  "ZoneID",
	}
	zoneSvc := services.New(&models.Zone{})
	var zoneList []models.Zone
	if err := zoneSvc.Find(&zoneList, &zoneIDsQueryer); err != nil {
		cw.Respond(err, nil)
		return
	}
	zoneNameDict := make(map[uint64]string)
	for _, zone := range zoneList {
		zoneNameDict[zone.ID] = zone.Name
	}

	type List struct {
		models.SchoolView
		DistrictName string `json:"district_name"`
		ZoneName     string `json:"zone_name"`
	}
	var lst []List
	for _, row := range schoolList {
		dst := List{}
		copier.CopyWithOption(&dst, row, copier.Option{IgnoreEmpty: true})
		dst.DistrictName = districtNameDict[row.DistrictID]
		dst.ZoneName = zoneNameDict[row.ZoneID]
		lst = append(lst, dst)
	}

	resp := map[string]interface{}{
		"list": lst,
		"prev": svc.GetPrev(),
		"next": svc.GetNext(),
	}

	cw.Respond(nil, resp)
}

// CreateSchool CreateSchool
func CreateSchool(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	var binding models.SchoolCreate
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.School{}
	svc := services.New(&mdl)
	if err := svc.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Create(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	var view models.SchoolView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// UpdateSchool UpdateSchool
func UpdateSchool(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	var binding models.SchoolUpdate
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.School{}
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

	var view models.SchoolView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// DeleteSchool DeleteSchool
func DeleteSchool(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	svc := services.New(&models.School{})
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
