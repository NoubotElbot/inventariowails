package stock_test

import (
	domainProduct "inventario/internal/domain/product"
	domain "inventario/internal/domain/stock"
)

// --- mock stock repo ---

type mockStockRepo struct {
	movements []domain.Movement
	nextID    int64
	saveErr   error
	findErr   error
}

func newMockStockRepo() *mockStockRepo {
	return &mockStockRepo{nextID: 1}
}

func (m *mockStockRepo) Save(mv *domain.Movement) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	mv.ID = m.nextID
	m.movements = append(m.movements, *mv)
	m.nextID++
	return nil
}

func (m *mockStockRepo) FindByProductID(productID int64) ([]domain.Movement, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	result := make([]domain.Movement, 0)
	for _, mv := range m.movements {
		if mv.ProductID == productID {
			result = append(result, mv)
		}
	}
	return result, nil
}

// --- mock product repo ---

type mockProductRepo struct {
	products       map[int64]domainProduct.Product
	nextID         int64
	findByIDErr    error
	updateStockErr error
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[int64]domainProduct.Product),
		nextID:   1,
	}
}

func (m *mockProductRepo) seed(p *domainProduct.Product) {
	p.ID = m.nextID
	m.products[m.nextID] = *p
	m.nextID++
}

func (m *mockProductRepo) FindAll() ([]domainProduct.Product, error) { return nil, nil }

func (m *mockProductRepo) FindByID(id int64) (*domainProduct.Product, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	p, ok := m.products[id]
	if !ok {
		return nil, domainProduct.ErrNotFound
	}
	return &p, nil
}

func (m *mockProductRepo) FindByCode(code string) (*domainProduct.Product, error) {
	return nil, nil
}

func (m *mockProductRepo) Save(p *domainProduct.Product) error { return nil }

func (m *mockProductRepo) Update(p *domainProduct.Product) error { return nil }

func (m *mockProductRepo) Delete(id int64) error { return nil }

func (m *mockProductRepo) UpdateStock(id int64, stock int) error {
	if m.updateStockErr != nil {
		return m.updateStockErr
	}
	p, ok := m.products[id]
	if !ok {
		return domainProduct.ErrNotFound
	}
	p.Stock = stock
	m.products[id] = p
	return nil
}
