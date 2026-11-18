package services_test

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"inventario/internal/modules/category/core/domain"
	"inventario/internal/modules/category/core/ports"
	"inventario/internal/modules/category/core/services"
	"inventario/internal/platform/pagination"
)

// errDB simula un fallo genérico de la base de datos.
var errDB = errors.New("db error")

var _ ports.CategoryRepository = (*mockRepo)(nil)

// mockRepo es un fake en memoria del repositorio. FindAll no filtra ni pagina:
// registra la query recibida y devuelve findAllResult, porque esa lógica
// pertenece al repositorio real, no al servicio.
type mockRepo struct {
	categories map[int64]domain.Category
	nextID     int64

	findAllResult pagination.Result[domain.Category]
	lastQuery     domain.FindAllQuery

	// hooks para forzar errores puntuales
	findAllErr     error
	findOptionsErr error
	findByIDErr    error
	saveErr        error
	updateErr      error
	softDeleteErr  error
	restoreErr     error
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		categories: make(map[int64]domain.Category),
		nextID:     1,
	}
}

// newTestService crea un servicio respaldado por un mockRepo vacío.
func newTestService(t *testing.T) (*services.CategoryService, *mockRepo) {
	t.Helper()
	repo := newMockRepo()
	return services.NewCategoryService(repo), repo
}

// seedCategory inserta una categoría activa directamente en el fake.
func seedCategory(t *testing.T, repo *mockRepo, name string) domain.Category {
	t.Helper()
	c := domain.Category{ID: repo.nextID, Name: name}
	repo.categories[c.ID] = c
	repo.nextID++
	return c
}

// seedDeletedCategory inserta una categoría ya eliminada (soft delete).
func seedDeletedCategory(t *testing.T, repo *mockRepo, name string) domain.Category {
	t.Helper()
	c := seedCategory(t, repo, name)
	now := time.Now()
	c.DeletedAt = &now
	repo.categories[c.ID] = c
	return c
}

func (m *mockRepo) FindAll(_ context.Context, q domain.FindAllQuery) (pagination.Result[domain.Category], error) {
	m.lastQuery = q
	if m.findAllErr != nil {
		return pagination.Result[domain.Category]{}, m.findAllErr
	}
	return m.findAllResult, nil
}

func (m *mockRepo) FindOptions(_ context.Context) ([]domain.Option, error) {
	if m.findOptionsErr != nil {
		return nil, m.findOptionsErr
	}
	opts := make([]domain.Option, 0, len(m.categories))
	for _, c := range m.categories {
		opts = append(opts, domain.Option{ID: c.ID, Name: c.Name})
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].Name < opts[j].Name })
	return opts, nil
}

func (m *mockRepo) FindByID(_ context.Context, id int64) (*domain.Category, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	c, ok := m.categories[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &c, nil
}

func (m *mockRepo) Save(_ context.Context, c *domain.Category) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	c.ID = m.nextID
	m.categories[m.nextID] = *c
	m.nextID++
	return nil
}

func (m *mockRepo) Update(_ context.Context, c *domain.Category) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.categories[c.ID]; !ok {
		return domain.ErrNotFound
	}
	m.categories[c.ID] = *c
	return nil
}

func (m *mockRepo) SoftDelete(_ context.Context, id int64) error {
	if m.softDeleteErr != nil {
		return m.softDeleteErr
	}
	c, ok := m.categories[id]
	if !ok {
		return domain.ErrNotFound
	}
	now := time.Now()
	c.DeletedAt = &now
	m.categories[id] = c
	return nil
}

func (m *mockRepo) Restore(_ context.Context, id int64) error {
	if m.restoreErr != nil {
		return m.restoreErr
	}
	c, ok := m.categories[id]
	if !ok {
		return domain.ErrNotFound
	}
	c.DeletedAt = nil
	m.categories[id] = c
	return nil
}
