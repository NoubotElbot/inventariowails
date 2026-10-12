package stock

type Repository interface {
	Save(m *Movement) error
	FindByProductID(productID int64) ([]Movement, error)
}
