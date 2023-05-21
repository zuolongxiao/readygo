package models

import (
	"database/sql"

	"readygo/pkg/errs"
	"readygo/pkg/global"
	"readygo/utils"

	"gorm.io/gorm"
)

// Admin model
type Admin struct {
	Base

	RoleID      uint64       `gorm:"type:uint;index:idx_role;not null"`
	Username    string       `gorm:"type:string;size:50;index:uk_username,unique;not null"`
	Password    string       `gorm:"type:char(60);not null"`
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
	IsLocked string `json:"is_locked"`
}

// AdminCreate binding
type AdminCreate struct {
	RoleID   *uint64 `json:"role_id" binding:"required"`
	Username *string `json:"username" binding:"required,alphanum,min=2,max=50"`
	Password *string `json:"password" binding:"required,min=2,max=50"`
	IsLocked *string `json:"is_locked" binding:"required,oneof=N Y"`
}

// AdminUpdate binding
type AdminUpdate struct {
	RoleID   *uint64 `json:"role_id" binding:"required"`
	Username *string `json:"username" binding:"required,alphanum,min=2,max=50"`
	Password *string `json:"password" binding:"required,min=0,max=50"`
	IsLocked *string `json:"is_locked" binding:"required,oneof=N Y"`
}

// Auth binding
type Auth struct {
	Username    *string `json:"username" binding:"required,alphanum,min=2,max=50"`
	Password    *string `json:"password" binding:"required,min=2,max=50"`
	CaptchaCode *string `json:"code" binding:"required,min=2,max=50"`
	CaptchaID   *string `json:"id" binding:"required,min=2,max=50"`
}

// ProfileView view
type ProfileView struct {
	BaseView

	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// ProfileUpdate binding
type ProfileUpdate struct {
	Password        *string `json:"password" binding:"required,eqfield=PasswordConfirm,min=2,max=50"`
	PasswordOld     *string `json:"password_old" binding:"required,min=2,max=50"`
	PasswordConfirm *string `json:"password_confirm"`
}

// BeforeSave hook
func (mdl *Admin) BeforeSave(tx *gorm.DB) error {
	var count int64

	if mdl.RoleID != 0 {
		if err := tx.Model(&Role{}).Where("id = ?", mdl.RoleID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count == 0 {
			return errs.ValidationError("role does not exist")
		}
	}

	if err := tx.Model(mdl).Where("id <> ? AND username = ?", mdl.ID, mdl.Username).Limit(1).Count(&count).Error; err != nil {
		return errs.DBError(err.Error())
	}
	if count > 0 {
		return errs.DuplicatedError("admin.username")
	}

	if mdl.Password == "" {
		tx.Statement.Omit("password")
	} else {
		hashedPassword, err := utils.HashPassword(mdl.Password)
		if err != nil {
			return errs.InternalServerError(err.Error())
		}
		mdl.Password = hashedPassword
	}

	return nil
}

// Filter filter
func (mdl *Admin) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if username := c.Query("username"); username != "" {
		db = db.Where("username LIKE ?", username+"%")
	}

	if isLocked := c.Query("is_locked"); isLocked != "" {
		db = db.Where("is_locked = ?", isLocked)
	}

	if roleID := c.Query("role_id"); roleID != "" {
		db = db.Where("role_id = ?", roleID)
	}

	return db
}
