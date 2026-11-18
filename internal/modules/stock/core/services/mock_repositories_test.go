package services_test

import (
	"context"
	"errors"
	"testing"

	productdomain "inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/stock/core/domain"
	"inventario/internal/modules/stock/core/ports"
	"inventario/internal/modules/stock/core/services"
)

// errDB simula un fallo genérico de la base de datos.
var errDB = errors.New("db error")

// --- fake unit of work ---

// fakeTxKey marca un contexto que ya está dentro de una transacción falsa.
type fakeTxKey struct{}

// fakeUoW ejecuta fn con el contexto del llamador. Como la unidad de trabajo
// real, es reentrante: un Do dentro de otro Do se une y no se cuenta, así que
// calls es el número de transacciones realmente abiertas. No hace rollback;
// esa semántica la cubren los tests de gormdb.UnitOfWork.
type fakeUoW struct {
	calls int
	err   error
}

func (f *fakeUoW) Do(ctx context.Context, fn func(context.Context) error) error {
	if ctx.Value(fakeTxKey{}) != nil {
		return fn(ctx)
	}
	f.calls++
	if f.err != nil {
		return f.err
	}
	return fn(context.WithValue(ctx, fakeTxKey{}, true))
}

// --- mock stock repo ---

var _ ports.StockRepository = (*mockStockRepo)(nil)

type mockStockRepo struct {
	movements []domain.Movement
	nextID    int64

	saveErr error
	findErr error
}

func newMockStockRepo() *mockStockRepo {
	return &mockStockRepo{nextID: 1}
}

func (m *mockStockRepo) Save(_ context.Context, mv *domain.Movement) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	mv.ID = m.nextID
	m.movements = append(m.movements, *mv)
	m.nextID++
	return nil
}

func (m *mockStockRepo) FindByProductID(_ context.Context, productID int64) ([]domain.Movement, error) {
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

var _ ports.ProductRepository = (*mockProductRepo)(nil)

type mockProductRepo struct {
	products map[int64]productdomain.Product
	nextID   int64

	findByIDErr    error
	updateStockErr error
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[int64]productdomain.Product),
		nextID:   1,
	}
}

// seed inserta un producto con el stock indicado y devuelve su ID.
func (m *mockProductRepo) seed(stock int) int64 {
	id := m.nextID
	m.products[id] = productdomain.Product{ID: id, Name: "Laptop", Code: "LAP-001", Price: 999, Stock: stock}
	m.nextID++
	return id
}

func (m *mockProductRepo) stockOf(id int64) int { return m.products[id].Stock }

func (m *mockProductRepo) FindByID(_ context.Context, id int64) (*productdomain.Product, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	p, ok := m.products[id]
	if !ok {
		return nil, productdomain.ErrNotFound
	}
	return &p, nil
}

func (m *mockProductRepo) UpdateStock(_ context.Context, id int64, stock int) error {
	if m.updateStockErr != nil {
		return m.updateStockErr
	}
	p, ok := m.products[id]
	if !ok {
		return productdomain.ErrNotFound
	}
	p.Stock = stock
	m.products[id] = p
	return nil
}

// --- helpers ---

type testDeps struct {
	uow      *fakeUoW
	stock    *mockStockRepo
	products *mockProductRepo
}

// newTestService crea un servicio respaldado por fakes vacíos.
func newTestService(t *testing.T) (*services.StockService, testDeps) {
	t.Helper()
	d := testDeps{uow: &fakeUoW{}, stock: newMockStockRepo(), products: newMockProductRepo()}
	return services.NewStockService(d.uow, d.stock, d.products), d
}
