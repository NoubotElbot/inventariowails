package services_test

import (
	"context"
	"errors"
	"sort"
	"testing"

	productdomain "inventario/internal/modules/product/core/domain"
	"inventario/internal/modules/sale/core/domain"
	"inventario/internal/modules/sale/core/ports"
	"inventario/internal/modules/sale/core/services"
	stockdomain "inventario/internal/modules/stock/core/domain"
	stockports "inventario/internal/modules/stock/core/ports"
	"inventario/internal/platform/pagination"
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

// --- mock sale repo ---

var _ ports.SaleRepository = (*mockSaleRepo)(nil)

type mockSaleRepo struct {
	sales  map[int64]domain.Sale
	nextID int64

	saveErr     error
	findAllErr  error
	findByIDErr error
}

func newMockSaleRepo() *mockSaleRepo {
	return &mockSaleRepo{sales: make(map[int64]domain.Sale), nextID: 1}
}

func (m *mockSaleRepo) Save(_ context.Context, s *domain.Sale) error {
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

// FindAll devuelve la página pedida con las ventas más nuevas (mayor ID) primero.
func (m *mockSaleRepo) FindAll(_ context.Context, q domain.FindAllQuery) (pagination.Result[domain.Sale], error) {
	if m.findAllErr != nil {
		return pagination.Result[domain.Sale]{}, m.findAllErr
	}
	all := make([]domain.Sale, 0, len(m.sales))
	for _, s := range m.sales {
		all = append(all, s)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID > all[j].ID })

	p := q.Params.Normalize()
	start := min(p.Offset(), len(all))
	end := min(start+p.PageSize, len(all))
	total := len(all)
	totalPages := (total + p.PageSize - 1) / p.PageSize

	return pagination.Result[domain.Sale]{
		Data: all[start:end],
		Meta: pagination.Meta{
			Page:       p.Page,
			PageSize:   p.PageSize,
			TotalItems: int64(total),
			TotalPages: totalPages,
			HasNext:    p.Page < totalPages,
			HasPrev:    p.Page > 1 && totalPages > 0,
		},
	}, nil
}

func (m *mockSaleRepo) FindByID(_ context.Context, id int64) (*domain.Sale, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	s, ok := m.sales[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &s, nil
}

// --- mock product repo ---

var _ ports.ProductRepository = (*mockProductRepo)(nil)

type mockProductRepo struct {
	products    map[int64]productdomain.Product
	nextID      int64
	findByIDErr error
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{products: make(map[int64]productdomain.Product), nextID: 1}
}

// seed inserta un producto con el precio, stock y estado indicados.
func (m *mockProductRepo) seed(price int64, stock int, active bool) int64 {
	id := m.nextID
	m.products[id] = productdomain.Product{
		ID: id, Name: "Producto", Code: "P", Price: price, Stock: stock, Active: active,
	}
	m.nextID++
	return id
}

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

// --- fake stock registrar ---

var _ ports.StockMovementRegistrar = (*fakeStockRegistrar)(nil)

// fakeStockRegistrar reemplaza al servicio de stock: registra cada
// movimiento solicitado y puede fallar a pedido.
type fakeStockRegistrar struct {
	calls []stockports.RegisterMovementInput
	err   error
}

func (f *fakeStockRegistrar) RegisterMovement(_ context.Context, in stockports.RegisterMovementInput) (*stockdomain.Movement, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.calls = append(f.calls, in)
	return &stockdomain.Movement{
		ID:           int64(len(f.calls)),
		ProductID:    in.ProductID,
		Quantity:     in.Quantity,
		MovementType: in.Type,
		Notes:        in.Notes,
	}, nil
}

// --- helpers ---

type testDeps struct {
	uow      *fakeUoW
	sales    *mockSaleRepo
	products *mockProductRepo
	stock    *fakeStockRegistrar
}

// newTestService crea un servicio respaldado por fakes vacíos.
func newTestService(t *testing.T) (*services.SaleService, testDeps) {
	t.Helper()
	d := testDeps{
		uow:      &fakeUoW{},
		sales:    newMockSaleRepo(),
		products: newMockProductRepo(),
		stock:    &fakeStockRegistrar{},
	}
	return services.NewSaleService(d.uow, d.sales, d.products, d.stock), d
}

func items(in ...ports.SaleItemInput) ports.CreateInput {
	return ports.CreateInput{Items: in}
}
