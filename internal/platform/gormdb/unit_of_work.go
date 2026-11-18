package gormdb

import (
	"context"

	"inventario/internal/platform/uow"

	"gorm.io/gorm"
)

// UnitOfWork is the GORM implementation of uow.UnitOfWork. Do opens a
// transaction and hands fn a context carrying it; repositories resolve their
// handle from that context with Conn.
type UnitOfWork struct {
	db *gorm.DB
}

var _ uow.UnitOfWork = (*UnitOfWork)(nil)

// NewUnitOfWork returns a UnitOfWork that opens transactions on db.
func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

// Do runs fn inside a transaction. If ctx already carries one, fn joins it
// (no nested BEGIN or savepoint) and the outermost Do decides the outcome.
//
// A non-nil error from fn rolls back and is returned. A panic in fn rolls
// back and is re-raised. Otherwise the transaction is committed.
func (u *UnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := txFrom(ctx); ok {
		return fn(ctx)
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}
