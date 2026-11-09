package product_test

import (
	"errors"
	appProduct "inventario/internal/application/product"
	domain "inventario/internal/domain/product"
	"inventario/internal/domain/pagination"
	"testing"
)

func ptrInt64(v int64) *int64  { return &v }
func ids(vs ...int64) []int64  { return vs }

func q(search string, catIDs []int64, p pagination.Params) domain.FindAllQuery {
	return domain.FindAllQuery{Search: search, CategoryIDs: catIDs, Params: p}
}

func TestListUseCase_Execute(t *testing.T) {
	t.Run("retorna lista vacía sin error", func(t *testing.T) {
		uc := appProduct.NewListUseCase(newMockRepo())

		result, err := uc.Execute(q("", nil, pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Items) != 0 || result.Total != 0 {
			t.Errorf("got %d items / total %d, want 0", len(result.Items), result.Total)
		}
	})

	t.Run("retorna productos existentes sin filtro", func(t *testing.T) {
		repo := newMockRepo()
		repo.Save(&domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999})
		repo.Save(&domain.Product{Name: "Mouse", Code: "MOU-001", Price: 29})
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute(q("", nil, pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Items) != 2 || result.Total != 2 {
			t.Errorf("got %d items / total %d, want 2", len(result.Items), result.Total)
		}
	})

	t.Run("filtra por nombre con search", func(t *testing.T) {
		repo := newMockRepo()
		repo.Save(&domain.Product{Name: "Laptop Gamer", Code: "LAP-001", Price: 999})
		repo.Save(&domain.Product{Name: "Mouse USB", Code: "MOU-001", Price: 29})
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute(q("laptop", nil, pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || result.Items[0].Name != "Laptop Gamer" {
			t.Errorf("unexpected items: %+v", result.Items)
		}
	})

	t.Run("filtra por código con search", func(t *testing.T) {
		repo := newMockRepo()
		repo.Save(&domain.Product{Name: "Laptop", Code: "LAP-001", Price: 999})
		repo.Save(&domain.Product{Name: "Mouse", Code: "MOU-001", Price: 29})
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute(q("MOU", nil, pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || result.Items[0].Code != "MOU-001" {
			t.Errorf("unexpected items: %+v", result.Items)
		}
	})

	t.Run("filtra por barcode con search", func(t *testing.T) {
		repo := newMockRepo()
		repo.Save(&domain.Product{Name: "Laptop", Code: "LAP-001", Barcode: "123456", Price: 999})
		repo.Save(&domain.Product{Name: "Mouse", Code: "MOU-001", Barcode: "789000", Price: 29})
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute(q("1234", nil, pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || result.Items[0].Barcode != "123456" {
			t.Errorf("unexpected items: %+v", result.Items)
		}
	})

	t.Run("filtra por category_id", func(t *testing.T) {
		repo := newMockRepo()
		catA := ptrInt64(1)
		catB := ptrInt64(2)
		repo.Save(&domain.Product{Name: "Laptop", Code: "LAP-001", CategoryID: catA, Price: 999})
		repo.Save(&domain.Product{Name: "Mouse", Code: "MOU-001", CategoryID: catA, Price: 29})
		repo.Save(&domain.Product{Name: "Desk", Code: "DSK-001", CategoryID: catB, Price: 200})
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute(q("", ids(1), pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 2 {
			t.Errorf("got total %d, want 2", result.Total)
		}
	})

	t.Run("filtra por múltiples category_ids", func(t *testing.T) {
		repo := newMockRepo()
		catA := ptrInt64(1)
		catB := ptrInt64(2)
		catC := ptrInt64(3)
		repo.Save(&domain.Product{Name: "Laptop", Code: "LAP-001", CategoryID: catA, Price: 999})
		repo.Save(&domain.Product{Name: "Mouse", Code: "MOU-001", CategoryID: catB, Price: 29})
		repo.Save(&domain.Product{Name: "Desk", Code: "DSK-001", CategoryID: catC, Price: 200})
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute(q("", ids(1, 2), pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 2 {
			t.Errorf("got total %d, want 2", result.Total)
		}
	})

	t.Run("combina search y category_ids", func(t *testing.T) {
		repo := newMockRepo()
		catA := ptrInt64(1)
		catB := ptrInt64(2)
		repo.Save(&domain.Product{Name: "Laptop Gamer", Code: "LAP-001", CategoryID: catA, Price: 999})
		repo.Save(&domain.Product{Name: "Laptop Básica", Code: "LAP-002", CategoryID: catB, Price: 500})
		repo.Save(&domain.Product{Name: "Mouse", Code: "MOU-001", CategoryID: catA, Price: 29})
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute(q("laptop", ids(1), pagination.Params{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || result.Items[0].Code != "LAP-001" {
			t.Errorf("unexpected items: %+v", result.Items)
		}
	})

	t.Run("pagina los resultados", func(t *testing.T) {
		repo := newMockRepo()
		for i, name := range []string{"A", "B", "C", "D", "E"} {
			repo.Save(&domain.Product{Name: name, Code: name, Price: int64(i + 1)})
		}
		uc := appProduct.NewListUseCase(repo)

		result, err := uc.Execute(q("", nil, pagination.Params{Page: 3, PageSize: 2}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Items) != 1 {
			t.Errorf("got %d items, want 1", len(result.Items))
		}
		if result.Total != 5 || result.TotalPages != 3 || result.Page != 3 {
			t.Errorf("unexpected metadata: %+v", result)
		}
	})

	t.Run("normaliza parámetros inválidos", func(t *testing.T) {
		uc := appProduct.NewListUseCase(newMockRepo())

		result, _ := uc.Execute(q("", nil, pagination.Params{Page: -1, PageSize: 9999}))

		if result.Page != 1 || result.PageSize != pagination.MaxPageSize {
			t.Errorf("got page=%d size=%d", result.Page, result.PageSize)
		}
	})

	t.Run("propaga error del repositorio", func(t *testing.T) {
		repo := newMockRepo()
		repo.findAllErr = errors.New("db error")
		uc := appProduct.NewListUseCase(repo)

		_, err := uc.Execute(q("", nil, pagination.Params{}))

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
