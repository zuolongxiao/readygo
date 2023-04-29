package models

import (
	"readygo/pkg/errs"
	"readygo/pkg/global"

	"gorm.io/gorm"
)

// Permission model
type Permission struct {
	Base

	Name      string `gorm:"type:string;size:50;index:uk_name,unique;not null"`
	Title     string `gorm:"type:string;size:100;not null"`
	Group     string `gorm:"type:string;size:50;not null"`
	Note      string `gorm:"type:string;size:200;not null"`
	IsEnabled string `gorm:"type:char(1);default:N;not null"`
}

// PermissionView view
type PermissionView struct {
	BaseView

	Name      string `json:"name"`
	Title     string `json:"title"`
	Group     string `json:"group"`
	IsEnabled string `json:"is_enabled"`
}

// PermissionCreate binding
type PermissionCreate struct {
	Name      string `json:"name" binding:"required,max=50"`
	Title     string `json:"title" binding:"required,max=50"`
	Group     string `json:"group" binding:"max=50"`
	Note      string `json:"note" binding:"max=50"`
	IsEnabled string `json:"is_enabled" binding:"omitempty,oneof=N Y"`
}

// PermissionUpdate binding
type PermissionUpdate struct {
	Name      string `json:"name" binding:"max=50"`
	Title     string `json:"title" binding:"max=50"`
	Group     string `json:"group" binding:"max=50"`
	Note      string `json:"note" binding:"max=50"`
	IsEnabled string `json:"is_enabled" binding:"omitempty,oneof=N Y"`
}

// BeforeSave hook
func (mdl *Permission) BeforeSave(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(mdl).Where("id <> ? AND name = ?", mdl.ID, mdl.Name).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}

	if count > 0 {
		return errs.DuplicatedError("permission.name")
	}

	return nil
}

// BeforeDelete hook
func (mdl *Permission) BeforeDelete(tx *gorm.DB) error {
	var count int64

	if err := tx.Model(&Authorization{}).Where("permission_id = ?", mdl.ID).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.ReferenceRestrictError("permission has been referenced by authorization")
	}

	return nil
}

// Filter filter
func (*Permission) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", name+"%")
	}

	if title := c.Query("title"); title != "" {
		db = db.Where("title LIKE ?", title+"%")
	}

	if group := c.Query("group"); group != "" {
		db = db.Where("`group` = ?", group)
	}

	return db
}
