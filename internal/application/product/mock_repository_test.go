package product_test

import (
	domain "inventario/internal/domain/product"
)

type mockRepo struct {
	products map[int64]domain.Product
	nextID   int64

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

func (m *mockRepo) FindAll() ([]domain.Product, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	result := make([]domain.Product, 0, len(m.products))
	for _, p := range m.products {
		result = append(result, p)
	}
	return result, nil
}

func (m *mockRepo) FindByID(id int64) (*domain.Product, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	p, ok := m.products[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &p, nil
}

func (m *mockRepo) FindByCode(code string) (*domain.Product, error) {
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

func (m *mockRepo) Save(p *domain.Product) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	p.ID = m.nextID
	m.products[m.nextID] = *p
	m.nextID++
	return nil
}

func (m *mockRepo) Update(p *domain.Product) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.products[p.ID]; !ok {
		return domain.ErrNotFound
	}
	m.products[p.ID] = *p
	return nil
}

func (m *mockRepo) Delete(id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.products[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.products, id)
	return nil
}

func (m *mockRepo) UpdateStock(id int64, stock int) error {
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
