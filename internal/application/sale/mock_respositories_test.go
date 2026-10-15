package sale_test

import (
	domainProduct "inventario/internal/domain/product"
	domainSale "inventario/internal/domain/sale"
	domainStock "inventario/internal/domain/stock"
)

// --- mock sale repo ---

type mockSaleRepo struct {
	sales   map[int64]domainSale.Sale
	nextID  int64
	saveErr error
	findErr error
}

func newMockSaleRepo() *mockSaleRepo {
	return &mockSaleRepo{
		sales:  make(map[int64]domainSale.Sale),
		nextID: 1,
	}
}

func (m *mockSaleRepo) Save(s *domainSale.Sale) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	s.ID = m.nextID
	for i := range s.Items {
		s.Items[i].SaleID = s.ID
		s.Items[i].ID = int64(i + 1)
	}
	m.sales[m.nextID] = *s
	m.nextID++
	return nil
}

func (m *mockSaleRepo) FindAll() ([]domainSale.Sale, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	result := make([]domainSale.Sale, 0, len(m.sales))
	for _, s := range m.sales {
		result = append(result, s)
	}
	return result, nil
}

func (m *mockSaleRepo) FindByID(id int64) (*domainSale.Sale, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	s, ok := m.sales[id]
	if !ok {
		return nil, domainSale.ErrNotFound
	}
	return &s, nil
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
func (m *mockProductRepo) FindByCode(code string) (*domainProduct.Product, error) {
	return nil, domainProduct.ErrNotFound
}
func (m *mockProductRepo) Save(p *domainProduct.Product) error   { return nil }
func (m *mockProductRepo) Update(p *domainProduct.Product) error { return nil }
func (m *mockProductRepo) Delete(id int64) error                 { return nil }

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

// --- mock stock repo ---

type mockStockRepo struct {
	movements []domainStock.Movement
	nextID    int64
	saveErr   error
}

func newMockStockRepo() *mockStockRepo {
	return &mockStockRepo{nextID: 1}
}

func (m *mockStockRepo) Save(mv *domainStock.Movement) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	mv.ID = m.nextID
	m.movements = append(m.movements, *mv)
	m.nextID++
	return nil
}

func (m *mockStockRepo) FindByProductID(productID int64) ([]domainStock.Movement, error) {
	return nil, nil
}
