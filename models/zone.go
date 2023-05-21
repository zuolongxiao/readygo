package models

import (
	"readygo/pkg/errs"
	"readygo/pkg/global"

	"gorm.io/gorm"
)

type Zone struct {
	Base

	Name       string `gorm:"type:string;size:50;index:uk_name,unique;not null"`
	DistrictID uint64 `gorm:"type:uint;index:idx_district;not null"`
}

// ZoneView view
type ZoneView struct {
	BaseView

	Name       string `json:"name"`
	DistrictID uint64 `json:"district_id"`
}

// ZoneCreate create binding
type ZoneCreate struct {
	Name       *string `json:"name" binding:"required,min=2,max=50"`
	DistrictID *uint64 `json:"district_id" binding:"required"`
}

// ZoneUpdate update binding
type ZoneUpdate struct {
	Name       *string `json:"name" binding:"required,min=2,max=50"`
	DistrictID *uint64 `json:"district_id" binding:"required"`
}

// BeforeSave hook
func (mdl *Zone) BeforeSave(tx *gorm.DB) error {
	var count int64

	if mdl.DistrictID != 0 {
		if err := tx.Model(&District{}).Where("id = ?", mdl.DistrictID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count == 0 {
			return errs.ValidationError("区域不存在")
		}
	}
	if err := tx.Model(mdl).Where("id <> ? AND name = ?", mdl.ID, mdl.Name).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}

	if count > 0 {
		return errs.DuplicatedError("zone.name")
	}

	return nil
}

// BeforeDelete hook
func (mdl *Zone) BeforeDelete(tx *gorm.DB) error {
	var count int64

	if err := tx.Model(&School{}).Where("zone_id = ?", mdl.ID).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.ReferenceRestrictError("zone has been referenced by school")
	}

	if err := tx.Model(&Community{}).Where("zone_id = ?", mdl.ID).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.ReferenceRestrictError("zone has been referenced by community")
	}

	return nil
}

// Filter filter
func (*Zone) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if districtID := c.Query("district_id"); districtID != "" {
		db = db.Where("district_id = ?", districtID)
	}

	return db
}
