package models

import (
	"readygo/pkg/errs"
	"readygo/pkg/global"

	"gorm.io/gorm"
)

// Role model
type Role struct {
	Base

	Name string `gorm:"type:string;size:50;index:uk_name,unique;not null"`
}

// RoleView view
type RoleView struct {
	BaseView

	Name string `json:"name"`
}

// RoleCreate binding
type RoleCreate struct {
	Name string `json:"name" binding:"required,max=50"`
}

// RoleUpdate binding
type RoleUpdate struct {
	Name string `json:"name" binding:"max=50"`
}

// BeforeSave hook
func (m *Role) BeforeSave(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(m).Where("id <> ? AND name = ?", m.ID, m.Name).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}

	if count > 0 {
		return errs.DuplicatedError("role.name")
	}

	return nil
}

// BeforeDelete hook
func (m *Role) BeforeDelete(tx *gorm.DB) error {
	var count int64

	if err := tx.Model(&Admin{}).Where("role_id = ?", m.ID).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.ReferenceRestrictError("role has been referenced by admin")
	}

	if err := tx.Model(&Authorization{}).Where("role_id = ?", m.ID).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.ReferenceRestrictError("role has been referenced by authorization")
	}

	return nil
}

// Filter filter
func (*Role) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", name+"%")
	}

	return db
}
