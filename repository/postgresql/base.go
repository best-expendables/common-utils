package postgresql

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/best-expendables/common-utils/model"
	"github.com/best-expendables/common-utils/repository"
	"github.com/best-expendables/common-utils/repository/filter"
	"github.com/best-expendables/common-utils/transaction"
	nrcontext "github.com/best-expendables/newrelic-context"
	"github.com/jinzhu/gorm"
)

type BaseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) *BaseRepo {
	return &BaseRepo{
		db: db,
	}
}

func (r *BaseRepo) GetDB(ctx context.Context) *gorm.DB {
	db := r.db
	if tnx := transaction.GetTnx(ctx); tnx != nil {
		db = tnx.(*gorm.DB)
	}
	db = nrcontext.SetTxnToGorm(ctx, db)
	return db
}

func (r *BaseRepo) FindByID(ctx context.Context, m model.Model, id string, preloadFields ...string) error {
	q := r.GetDB(ctx)

	for _, p := range preloadFields {
		q = q.Preload(p)
	}

	err := q.Where("id = ?", id).Take(m).Error

	if err == gorm.ErrRecordNotFound {
		return repository.RecordNotFound
	}
	return err
}

func (r *BaseRepo) CreateOrUpdate(ctx context.Context, m model.Model, query interface{}, attrs ...interface{}) error {
	return r.GetDB(ctx).Where(query).Assign(attrs...).FirstOrCreate(m).Error
}

func (r *BaseRepo) Update(ctx context.Context, m model.Model, attrs ...interface{}) error {
	return r.GetDB(ctx).Model(m).Update(attrs...).Error
}

func (r *BaseRepo) Updates(ctx context.Context, m model.Model, params interface{}) error {
	return r.GetDB(ctx).Model(m).Updates(params).Error
}

func (r *BaseRepo) Create(ctx context.Context, m model.Model) error {
	return r.GetDB(ctx).Create(m).Error
}

func (r *BaseRepo) Search(ctx context.Context, val interface{}, f filter.Filter, preloadFields ...string) error {
	q := r.GetDB(ctx).Model(val)
	for query, val := range f.GetWhere() {
		q = q.Where(query, val...)
	}

	for query, val := range f.GetJoins() {
		q = q.Joins(query, val...)
	}

	if f.GetGroups() != "" {
		q = q.Group(f.GetGroups())
	}

	if f.GetLimit() > 0 {
		q = q.Limit(f.GetLimit())
	}

	if len(f.GetOrderBy()) > 0 {
		for _, order := range f.GetOrderBy() {
			q = q.Order(order)
		}
	}

	for _, p := range preloadFields {
		q = q.Preload(p)
	}

	return q.Offset(f.GetOffset()).Find(val).Error
}

func (r *BaseRepo) Save(ctx context.Context, m model.Model) error {
	return r.GetDB(ctx).Model(m).Save(m).Error
}

func (r *BaseRepo) DeleteByID(ctx context.Context, m model.Model, id string) error {
	db := r.GetDB(ctx).Where("id = ?", id).Take(m)
	if db.Error != nil || m.GetId() == "" {
		return repository.RecordNotFound
	}
	return db.Delete(m).Error
}

func (r *BaseRepo) BulkCreate(ctx context.Context, arr []model.Model) error {
	if len(arr) == 0 {
		return nil
	}

	var valueStrings []string
	var valueArgs []interface{}
	properties := getStructProperties(arr[0])
	for _, val := range arr {
		_ = val.BeforeCreate(r.GetDB(ctx).NewScope(val))
		ri := redirectReflectPtrToElem(reflect.ValueOf(val))

		var valueKeys []string
		for _, property := range properties {
			valueKeys = append(valueKeys, "?")
			valueArgs = append(valueArgs, ri.FieldByName(property).Interface())
		}
		valueStrings = append(valueStrings, strings.Join(valueKeys, ","))
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.GetDB(ctx).NewScope(arr[0]).TableName(),
		strings.Join(transformPropertiesToFieldNames(properties), ","),
		strings.Join(valueStrings, "),("))

	return r.GetDB(ctx).Exec(sql, valueArgs...).Error
}

func transformPropertiesToFieldNames(properties []string) []string {
	var fieldNames []string
	for _, property := range properties {
		fieldNames = append(fieldNames, gorm.ToDBName(property))
	}

	return fieldNames
}

func getStructProperties(val interface{}) []string {
	var fields []string
	ri := redirectReflectPtrToElem(reflect.ValueOf(val))
	ri.FieldByNameFunc(func(name string) bool {
		if (ri.FieldByName(name).Kind() == reflect.Slice) ||
			((ri.FieldByName(name).Kind() == reflect.Struct) && (reflect.TypeOf(ri.FieldByName(name).Interface()).String() != "time.Time")) {
			return false
		}

		fields = append(fields, name)
		return false
	})

	return fields
}

func redirectReflectPtrToElem(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}
