package repository

import (
	"inventario/internal/modules/category/core/domain"

	"gorm.io/gorm"
)

func active(db *gorm.DB) *gorm.DB  { return db.Where("deleted_at IS NULL") }
func deleted(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NOT NULL") }

// nameContains filters by name substring; an empty term doesn't filter.
func nameContains(term string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if term == "" {
			return db
		}
		return db.Where("name LIKE ?", "%"+term+"%")
	}
}

// findAllFilters turns a FindAllQuery into the scopes to apply.
func findAllFilters(q domain.FindAllQuery) func(*gorm.DB) *gorm.DB {
	state := active
	if q.OnlyInactive {
		state = deleted
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(state, nameContains(q.Search))
	}
}
