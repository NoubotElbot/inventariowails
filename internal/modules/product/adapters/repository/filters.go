package repository

import (
	"inventario/internal/modules/product/core/domain"

	"gorm.io/gorm"
)

// searchTerm filters by a substring of name, code or barcode; an empty term
// doesn't filter.
func searchTerm(term string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if term == "" {
			return db
		}
		like := "%" + term + "%"
		return db.Where("(name LIKE ? OR code LIKE ? OR barcode LIKE ?)", like, like, like)
	}
}

// inCategories filters by any of the given category IDs; an empty list
// doesn't filter.
func inCategories(ids []int64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(ids) == 0 {
			return db
		}
		return db.Where("category_id IN ?", ids)
	}
}

// findAllFilters turns a FindAllQuery into the scopes to apply.
func findAllFilters(q domain.FindAllQuery) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(searchTerm(q.Search), inCategories(q.CategoryIDs))
	}
}
