package transaction

import (
	"context"
	"gorm.io/gorm"
)

type contextKey string

var tnxKey contextKey = "tnxKey"

func GetTnx(ctx context.Context) interface{} {
	return ctx.Value(tnxKey)
}

type TnxManager interface {
	Start(ctx context.Context) (Transaction, context.Context)
}

type Transaction interface {
	Commit() error
	RollBack()
	Finish(err error) error
}

type transaction struct {
	tnx *gorm.DB
}

func (t transaction) Commit() error {
	return t.tnx.Commit().Error
}

func (t transaction) RollBack() {
	t.tnx.Rollback()
}
func (t transaction) Finish(err error) error {
	if r := recover(); r != nil {
		t.RollBack()
		panic(r)
	}
	if err != nil {
		t.RollBack()
		return err
	}
	return t.tnx.Commit().Error
}

type tnxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) TnxManager {
	return tnxManager{db: db}
}

func (t tnxManager) Start(ctx context.Context) (Transaction, context.Context) {
	if GetTnx(ctx) != nil {
		return DummyTransaction{}, ctx
	}
	tnx := t.db.Begin()
	return transaction{tnx: tnx}, context.WithValue(ctx, tnxKey, tnx)
}

// Dummy Transaction - Support Multiple Level Transaction
type DummyTransaction struct {
}

func (d DummyTransaction) Commit() error {
	return nil
}

func (d DummyTransaction) RollBack() {

}
func (d DummyTransaction) Finish(err error) error {
	return err
}
