package models

import (
	"readygo/pkg/errs"
	"readygo/pkg/global"

	"gorm.io/gorm"
)

// Community model
type Community struct {
	Base

	DistrictID uint64 `gorm:"type:uint;index:idx_district;not null"`
	ZoneID     uint64 `gorm:"type:uint;index:idx_zone;not null"`
	PrimaryID  uint64 `gorm:"type:uint;index:idx_primary;not null"`
	MiddleID   uint64 `gorm:"type:uint;index:idx_middle;not null"`
	Name       string `gorm:"type:string;size:50;index:uk_name,unique;not null"`
	Addr       string `gorm:"type:string;size:250;not null"`
	BuildYear  string `gorm:"type:string;size:50;not null"`
	HouseNum   int    `gorm:"type:int;not null"`
	LotNum     int    `gorm:"type:int;not null"`
	Note       string `gorm:"type:string;size:250;not null"`
}

// CommunityView view
type CommunityView struct {
	BaseView

	DistrictID uint64 `json:"district_id"`
	ZoneID     uint64 `json:"zone_id"`
	PrimaryID  uint64 `json:"primary_id"`
	MiddleID   uint64 `json:"middle_id"`
	Name       string `json:"name"`
	Addr       string `json:"addr"`
	BuildYear  string `json:"build_year"`
	HouseNum   int    `json:"house_num"`
	LotNum     int    `json:"lot_num"`
	Note       string `json:"note"`
}

// CommunityCreate binding
type CommunityCreate struct {
	DistrictID *uint64 `json:"district_id" binding:"required"`
	ZoneID     *uint64 `json:"zone_id" binding:"required"`
	PrimaryID  *uint64 `json:"primary_id" binding:"required"`
	MiddleID   *uint64 `json:"middle_id" binding:"required"`
	Name       *string `json:"name" binding:"required,min=2,max=50"`
	Addr       *string `json:"addr" binding:"required,min=0,max=250"`
	BuildYear  *string `json:"build_year" binding:"required,min=0,max=20"`
	HouseNum   *int    `json:"house_num" binding:"required,min=0,max=20000"`
	LotNum     *int    `json:"lot_num" binding:"required,min=0,max=20000"`
	Note       *string `json:"note" binding:"required,max=250"`
}

// CommunityUpdate binding
type CommunityUpdate struct {
	DistrictID *uint64 `json:"district_id" binding:"required"`
	ZoneID     *uint64 `json:"zone_id" binding:"required"`
	PrimaryID  *uint64 `json:"primary_id" binding:"required"`
	MiddleID   *uint64 `json:"middle_id" binding:"required"`
	Name       *string `json:"name" binding:"required,min=2,max=50"`
	Addr       *string `json:"addr" binding:"required,min=0,max=250"`
	BuildYear  *string `json:"build_year" binding:"required,min=0,max=20"`
	HouseNum   *int    `json:"house_num" binding:"required,min=0,max=20000"`
	LotNum     *int    `json:"lot_num" binding:"required,min=0,max=20000"`
	Note       *string `json:"note" binding:"required,max=250"`
}

// BeforeSave hook
func (mdl *Community) BeforeSave(tx *gorm.DB) error {
	var count int64
	if mdl.DistrictID != 0 {
		if err := tx.Model(&District{}).Where("id = ?", mdl.DistrictID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count == 0 {
			return errs.ValidationError("区域不存在")
		}
	}
	if mdl.ZoneID != 0 {
		if err := tx.Model(&Zone{}).Where("id = ?", mdl.ZoneID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count == 0 {
			return errs.ValidationError("商圈不存在")
		}
	}
	if mdl.PrimaryID != 0 {
		if err := tx.Model(&School{}).Where("id = ?", mdl.PrimaryID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count == 0 {
			return errs.ValidationError("小学不存在")
		}
	}
	if mdl.MiddleID != 0 {
		if err := tx.Model(&School{}).Where("id = ?", mdl.MiddleID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count == 0 {
			return errs.ValidationError("初中不存在")
		}
	}

	if err := tx.Model(mdl).Where("id <> ? AND name = ?", mdl.ID, mdl.Name).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.DuplicatedError("小区已存在")
	}

	return nil
}

// Filter filter
func (mdl *Community) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if districtID := c.Query("district_id"); districtID != "" {
		db = db.Where("district_id = ?", districtID)
	}

	if zoneID := c.Query("zone_id"); zoneID != "" {
		db = db.Where("zone_id = ?", zoneID)
	}

	if primaryID := c.Query("primary_id"); primaryID != "" {
		db = db.Where("primary_id = ?", primaryID)
	}

	if middleID := c.Query("middle_id"); middleID != "" {
		db = db.Where("middle_id = ?", middleID)
	}

	return db
}
