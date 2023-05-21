package api

import (
	"readygo/models"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// ListCommunities ListCommunities
func ListCommunities(c *gin.Context) {
	cw := utils.NewContextWrapper(c)
	svc := services.New(&models.Community{})

	var schoolList []models.CommunityView
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

	primaryIDsQueryer := models.IDsQueryer{
		List: schoolList,
		Key:  "PrimaryID",
	}
	primarySvc := services.New(&models.School{})
	var primaryList []models.School
	if err := primarySvc.Find(&primaryList, &primaryIDsQueryer); err != nil {
		cw.Respond(err, nil)
		return
	}
	primaryNameDict := make(map[uint64]string)
	for _, primary := range primaryList {
		primaryNameDict[primary.ID] = primary.Name
	}

	middleIDsQueryer := models.IDsQueryer{
		List: schoolList,
		Key:  "MiddleID",
	}
	middleSvc := services.New(&models.School{})
	var middleList []models.School
	if err := middleSvc.Find(&middleList, &middleIDsQueryer); err != nil {
		cw.Respond(err, nil)
		return
	}
	middleNameDict := make(map[uint64]string)
	for _, middle := range middleList {
		middleNameDict[middle.ID] = middle.Name
	}

	type List struct {
		models.CommunityView
		DistrictName string `json:"district_name"`
		ZoneName     string `json:"zone_name"`
		PrimaryName  string `json:"primary_name"`
		MiddleName   string `json:"middle_name"`
	}
	var lst []List
	for _, row := range schoolList {
		dst := List{}
		copier.CopyWithOption(&dst, row, copier.Option{IgnoreEmpty: true})
		dst.DistrictName = districtNameDict[row.DistrictID]
		dst.ZoneName = zoneNameDict[row.ZoneID]
		dst.PrimaryName = primaryNameDict[row.PrimaryID]
		dst.MiddleName = middleNameDict[row.MiddleID]
		lst = append(lst, dst)
	}

	resp := map[string]interface{}{
		"list": lst,
		"prev": svc.GetPrev(),
		"next": svc.GetNext(),
	}

	cw.Respond(nil, resp)
}

// CreateCommunity CreateCommunity
func CreateCommunity(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	var binding models.CommunityCreate
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Community{}
	svc := services.New(&mdl)
	if err := svc.Fill(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	if err := svc.Create(cw); err != nil {
		cw.Respond(err, nil)
		return
	}

	var view models.CommunityView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// UpdateCommunity UpdateCommunity
func UpdateCommunity(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	var binding models.CommunityUpdate
	if err := cw.Bind(&binding); err != nil {
		cw.Respond(err, nil)
		return
	}

	mdl := models.Community{}
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

	var view models.CommunityView
	copier.Copy(&view, &mdl)

	cw.Respond(nil, view)
}

// DeleteCommunity DeleteCommunity
func DeleteCommunity(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	svc := services.New(&models.Community{})
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
