package models

import (
	"readygo/pkg/errs"
	"readygo/pkg/global"

	"gorm.io/gorm"
)

type District struct {
	Base

	Name string `gorm:"type:string;size:50;index:uk_name,unique;not null"`
}

// DistrictView view
type DistrictView struct {
	BaseView

	Name string `json:"name"`
}

// DistrictCreate create binding
type DistrictCreate struct {
	Name *string `json:"name" binding:"required,min=2,max=50"`
}

// DistrictUpdate update binding
type DistrictUpdate struct {
	Name *string `json:"name" binding:"required,min=2,max=50"`
}

// BeforeSave hook
func (mdl *District) BeforeSave(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(mdl).Where("id <> ? AND name = ?", mdl.ID, mdl.Name).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}

	if count > 0 {
		return errs.DuplicatedError("区域已存在")
	}

	return nil
}

// BeforeDelete hook
func (mdl *District) BeforeDelete(tx *gorm.DB) error {
	var count int64

	if err := tx.Model(&Zone{}).Where("district_id = ?", mdl.ID).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.ReferenceRestrictError("district has been referenced by zone")
	}

	if err := tx.Model(&School{}).Where("district_id = ?", mdl.ID).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.ReferenceRestrictError("district has been referenced by school")
	}

	if err := tx.Model(&Community{}).Where("district_id = ?", mdl.ID).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.ReferenceRestrictError("district has been referenced by community")
	}

	return nil
}

// Filter filter
func (*District) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	return db
}
