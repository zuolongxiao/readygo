package services

import (
	"errors"
	"fmt"
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
	model interface{}
	prev  uint64
	next  uint64
}

// New new service
func New(m interface{}) *Base {
	return &Base{model: m}
}

// Find Find
func (s *Base) Find(o interface{}, c global.Queryer) error {
	sort := c.Query("sort") // +id or -id
	dir := c.Query("dir")   // next or prev
	size, _ := strconv.Atoi(c.Query("size"))
	offset, _ := strconv.ParseUint(c.Query("offset"), 10, 0)
	IDs := c.Query("IDs") // 1,2,3

	fields := []string{"id"}
	maxSize := 500

	field := ""
	sortDir := ""
	op := ""
	sym := "+"

	if dir != "next" && dir != "prev" {
		dir = "next"
	}

	if len(sort) > 0 {
		ch := sort[0:1]
		if ch == "+" || ch == "-" {
			sym = ch
			field = sort[1:]
		} else {
			field = sort
		}
	}

	if !utils.StrInSlice(field, fields) {
		field = fields[0]
	}

	isReverse := false
	if sym == "+" {
		if dir == "next" {
			op = ">"
			sortDir = "ASC"
		} else {
			op = "<"
			sortDir = "DESC"
			isReverse = true
		}
	} else {
		if dir == "next" {
			op = "<"
			sortDir = "DESC"
		} else {
			op = ">"
			sortDir = "ASC"
			isReverse = true
		}
	}

	order := field + " " + sortDir

	pageSize := 20
	if p, ok := s.model.(Pager); ok {
		pageSize = p.Size()
	}

	if size <= 0 {
		size = pageSize
	}
	if size > maxSize {
		size = maxSize
	}

	session := db.DB.Session(&gorm.Session{})

	if IDs != "" {
		ids := strings.Split(IDs, ",")
		session = session.Where("id IN ?", ids)
	}

	if f, ok := s.model.(Filterer); ok {
		session = f.Filter(session, c)
	}

	minIdKey := "MIN(id)"
	maxIdKey := "MAX(id)"
	result := map[string]interface{}{}
	session.Model(s.model).Select(fmt.Sprintf("%s, %s", minIdKey, maxIdKey)).Take(result)
	var minId uint64
	var maxId uint64
	if reflect.ValueOf(result[minIdKey]).IsValid() {
		minId = uint64(reflect.ValueOf(result[minIdKey]).Int())
	}
	if reflect.ValueOf(result[maxIdKey]).IsValid() {
		maxId = uint64(reflect.ValueOf(result[maxIdKey]).Int())
	}
	// fmt.Println(minId, maxId)

	if offset > 0 {
		where := "id " + op + " ?"
		session = session.Where(where, offset)
	}

	err := session.Model(s.model).Select("*").Order(order).Limit(size).Find(o).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// return errs.NotFoundError(err.Error())
			return nil
		}
		return errs.DBError(err.Error())
	}

	if isReverse {
		reverse(o)
	}

	v := reflect.ValueOf(o).Elem()
	len := v.Len()
	if len > 0 {
		first := v.Index(0).FieldByName("ID").Uint()
		last := v.Index(len - 1).FieldByName("ID").Uint()
		if offset == 0 && len < size {
			s.next = 0
			s.prev = 0
		} else {
			if sym == "+" {
				if last < maxId {
					s.next = last
				}
				if first > minId {
					s.prev = first
				}
			} else {
				if last > minId {
					s.next = last
				}
				if first < maxId {
					s.prev = first
				}
			}
		}
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
func (s *Base) GetOffset() (uint64, uint64) {
	return s.prev, s.next
}

// GetPrev GetPrev
func (s *Base) GetPrev() uint64 {
	return s.prev
}

// GetNext GetNext
func (s *Base) GetNext() uint64 {
	return s.next
}

func (s *Base) getTypeName() string {
	v := reflect.ValueOf(s.model)
	typ := v.Elem().Type().String()
	ss := strings.Split(typ, ".")
	return strings.ToLower(ss[len(ss)-1:][0])
}

func reverse(s interface{}) {
	v := reflect.ValueOf(s).Elem()
	swap := reflect.Swapper(v.Interface())
	n := v.Len()
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
