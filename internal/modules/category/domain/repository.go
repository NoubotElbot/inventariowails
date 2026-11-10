package domain

import "inventario/internal/platform/pagination"

// Option es una vista mínima (id, name) para componentes de selección.
type Option struct {
	ID   int64
	Name string
}

type Repository interface {
	FindOptions() ([]Option, error)
	FindAll(q FindAllQuery) (pagination.Result[Category], error)
	FindByID(id int64) (*Category, error)
	Save(c *Category) error
	Update(c *Category) error
	SoftDelete(id int64) error
	Restore(id int64) error
}
