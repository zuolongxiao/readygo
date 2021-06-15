package models

import (
	"database/sql"
	"time"
)

// Base base model
type Base struct {
	ID        uint64       `gorm:"type:uint;primaryKey"`
	CreatedAt sql.NullTime `gorm:"type:timestamp"`
	CreatedBy string       `gorm:"type:string;size:100;not null"`
	UpdatedAt sql.NullTime `gorm:"type:timestamp"`
	UpdatedBy string       `gorm:"type:string;size:100;not null"`
}

// Size implements services.Pager
func (*Base) Size() int {
	return 20
}

// BaseView base view
type BaseView struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
