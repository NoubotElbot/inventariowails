package services_test

import (
	"context"
	"errors"
	"testing"

	"inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/product/core/ports"
	"inventario/internal/modules/product/core/services"
	"inventario/internal/platform/pagination"
)

// errDB simula un fallo genérico de la base de datos.
var errDB = errors.New("db error")

var _ ports.ProductRepository = (*mockRepo)(nil)

// mockRepo es un fake en memoria del repositorio. FindAll no filtra ni pagina:
// registra la query recibida y devuelve findAllResult, porque esa lógica
// pertenece al repositorio real, no al servicio.
type mockRepo struct {
	products map[int64]domain.Product
	nextID   int64

	findAllResult pagination.Result[domain.Product]
	lastQuery     domain.FindAllQuery

	// hooks para forzar errores puntuales
	findAllErr     error
	findByIDErr    error
	findByCodeErr  error
	saveErr        error
	updateErr      error
	deleteErr      error
	updateStockErr error
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		products: make(map[int64]domain.Product),
		nextID:   1,
	}
}

// newTestService crea un servicio respaldado por un mockRepo vacío.
func newTestService(t *testing.T) (*services.ProductService, *mockRepo) {
	t.Helper()
	repo := newMockRepo()
	return services.NewProductService(repo), repo
}

// seedProduct inserta un producto válido y activo directamente en el fake.
func seedProduct(t *testing.T, repo *mockRepo, code string) domain.Product {
	t.Helper()
	p := domain.Product{
		ID:     repo.nextID,
		Name:   "Producto " + code,
		Code:   code,
		Price:  999,
		Stock:  10,
		Active: true,
	}
	repo.products[p.ID] = p
	repo.nextID++
	return p
}

func (m *mockRepo) FindAll(_ context.Context, q domain.FindAllQuery) (pagination.Result[domain.Product], error) {
	m.lastQuery = q
	if m.findAllErr != nil {
		return pagination.Result[domain.Product]{}, m.findAllErr
	}
	return m.findAllResult, nil
}

func (m *mockRepo) FindByID(_ context.Context, id int64) (*domain.Product, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	p, ok := m.products[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &p, nil
}

func (m *mockRepo) FindByCode(_ context.Context, code string) (*domain.Product, error) {
	if m.findByCodeErr != nil {
		return nil, m.findByCodeErr
	}
	for _, p := range m.products {
		if p.Code == code {
			return &p, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockRepo) Save(_ context.Context, p *domain.Product) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	p.ID = m.nextID
	m.products[m.nextID] = *p
	m.nextID++
	return nil
}

func (m *mockRepo) Update(_ context.Context, p *domain.Product) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.products[p.ID]; !ok {
		return domain.ErrNotFound
	}
	m.products[p.ID] = *p
	return nil
}

func (m *mockRepo) Delete(_ context.Context, id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.products[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.products, id)
	return nil
}

func (m *mockRepo) UpdateStock(_ context.Context, id int64, stock int) error {
	if m.updateStockErr != nil {
		return m.updateStockErr
	}
	p, ok := m.products[id]
	if !ok {
		return domain.ErrNotFound
	}
	p.Stock = stock
	m.products[id] = p
	return nil
}
