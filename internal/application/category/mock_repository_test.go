package category_test

import (
	domain "inventario/internal/domain/category"
)

type mockRepo struct {
	categories map[int64]domain.Category
	nextID     int64

	// hooks para forzar errores puntuales
	findAllErr  error
	findByIDErr error
	saveErr     error
	updateErr   error
	deleteErr   error
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		categories: make(map[int64]domain.Category),
		nextID:     1,
	}
}

func (m *mockRepo) FindAll() ([]domain.Category, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	result := make([]domain.Category, 0, len(m.categories))
	for _, c := range m.categories {
		result = append(result, c)
	}
	return result, nil
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

func (m *mockRepo) Delete(id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.categories[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.categories, id)
	return nil
}
