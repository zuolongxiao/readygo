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
func New(mdl interface{}) *Base {
	return &Base{model: mdl}
}

// Find Find
func (svc *Base) Find(o interface{}, c global.Queryer) error {
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
	if p, ok := svc.model.(Pager); ok {
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

	if f, ok := svc.model.(Filterer); ok {
		session = f.Filter(session, c)
	}

	minIdKey := "MIN(id)"
	maxIdKey := "MAX(id)"
	result := map[string]interface{}{}
	session.Model(svc.model).Select(fmt.Sprintf("%s, %s", minIdKey, maxIdKey)).Take(result)
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

	err := session.Model(svc.model).Select("*").Order(order).Limit(size).Find(o).Error
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
			svc.next = 0
			svc.prev = 0
		} else {
			if sym == "+" {
				if last < maxId {
					svc.next = last
				}
				if first > minId {
					svc.prev = first
				}
			} else {
				if last > minId {
					svc.next = last
				}
				if first < maxId {
					svc.prev = first
				}
			}
		}
	}

	return nil
}

// Get get row
func (svc *Base) Get(o interface{}, cond interface{}) error {
	err := db.DB.Model(svc.model).Where(cond).First(o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := svc.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// GetRows get rows
func (svc *Base) GetRows(o interface{}, cond interface{}) error {
	err := db.DB.Model(svc.model).Where(cond).Find(o).Error
	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// GetByID GetByID
func (svc *Base) GetByID(o interface{}, qs string) error {
	id, _ := strconv.ParseUint(qs, 10, 0)
	if id <= 0 {
		return errs.ValidationError("invalid ID")
	}

	err := db.DB.Model(svc.model).Where("id = ?", id).First(o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := svc.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// GetByKey GetByKey
func (svc *Base) GetByKey(o interface{}, k string, v interface{}) error {
	where := k + " = ?"
	err := db.DB.Model(svc.model).Where(where, v).First(o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := svc.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// Load get and fill the model
func (svc *Base) Load() error {
	err := db.DB.Model(svc.model).Where(svc.model).Take(svc.model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := svc.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// LoadByID get and fill the model
func (svc *Base) LoadByID(qs string) error {
	id, _ := strconv.ParseUint(qs, 10, 0)
	if id <= 0 {
		return errs.ValidationError("invalid ID")
	}

	err := db.DB.Model(svc.model).Where("id = ?", id).First(svc.model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := svc.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// LoadByKey LoadByKey
func (svc *Base) LoadByKey(key string, v interface{}) error {
	where := key + " = ?"
	err := db.DB.Model(svc.model).Where(where, v).First(svc.model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := svc.getTypeName()
		return errs.NotFoundError(typ)
	}

	if err != nil {
		return errs.DBError(err.Error())
	}

	return nil
}

// Fill fill model
func (svc *Base) Fill(o interface{}) error {
	return copier.CopyWithOption(svc.model, o, copier.Option{IgnoreEmpty: false})
}

// Create Create
func (svc *Base) Create(cw global.IContextWrapper) error {
	mdl := reflect.ValueOf(svc.model).Elem()
	mdl.FieldByName("CreatedBy").SetString(cw.GetUsername())
	return db.DB.Create(svc.model).Error
}

// Save Save
func (svc *Base) Save(cw global.IContextWrapper) error {
	mdl := reflect.ValueOf(svc.model).Elem()
	mdl.FieldByName("UpdatedBy").SetString(cw.GetUsername())
	return db.DB.Select("*").Updates(svc.model).Error
}

// Update Update
func (svc *Base) Update(cw global.IContextWrapper, query interface{}, args ...interface{}) error {
	mdl := reflect.ValueOf(svc.model).Elem()
	mdl.FieldByName("UpdatedBy").SetString(cw.GetUsername())
	return db.DB.Select(query, args...).UpdateColumns(svc.model).Error
}

// Delete delete from database and destroy the model
func (svc *Base) Delete() error {
	err := db.DB.Delete(svc.model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		typ := svc.getTypeName()
		return errs.NotFoundError(typ)
	}

	return err
}

// GetOffset GetOffset
func (svc *Base) GetOffset() (uint64, uint64) {
	return svc.prev, svc.next
}

// GetPrev GetPrev
func (svc *Base) GetPrev() uint64 {
	return svc.prev
}

// GetNext GetNext
func (svc *Base) GetNext() uint64 {
	return svc.next
}

func (svc *Base) getTypeName() string {
	v := reflect.ValueOf(svc.model)
	typ := v.Elem().Type().String()
	ss := strings.Split(typ, ".")
	return strings.ToLower(ss[len(ss)-1:][0])
}

func reverse(slc interface{}) {
	v := reflect.ValueOf(slc).Elem()
	swap := reflect.Swapper(v.Interface())
	n := v.Len()
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
