package models

import (
	"database/sql"
	"readygo/pkg/settings"
	"reflect"
	"strconv"
	"strings"
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
	return int(settings.App.PageSize)
}

// BaseView base view
type BaseView struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IDsQueryer
type IDsQueryer struct {
	List interface{}
	Key  string
}

// Query implements global.queryer
func (q *IDsQueryer) Query(s string) string {
	if s != "IDs" {
		return ""
	}
	m := make(map[string]bool)
	var ids []string

	lst := reflect.ValueOf(q.List)
	for i := 0; i < lst.Len(); i++ {
		v := reflect.ValueOf(lst.Index(i).Interface())
		id := v.FieldByName(q.Key).Uint()
		if id <= 0 {
			continue
		}
		val := strconv.Itoa(int(id))
		if _, ok := m[val]; !ok {
			m[val] = true
			ids = append(ids, val)
		}
	}

	return strings.Join(ids, ",")
}
