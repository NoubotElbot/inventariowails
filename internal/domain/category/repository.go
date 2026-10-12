package category

type Repository interface {
	FindAll() ([]Category, error)
	FindByID(id int64) (*Category, error)
	Save(c *Category) error
	Update(c *Category) error
	Delete(id int64) error
}
