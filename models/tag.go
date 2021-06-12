package models

import (
	"readygo/pkg/errs"
	"readygo/pkg/global"

	"gorm.io/gorm"
)

// Tag model
type Tag struct {
	Base

	Name  string `gorm:"type:string;size:50;index:uk_name,unique;not null"`
	State string `gorm:"type:char(1);default:N;not null"`
}

// TagView view
type TagView struct {
	BaseView

	Name  string `json:"name"`
	State string `json:"state"`
}

// TagCreate binding
type TagCreate struct {
	Name  string `json:"name" binding:"required,max=50"`
	State string `json:"state" binding:"omitempty,oneof=N Y"`
}

// TagUpdate binding
type TagUpdate struct {
	Name  string `json:"name" binding:"max=50"`
	State string `json:"state" binding:"omitempty,oneof=N Y"`
}

// BeforeSave hook
func (m *Tag) BeforeSave(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(m).Where("id <> ? AND name = ?", m.ID, m.Name).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}

	if count > 0 {
		return errs.DuplicatedError("tag.name")
	}

	return nil
}

// Filter filter
func (*Tag) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if name := c.Query("name"); name != "" {
		db = db.Where("name LIKE ?", name+"%")
	}

	if state := c.Query("state"); state != "" {
		db = db.Where("state = ?", state)
	}

	return db
}
