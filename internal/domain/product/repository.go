package product

type Repository interface {
	FindAll() ([]Product, error)
	FindByID(id int64) (*Product, error)
	FindByCode(code string) (*Product, error)
	Save(p *Product) error
	Update(p *Product) error
	Delete(id int64) error
	UpdateStock(id int64, stock int) error
}
