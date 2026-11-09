package category_test

import (
	"errors"
	appCategory "inventario/internal/application/category"
	domain "inventario/internal/domain/category"
	"inventario/internal/domain/pagination"
	"testing"
)

func query(p pagination.Params) domain.FindAllQuery { return domain.FindAllQuery{Params: p} }
func queryInactive(p pagination.Params) domain.FindAllQuery {
	return domain.FindAllQuery{Params: p, OnlyInactive: true}
}

func TestListUseCase_Execute(t *testing.T) {
	t.Run("retorna lista vacía sin error", func(t *testing.T) {
		uc := appCategory.NewListUseCase(newMockRepo())

		result, err := uc.Execute(query(pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Items) != 0 || result.Total != 0 {
			t.Errorf("got %d items / total %d, want 0", len(result.Items), result.Total)
		}
	})

	t.Run("retorna categorías existentes", func(t *testing.T) {
		repo := newMockRepo()
		repo.Save(&domain.Category{Name: "Electrónica"})
		repo.Save(&domain.Category{Name: "Ropa"})
		uc := appCategory.NewListUseCase(repo)

		result, err := uc.Execute(query(pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Items) != 2 || result.Total != 2 {
			t.Errorf("got %d items / total %d, want 2", len(result.Items), result.Total)
		}
	})

	t.Run("pagina los resultados", func(t *testing.T) {
		repo := newMockRepo()
		for _, n := range []string{"A", "B", "C", "D", "E"} {
			repo.Save(&domain.Category{Name: n})
		}
		uc := appCategory.NewListUseCase(repo)

		result, err := uc.Execute(query(pagination.Params{Page: 3, PageSize: 2}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Items) != 1 || result.Items[0].Name != "E" {
			t.Errorf("unexpected items: %+v", result.Items)
		}
		if result.Total != 5 || result.TotalPages != 3 || result.Page != 3 {
			t.Errorf("unexpected metadata: %+v", result)
		}
	})

	t.Run("normaliza parámetros inválidos", func(t *testing.T) {
		uc := appCategory.NewListUseCase(newMockRepo())

		result, _ := uc.Execute(query(pagination.Params{Page: -1, PageSize: 9999}))

		if result.Page != 1 || result.PageSize != pagination.MaxPageSize {
			t.Errorf("got page=%d size=%d", result.Page, result.PageSize)
		}
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		repo := newMockRepo()
		repo.findAllErr = errors.New("db error")
		uc := appCategory.NewListUseCase(repo)

		_, err := uc.Execute(query(pagination.Params{}))

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("excluye categorías eliminadas por defecto", func(t *testing.T) {
		repo := newMockRepo()
		active := &domain.Category{Name: "Activa"}
		inactive := &domain.Category{Name: "Inactiva"}
		repo.Save(active)
		repo.Save(inactive)
		repo.SoftDelete(inactive.ID)
		uc := appCategory.NewListUseCase(repo)

		result, err := uc.Execute(query(pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || result.Items[0].Name != "Activa" {
			t.Errorf("expected only active categories, got %+v", result.Items)
		}
	})

	t.Run("OnlyInactive retorna solo categorías eliminadas", func(t *testing.T) {
		repo := newMockRepo()
		active := &domain.Category{Name: "Activa"}
		inactive := &domain.Category{Name: "Inactiva"}
		repo.Save(active)
		repo.Save(inactive)
		repo.SoftDelete(inactive.ID)
		uc := appCategory.NewListUseCase(repo)

		result, err := uc.Execute(queryInactive(pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || result.Items[0].Name != "Inactiva" {
			t.Errorf("expected only inactive categories, got %+v", result.Items)
		}
	})
}
