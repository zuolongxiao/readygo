package services

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"readygo/pkg/db"
	"readygo/pkg/errs"
	"readygo/pkg/global"
	"readygo/utils"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// Pager interface
type Pager interface {
	Size() int
}

// Filterer interface
type Filterer interface {
	Filter(*gorm.DB, global.Queryer) *gorm.DB
}

// Base base service
type Base struct {
	model  interface{}
	offset int
}

// New new service
func New(m interface{}) *Base {
	return &Base{model: m}
}

// Find Find
func (s *Base) Find(o interface{}, c global.Queryer) error {
	fields := []string{"id"}
	maxSize := 100

	var field string
	dir := "ASC"
	op := ">"
	sort := c.Query("sort")
	if len(sort) > 0 {
		sym := sort[0:1]
		if sym == "+" || sym == "-" {
			if sym == "-" {
				dir = "DESC"
				op = "<"
			}
			field = sort[1:]
		} else {
			field = sort
		}
	}
	if !utils.StrInSlice(field, fields) {
		field = fields[0]
	}
	order := field + " " + dir

	offset, _ := strconv.Atoi(c.Query("offset"))
	if offset < 0 {
		offset = 0
	}

	pageSize := 20
	if p, ok := s.model.(Pager); ok {
		pageSize = p.Size()
	}
	size, _ := strconv.Atoi(c.Query("size"))
	if size <= 0 {
		size = pageSize
	}
	if size > maxSize {
		size = maxSize
	}

	db := db.DB.Session(&gorm.Session{})

	if offset > 0 {
		where := "id " + op + " ?"
		db = db.Where(where, offset)
	}

	if IDs := c.Query("IDs"); IDs != "" {
		ids := strings.Split(IDs, ",")
		db = db.Where("id IN ?", ids)
	}

	if f, ok := s.model.(Filterer); ok {
		db = f.Filter(db, c)
	}

	err := db.Model(s.model).Order(order).Limit(size).Find(o).Error
	if err != nil {
		return err
	}

	v := reflect.ValueOf(o).Elem()
	l := v.Len()
	if l == size {
		e := v.Index(l - 1)
		s.offset = int(e.FieldByName("ID").Uint())
	}

	return nil
}

// Get get row
func (s *Base) Get(o interface{}, cond interface{}) error {
	err := db.DB.Model(s.model).Where(cond).First(o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := s.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// GetRows get rows
func (s *Base) GetRows(o interface{}, cond interface{}) error {
	err := db.DB.Model(s.model).Where(cond).Find(o).Error
	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// GetByID GetByID
func (s *Base) GetByID(o interface{}, vs string) error {
	id, _ := strconv.Atoi(vs)
	if id <= 0 {
		return errs.ValidationError("invalid ID")
	}

	err := db.DB.Model(s.model).Where("id = ?", id).First(o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := s.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// GetByKey GetByKey
func (s *Base) GetByKey(o interface{}, k string, v interface{}) error {
	where := k + " = ?"
	err := db.DB.Model(s.model).Where(where, v).First(o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := s.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// Load get and fill the model
func (s *Base) Load() error {
	err := db.DB.Model(s.model).Where(s.model).Take(s.model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := s.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// LoadByID get and fill the model
func (s *Base) LoadByID(vs string) error {
	id, _ := strconv.Atoi(vs)
	if id <= 0 {
		return errs.ValidationError("invalid ID")
	}

	err := db.DB.Model(s.model).Where("id = ?", id).First(s.model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := s.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// LoadByKey LoadByKey
func (s *Base) LoadByKey(key string, v interface{}) error {
	where := key + " = ?"
	err := db.DB.Model(s.model).Where(where, v).First(s.model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := s.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// Fill fill model
func (s *Base) Fill(o interface{}) error {
	return copier.CopyWithOption(s.model, o, copier.Option{IgnoreEmpty: true})
}

// Create Create
func (s *Base) Create() error {
	return db.DB.Create(s.model).Error
}

// Save Save
func (s *Base) Save() error {
	return db.DB.Updates(s.model).Error
}

// Update Update
func (s *Base) Update(query interface{}, args ...interface{}) error {
	return db.DB.Select(query, args...).UpdateColumns(s.model).Error
}

// Delete delete from database and destroy the model
func (s *Base) Delete() error {
	err := db.DB.Delete(s.model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := s.getTypeName()
		return errs.NotFoundError(typ)
	}

	return err
}

// GetOffset GetOffset
func (s *Base) GetOffset() int {
	return s.offset
}

func (s *Base) getTypeName() string {
	v := reflect.ValueOf(s.model)
	typ := v.Elem().Type().String()
	ss := strings.Split(typ, ".")
	return strings.ToLower(ss[len(ss)-1:][0])
}
