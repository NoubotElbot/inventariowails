package gormdb

import (
	"context"

	"gorm.io/gorm"
)

// txKey is the context key under which UnitOfWork stores the active
// transaction handle.
type txKey struct{}

// txFrom returns the transaction carried by ctx, if any.
func txFrom(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok
}

// Conn returns the handle repositories must query through: the transaction
// carried by ctx (set by UnitOfWork.Do) or, outside a unit of work, db. The
// result is always a fresh session bound to ctx.
func Conn(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := txFrom(ctx); ok {
		return tx.WithContext(ctx)
	}
	return db.WithContext(ctx)
}
