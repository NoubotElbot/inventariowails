// Package uow defines the Unit of Work port used by application use cases to
// run several repository operations atomically. It has no infrastructure
// dependencies so use cases stay decoupled from the database driver.
package uow

import "context"

// UnitOfWork marks the boundaries of a transaction. Do runs fn with a context
// that carries the transaction; repositories called with that context take
// part in it. Repositories called with any other context do not.
//
// Returning a non-nil error from fn (or panicking) rolls the work back;
// returning nil commits it. Calling Do with a context that already carries a
// transaction joins it instead of starting a new one.
type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
