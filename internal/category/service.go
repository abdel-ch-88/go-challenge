package category

import "errors"

type CategoryService struct {
	catRepo Repository
}

func NewService(r Repository) Service {
	return &CategoryService{
		catRepo: r,
	}
}

func (srv *CategoryService) GetAllCategories() ([]Category, error) {
	return srv.catRepo.GetAllCategories()
}

func (srv *CategoryService) CreateCategory(newCat Category) error {
	if newCat.Code == "" || newCat.Name == "" {
		return errors.New("invalid Category data")
	}

	existCat, err := srv.catRepo.GetCategoryByCode(newCat.Code)
	if err != nil {
		return err
	}
	if existCat != nil && existCat.ID != 0 {
		return errors.New("categody code already exists")
	}

	return srv.catRepo.Save(newCat)
}
