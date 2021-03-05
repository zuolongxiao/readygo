package models

import (
	"strconv"

	"github.com/zuolongxiao/readygo/pkg/errs"
	"github.com/zuolongxiao/readygo/pkg/global"
	"gorm.io/gorm"
)

// Authorization model
type Authorization struct {
	Base

	RoleID       uint64 `gorm:"type:uint;index:uk_role_permission,unique,priority:1;index:uk_permission_role,unique,priority:2;not null"`
	PermissionID uint64 `gorm:"type:uint;index:uk_role_permission,unique,priority:2;index:uk_permission_role,unique,priority:1;not null"`
}

// AuthorizationView view
type AuthorizationView struct {
	BaseView

	RoleID       uint64 `json:"role_id"`
	PermissionID uint64 `json:"permission_id"`
}

// AuthorizationBinding binding
type AuthorizationBinding struct {
	RoleID       uint64 `json:"role_id" binding:"min=0"`
	PermissionID uint64 `json:"permission_id" binding:"min=0"`
}

// AuthorizationPermission model
type AuthorizationPermission struct {
	PermissionID uint64 `json:"permission_id"`
}

// BeforeSave hook
func (m *Authorization) BeforeSave(tx *gorm.DB) error {
	var count int64

	if m.RoleID > 0 {
		if err := tx.Model(&Role{}).Where("id = ?", m.RoleID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count == 0 {
			return errs.ValidationError("role does not exist")
		}
	}
	if m.PermissionID > 0 {
		if err := tx.Model(&Permission{}).Where("id = ?", m.PermissionID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count == 0 {
			return errs.ValidationError("permission does not exist")
		}
	}

	if m.RoleID > 0 && m.PermissionID > 0 {
		if err := tx.Model(m).Where("role_id = ? AND permission_id = ?", m.RoleID, m.PermissionID).Limit(1).Count(&count).Error; err != nil {
			return errs.DBError(err.Error())
		}
		if count > 0 {
			return errs.ValidationError("authorization already existed")
		}
	}

	return nil
}

// Filter filter
func (*Authorization) Filter(db *gorm.DB, c global.Queryer) *gorm.DB {
	if roleID, _ := strconv.Atoi(c.Query("role_id")); roleID > 0 {
		db = db.Where("role_id = ?", roleID)
	}

	if permissionID, _ := strconv.Atoi(c.Query("permission_id")); permissionID > 0 {
		db = db.Where("permission_id = ?", permissionID)
	}

	return db
}
