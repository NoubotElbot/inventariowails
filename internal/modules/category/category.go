// Package category is the composition facade for the category module.
// It wires the module's own layers (domain, application, infrastructure,
// handler) and exposes the Wails-facing handler to the application root.
package category

import (
	"inventario/internal/modules/category/application"
	"inventario/internal/modules/category/handler"
	"inventario/internal/modules/category/infrastructure"

	"github.com/jmoiron/sqlx"
)

// Module holds the public surface of the category module.
type Module struct {
	Handler *handler.CategoryHandler
}

// New builds the category module's object graph from a shared DB handle.
func New(db *sqlx.DB) *Module {
	repo := infrastructure.NewCategoryRepository(db)

	h := handler.NewCategoryHandler(
		application.NewListUseCase(repo),
		application.NewListOptionsUseCase(repo),
		application.NewFindUseCase(repo),
		application.NewCreateUseCase(repo),
		application.NewUpdateUseCase(repo),
		application.NewDeleteUseCase(repo),
		application.NewRestoreUseCase(repo),
	)

	return &Module{Handler: h}
}
