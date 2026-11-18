// Package pagination ofrece paginación reutilizable sobre GORM usando generics.
//
// Uso básico (repositorio):
//
//	func (r *UserRepo) List(p pagination.Params) (*pagination.Result[User], error) {
//		return pagination.Paginate[User](r.db.Where("active = ?", true).Order("id DESC"), p)
//	}
//
// Uso en handler Gin:
//
//	var p pagination.Params
//	_ = c.ShouldBindQuery(&p) // ?page=2&page_size=20
//	res, err := repo.List(p)
//
// Solo el scope (si quieres controlar el Find tú mismo):
//
//	db.Scopes(pagination.Scope(p)).Find(&users)
package pagination

import (
	"math"

	"gorm.io/gorm"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Params son los parámetros de entrada (query string o JSON).
// Los tags `default` hacen que Gin rellene los valores si faltan en la query.
type Params struct {
	Page     int `form:"page,default=1" json:"page"`
	PageSize int `form:"page_size,default=10" json:"page_size"`
}

// NewParams crea Params ya con los valores por defecto.
// Útil para JSON (donde no aplican los tags `default`) o para uso interno:
//
//	p := pagination.NewParams()
//	_ = c.ShouldBindJSON(&p) // los campos ausentes conservan el default
func NewParams() Params {
	return Params{Page: DefaultPage, PageSize: DefaultPageSize}
}

// Normalize aplica valores por defecto y límites.
func (p Params) Normalize() Params {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
	return p
}

// Offset calcula el OFFSET SQL a partir de los parámetros normalizados.
func (p Params) Offset() int {
	p = p.Normalize()
	return (p.Page - 1) * p.PageSize
}

// Meta describe el estado de la paginación.
type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Result es la respuesta paginada genérica.
type Result[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

// Scope devuelve un scope de GORM que aplica LIMIT y OFFSET.
func Scope(p Params) func(*gorm.DB) *gorm.DB {
	p = p.Normalize()
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.Offset()).Limit(p.PageSize)
	}
}

// Paginate ejecuta COUNT + SELECT sobre la query recibida y devuelve el resultado paginado.
// La query `db` conserva sus Where/Joins/Preload/Order; no se modifica el original.
func Paginate[T any](db *gorm.DB, p Params) (Result[T], error) {
	p = p.Normalize()

	var (
		model T
		total int64
	)
	if err := db.Session(&gorm.Session{}).Model(&model).Count(&total).Error; err != nil {
		return Result[T]{}, err
	}

	items := make([]T, 0, p.PageSize)
	if total > 0 {
		if err := db.Session(&gorm.Session{}).Scopes(Scope(p)).Find(&items).Error; err != nil {
			return Result[T]{}, err
		}
	}

	return Result[T]{Data: items, Meta: buildMeta(p, total)}, nil
}

func buildMeta(p Params, total int64) Meta {
	totalPages := int(math.Ceil(float64(total) / float64(p.PageSize)))
	return Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
		HasNext:    p.Page < totalPages,
		HasPrev:    p.Page > 1 && totalPages > 0,
	}
}
