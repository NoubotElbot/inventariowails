package sale

type Repository interface {
	Save(s *Sale) error
	FindAll() ([]Sale, error)
	FindByID(id int64) (*Sale, error)
}
