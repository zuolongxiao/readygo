package models

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"readygo/pkg/settings"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type LocalTime time.Time

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format(settings.App.TimeFormat))), nil
}
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}
func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
func (t *LocalTime) String() string {
	tTime := time.Time(*t)
	return tTime.Format(settings.App.TimeFormat)
}

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
	CreatedAt LocalTime `json:"created_at"`
	UpdatedAt LocalTime `json:"updated_at"`
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

// KeyValueQueryer
type KeyValueQueryer struct {
	Entries map[string]string
}

// Query implements global.queryer
func (q *KeyValueQueryer) Query(k string) string {
	return q.Entries[k]
}
