package models

import (
	"readygo/pkg/errs"
	"readygo/pkg/global"

	"gorm.io/gorm"
)

type School struct {
	Base

	DistrictID uint64  `gorm:"type:uint;index:idx_district;not null"`
	ZoneID     uint64  `gorm:"type:uint;index:idx_zone;not null"`
	Name       string  `gorm:"type:string;size:50;index:uk_name,unique;not null"`
	Addr       string  `gorm:"type:string;size:250;not null"`
	Type       int     `gorm:"type:int;not null"`
	Rank       int     `gorm:"type:int;not null"`
	Lat        float64 `gorm:"type:decimal(10,8)"`
	Lng        float64 `gorm:"type:decimal(11,8)"`
	Note       string  `gorm:"type:string;size:250;not null"`
}

// SchoolView view
type SchoolView struct {
	BaseView

	DistrictID uint64  `json:"district_id"`
	ZoneID     uint64  `json:"zone_id"`
	Name       string  `json:"name"`
	Addr       string  `json:"addr"`
	Type       int     `json:"type"`
	Rank       int     `json:"rank"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Note       string  `json:"note"`
}

// SchoolCreate create binding
type SchoolCreate struct {
	DistrictID *uint64  `json:"district_id" binding:"required"`
	ZoneID     *uint64  `json:"zone_id" binding:"required"`
	Name       *string  `json:"name" binding:"required,min=2,max=50"`
	Addr       *string  `json:"addr" binding:"required,min=0,max=250"`
	Type       *int     `json:"type" binding:"required,oneof=1 2"` // 1: 小学(primary), 2: 初中(middle)
	Rank       *int     `json:"rank" binding:"required"`
	Lat        *float64 `json:"lat" binding:"required"`
	Lng        *float64 `json:"lng" binding:"required"`
	Note       *string  `json:"note" binding:"required,max=250"`
}

// SchoolUpdate update binding
type SchoolUpdate struct {
	DistrictID *uint64  `json:"district_id" binding:"required"`
	ZoneID     *uint64  `json:"zone_id" binding:"required"`
	Name       *string  `json:"name" binding:"required,min=2,max=50"`
	Addr       *string  `json:"addr" binding:"required,min=0,max=250"`
	Type       *int     `json:"type" binding:"required,oneof=1 2"` // 1: 小学(primary), 2: 初中(middle)
	Rank       *int     `json:"rank" binding:"required"`
	Lat        *float64 `json:"lat" binding:"required"`
	Lng        *float64 `json:"lng" binding:"required"`
	Note       *string  `json:"note" binding:"required,max=250"`
}

// BeforeSave hook
func (mdl *School) BeforeSave(tx *gorm.DB) error {
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

	if err := tx.Model(mdl).Where("id <> ? AND name = ?", mdl.ID, mdl.Name).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}

	if count > 0 {
		return errs.DuplicatedError("school.name")
	}

	return nil
}

// BeforeDelete hook
func (mdl *School) BeforeDelete(tx *gorm.DB) error {
	return nil
}

// Filter filter
func (*School) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if districtID := c.Query("district_id"); districtID != "" {
		db = db.Where("district_id = ?", districtID)
	}

	if zoneID := c.Query("zone_id"); zoneID != "" {
		db = db.Where("zone_id = ?", zoneID)
	}

	if typ := c.Query("type"); typ != "" {
		db = db.Where("`type` = ?", typ)
	}

	return db
}
