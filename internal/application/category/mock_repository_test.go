package category_test

import (
	domain "inventario/internal/domain/category"
	"inventario/internal/domain/pagination"
	"sort"
	"strings"
	"time"
)

func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

type mockRepo struct {
	categories map[int64]domain.Category
	nextID     int64

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

func (m *mockRepo) FindAll(q domain.FindAllQuery) (pagination.Result[domain.Category], error) {
	if m.findAllErr != nil {
		return pagination.Result[domain.Category]{}, m.findAllErr
	}
	all := make([]domain.Category, 0, len(m.categories))
	for _, c := range m.categories {
		if q.OnlyInactive && c.DeletedAt == nil {
			continue
		}
		if !q.OnlyInactive && c.DeletedAt != nil {
			continue
		}
		if q.Search == "" || containsFold(c.Name, q.Search) {
			all = append(all, c)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })

	p := q.Params.Normalize()
	start := min(p.Offset(), len(all))
	end := min(start+p.Limit(), len(all))
	return pagination.NewResult(all[start:end], int64(len(all)), p), nil
}

func (m *mockRepo) FindOptions() ([]domain.Option, error) {
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

func (m *mockRepo) FindByID(id int64) (*domain.Category, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	c, ok := m.categories[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &c, nil
}

func (m *mockRepo) Save(c *domain.Category) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	c.ID = m.nextID
	m.categories[m.nextID] = *c
	m.nextID++
	return nil
}

func (m *mockRepo) Update(c *domain.Category) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.categories[c.ID]; !ok {
		return domain.ErrNotFound
	}
	m.categories[c.ID] = *c
	return nil
}

func (m *mockRepo) SoftDelete(id int64) error {
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

func (m *mockRepo) Restore(id int64) error {
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
