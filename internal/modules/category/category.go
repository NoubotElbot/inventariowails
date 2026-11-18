// Package category is the composition root for the category module.
// It wires the hexagonal layers: the GORM repository adapter (driven) is
// injected into the core service through the outbound port, and the service
// is injected into the Wails handler adapter (driving) through the inbound
// port. Only the handler is exposed to the application root.
package category

import (
	"inventario/internal/modules/category/adapters/handler"
	"inventario/internal/modules/category/adapters/repository"
	"inventario/internal/modules/category/core/services"

	"gorm.io/gorm"
)

// Module holds the public surface of the category module.
type Module struct {
	Handler *handler.CategoryHandler
}

// New builds the category module's object graph from a shared DB handle.
func New(db *gorm.DB) *Module {
	repo := repository.NewCategoryRepository(db)
	svc := services.NewCategoryService(repo)

	return &Module{Handler: handler.NewCategoryHandler(svc)}
}
