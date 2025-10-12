package product

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (srv *Service) FindProducts() ([]Product, error) {
	return srv.repository.GetAllProducts()
}
