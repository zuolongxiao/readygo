package models

import (
	"database/sql"

	"github.com/zuolongxiao/readygo/pkg/errs"
	"github.com/zuolongxiao/readygo/pkg/global"
	"github.com/zuolongxiao/readygo/pkg/utils"
	"gorm.io/gorm"
)

// Admin model
type Admin struct {
	Base

	RoleID      uint64       `gorm:"type:uint;index:idx_role;not null"`
	Username    string       `gorm:"type:string;size:50;index:uk_username,unique;not null"`
	Password    string       `gorm:"type:char(60);;not null"`
	IsLocked    string       `gorm:"type:char(1);default:N;not null"`
	IPAddr      string       `gorm:"type:string;size:100;not null"`
	LastLoginIP string       `gorm:"type:string;size:100;not null"`
	LastLoginAt sql.NullTime `gorm:"type:timestamp"`
}

// AdminView view
type AdminView struct {
	BaseView

	RoleID   uint64 `json:"role_id"`
	Username string `json:"username"`
}

// AdminCreate binding
type AdminCreate struct {
	RoleID   uint64 `json:"role_id" binding:"min=0"`
	Username string `json:"username" binding:"required,alphanum,min=2,max=50"`
	Password string `json:"password" binding:"required,min=2,max=50"`
	IsLocked string `json:"is_locked" binding:"omitempty,oneof=N Y"`
}

// AdminUpdate binding
type AdminUpdate struct {
	RoleID   uint64 `json:"role_id" binding:"min=0"`
	Username string `json:"username" binding:"omitempty,alphanum,min=2,max=50"`
	Password string `json:"password" binding:"omitempty,min=2,max=50"`
	IsLocked string `json:"is_locked" binding:"omitempty,oneof=N Y"`
}

// Auth binding
type Auth struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=50"`
	Password string `json:"password" binding:"required,min=2,max=50"`
}

// ProfileView view
type ProfileView struct {
	BaseView

	Username string `json:"username"`
	Role     string `json:"role"`
}

// ProfileUpdate binding
type ProfileUpdate struct {
	Password        string `json:"password" binding:"required,eqfield=PasswordConfirm,min=2,max=50"`
	PasswordOld     string `json:"password_old" binding:"required,min=2,max=50"`
	PasswordConfirm string `json:"password_confirm"`
}

// BeforeSave hook
func (m *Admin) BeforeSave(tx *gorm.DB) error {
	var count int64

	if m.RoleID != 0 {
		if err := tx.Model(&Role{}).Where("id = ?", m.RoleID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count == 0 {
			return errs.ValidationError("role does not exist")
		}
	}

	if err := tx.Model(m).Where("id <> ? AND username = ?", m.ID, m.Username).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.DuplicatedError("admin.username")
	}

	if m.Password != "" {
		password, err := utils.HashPassword(m.Password)
		if err != nil {
			return err
		}
		m.Password = password
	}

	return nil
}

// Filter filter
func (*Admin) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if username := c.Query("username"); username != "" {
		db = db.Where("username LIKE ?", username+"%")
	}

	if isLocked := c.Query("is_locked"); isLocked != "" {
		db = db.Where("is_locked = ?", isLocked)
	}

	return db
}
